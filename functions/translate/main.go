package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"blinders/packages/auth"
	"blinders/packages/collecting"
	"blinders/packages/db"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	translator = translate.YandexTranslator{APIKey: os.Getenv("YANDEX_API_KEY")}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("failed too load aws config", cfg)
	}

	transporter = transport.NewLambdaTransport(cfg)
	consumerMap = transport.ConsumerMap{
		transport.CollectingPush: os.Getenv("COLLECTING_PUSH_FUNCTION_NAME"),
	}

	url := fmt.Sprintf(
		db.MongoURLTemplate,
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_DATABASE"),
	)

	database := db.NewMongoManager(url, os.Getenv("MONGO_DATABASE"))
	if database == nil {
		log.Fatal("cannot create database manager")
	}
	adminConfig, err := utils.GetFile("firebase.admin.json")
	if err != nil {
		log.Fatal(err)
	}
	authManager, err := auth.NewFirebaseManager(adminConfig)
	if err != nil {
		log.Fatal(err)
	}

	// mock this
	authMiddleware = auth.LambdaAuthMiddleware(authManager, database.Users, auth.MiddlewareOptions{CheckUser: false})
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
			StatusCode: 400,
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
			StatusCode: 400,
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

	//  push event to collecting service
	translateEvent := collecting.TranslateEvent{
		UserID: userOID,
		Request: collecting.TranslateRequest{
			Text: text,
		},
		Response: collecting.TranslateResponse{
			Translate: translated,
		},
	}

	event := transport.CollectEventRequest{
		Request: transport.Request{Type: transport.CollectEvent},
		Data:    collecting.NewGenericEvent(collecting.EventTypeTranslate, translateEvent),
	}

	eventPayload, err := json.Marshal(event)
	if err != nil {
		log.Printf("cannot marshal collect event request, err: %v\n", err)
	} else {
		if err := transporter.Push(
			ctx,
			consumerMap[transport.CollectingPush],
			eventPayload,
		); err != nil {
			log.Printf("cannot push event to collecting service, err: %v\n", err)
		}
	}

	// respond result to user
	resInBytes, _ := json.Marshal(res)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(resInBytes),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
	}, nil
}

func main() {
	lambda.Start(authMiddleware(HandleRequest))
}
