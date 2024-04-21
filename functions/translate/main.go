package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"blinders/packages/collecting"
	"blinders/packages/translate"
	"blinders/packages/transport"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

type TranslatePayload struct {
	Text string `json:"text"`
}

type TranslateResponse struct {
	Text       string `json:"text"`
	Translated string `json:"translated"`
	Languages  string `json:"languages"`
}

var (
	translator  translate.Translator
	transporter transport.Transport
	consumerMap transport.ConsumerMap
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	translator = translate.YandexTranslator{APIKey: os.Getenv("YANDEX_API_KEY")}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("failed too load aws config", cfg)
	}

	transporter = transport.NewLambdaTransport(cfg)
	consumerMap = transport.ConsumerMap{
		transport.Collecting: os.Getenv("COLLECTING_FUNCTION_NAME"),
	}
}

func HandleRequest(
	ctx context.Context,
	event events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	text, ok := event.QueryStringParameters["text"]
	if !ok {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       "required text param",
			Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		}, nil
	}

	langs, ok := event.QueryStringParameters["languages"]
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

	//  push event to collecting service
	translateEvent := collecting.TranslateEvent{
		Request: collecting.TranslateRequest{
			Text: text,
		},
		Response: collecting.TranslateResponse{
			Translate: translated,
		},
	}

	eventPayload, _ := json.Marshal(
		transport.CollectEventRequest{
			Request: transport.Request{
				Type: transport.CollectEvent,
			},
			Data: collecting.NewGenericEvent(collecting.EventTypeTranslate, translateEvent),
		},
	)

	if err := transporter.Push(
		ctx,
		consumerMap[transport.Collecting],
		eventPayload,
	); err != nil {
		log.Printf("cannot push event to collecting service, err: %v\n", err)
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
	lambda.Start(HandleRequest)
}
