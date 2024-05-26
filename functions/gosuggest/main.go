package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"blinders/functions/suggest/core"
	"blinders/packages/auth"
	"blinders/packages/db/collectingdb"
	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"
	"blinders/packages/transport"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	transportCh chan transport.Transport
	authCh      chan auth.Manager
	usersRepoCh chan *usersdb.UsersRepo
	brrcCh      chan *bedrockruntime.Client
)

func init() {
	st := time.Now()
	env := os.Getenv("ENVIRONMENT")
	log.Printf("GoSuggest api running on %s environment\n", env)

	transportCh = InitTransport(context.Background())
	usersRepoCh = make(chan *usersdb.UsersRepo)
	go func() {
		dbCh := dbutils.InitMongoDatabaseChanFromEnv("USERS")
		usersDB := <-dbCh
		usersRepoCh <- usersdb.NewUsersRepo(usersDB)
	}()

	authCh = auth.InitFirebaseManagerFromFile("firebase.admin.json")
	brrcCh = core.InitBedrockRuntimeClient(context.Background())

	en := time.Now()
	log.Printf("Cold-start duration: %v ms", en.Sub(st).Milliseconds())
}

func HandleRequest(
	ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	authUser, ok := ctx.Value(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panicln("unexpected err, authUser not included in ctx")
	}

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

	brrc := <-brrcCh
	explanation, err := core.ExplainPhraseInSentence(*brrc, phrase, sentence)
	if err != nil {
		log.Println("error when explaining: ", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("cannot explain \"%s\"", phrase),
			Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		}, nil
	}

	userID, _ := primitive.ObjectIDFromHex(authUser.ID)
	event := transport.AddExplainLogEvent{
		Event: transport.Event{Type: transport.AddTranslateLog},
		Payload: collectingdb.ExplainLog{
			UserID:   userID,
			Request:  collectingdb.ExplainRequest{Text: phrase, Sentence: sentence},
			Response: *explanation,
		},
	}

	transporter := <-transportCh
	eventPayload, _ := json.Marshal(event)
	if err := transporter.Push(
		ctx,
		transporter.ConsumerID(transport.CollectingPush),
		eventPayload,
	); err != nil {
		log.Printf("cannot push event to collecting service, err: %v\n", err)
	}

	resInBytes, _ := json.Marshal(explanation)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(resInBytes),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
	}, nil
}

func main() {
	lambda.Start(auth.LambdaAuthMiddlewareFromChan(authCh, usersRepoCh)(HandleRequest))
}
