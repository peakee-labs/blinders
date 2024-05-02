package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"blinders/packages/auth"
	"blinders/packages/db/collectingdb"
	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"
	"blinders/packages/translate"
	"blinders/packages/transport"
	"blinders/packages/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TranslateResponse struct {
	Text       string `json:"text"`
	Translated string `json:"translated"`
	Languages  string `json:"languages"`
}

var (
	translator     translate.Translator
	transporter    transport.Transport
	consumerMap    transport.ConsumerMap
	authMiddleware auth.LambdaMiddleware
)

func init() {
	env := os.Getenv("ENVIRONMENT")
	log.Printf("Translate api running on %s environment\n", env)

	translator = translate.YandexTranslator{APIKey: os.Getenv("YANDEX_API_KEY")}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("failed too load aws config", cfg)
	}
	transporter = transport.NewLambdaTransport(cfg)
	consumerMap = transport.ConsumerMap{
		transport.CollectingPush: os.Getenv("COLLECTING_PUSH_FUNCTION_NAME"),
	}

	mongoInfo := dbutils.GetMongoInfoFromEnv()
	client, err := dbutils.InitMongoClient(mongoInfo.URL)
	if err != nil {
		log.Fatal(err)
	}

	adminConfig, err := utils.GetFile("firebase.admin.json")
	if err != nil {
		log.Fatal(err)
	}
	authManager, err := auth.NewFirebaseManager(adminConfig)
	if err != nil {
		log.Fatal(err)
	}

	authMiddleware = auth.LambdaAuthMiddleware(
		authManager,
		usersdb.NewUsersRepo(client.Database(mongoInfo.DBName)),
	)
}

func HandleRequest(
	ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	authUser, ok := ctx.Value(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panicln("unexpected err, authUser not included in ctx")
	}

	text, ok := req.QueryStringParameters["text"]
	if !ok {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "required text param",
			Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		}, nil
	}

	langs, ok := req.QueryStringParameters["languages"]
	if !ok {
		langs = "en-vi"
	}

	translated, err := translator.Translate(text, translate.Languages(langs))
	if err != nil {
		log.Println("error translating: ", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("cannot translate \"%s\"", text),
			Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		}, nil
	}

	res := TranslateResponse{
		Text:       text,
		Translated: translated,
		Languages:  langs,
	}

	userOID, _ := primitive.ObjectIDFromHex(authUser.ID)

	event := transport.Event{
		Type: transport.AddTranslateLog,
		Payload: collectingdb.TranslateLog{
			UserID:   userOID,
			Request:  collectingdb.TranslateRequest{Text: text},
			Response: collectingdb.TranslateResponse{Translate: translated},
		},
	}

	eventPayload, _ := json.Marshal(event)
	if err := transporter.Push(
		ctx,
		consumerMap[transport.CollectingPush],
		eventPayload,
	); err != nil {
		log.Printf("cannot push event to collecting service, err: %v\n", err)
	}

	resInBytes, _ := json.Marshal(res)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(resInBytes),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
	}, nil
}

func main() {
	lambda.Start(authMiddleware(HandleRequest))
}
