package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"blinders/packages/auth"
	"blinders/packages/db/collectingdb"
	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"
	"blinders/packages/transport"
	"blinders/services/suggest/core"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	transporter    transport.Transport
	brrc           *bedrockruntime.Client
	authMiddleware auth.LambdaMiddleware
)

func init() {
	st := time.Now()
	env := os.Getenv("ENVIRONMENT")
	log.Printf("GoSuggest api running on %s environment\n", env)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		transporter = InitTransportSync(context.Background())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		userDB, err := dbutils.InitMongoDatabaseFromEnv("USERS")
		if err != nil {
			log.Fatalf("can not init database: %v", userDB)
		}
		usersRepo := usersdb.NewUsersRepo(userDB)
		am, err := auth.NewFirebaseManagerFromFile("firebase.admin.json")
		if err != nil {
			log.Fatalf("can not create auth manager: %v", err)
		}
		authMiddleware = auth.LambdaAuthMiddleware(am, usersRepo)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		brrc = core.InitBedrockRuntimeClientSync(context.Background())
	}()

	wg.Wait()
	en := time.Now()
	log.Printf("Cold-start duration: %v ms", en.Sub(st).Milliseconds())
}

func HandleRequest(
	ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	phrase, ok := req.QueryStringParameters["phrase"]
	if !ok {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "required phrase param",
			Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		}, nil
	}
	sentence, ok := req.QueryStringParameters["sentence"]
	if !ok {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "required sentence param",
			Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		}, nil
	}

	explanation, err := core.ExplainPhraseInSentence(brrc, phrase, sentence)
	if err != nil {
		log.Println("error when explaining: ", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("cannot explain \"%s\"", phrase),
			Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		}, nil
	}

	authUser, ok := ctx.Value(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panicln("unexpected err, authUser not included in ctx")
	}
	userID, _ := primitive.ObjectIDFromHex(authUser.ID)
	event := transport.AddExplainLogEvent{
		Event: transport.Event{Type: transport.AddExplainLog},
		Payload: collectingdb.ExplainLog{
			UserID:   userID,
			Request:  collectingdb.ExplainRequest{Text: phrase, Sentence: sentence},
			Response: *explanation,
		},
	}

	eventPayload, _ := json.Marshal(event)
	if err := transporter.Push(
		ctx,
		transporter.ConsumerID(transport.CollectingPush),
		eventPayload,
	); err != nil {
		log.Printf("cannot push event to collecting service, err: %v\n", err)
	}

	resInBytes, _ := json.Marshal(*explanation)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(resInBytes),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Content-Type":                "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(authMiddleware(HandleRequest))
}
