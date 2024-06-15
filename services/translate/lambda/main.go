package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"blinders/packages/translate"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type TranslateResponse struct {
	Text       string `json:"text"`
	Translated string `json:"translated"`
	Languages  string `json:"languages"`
}

var translator translate.Translator

func init() {
	env := os.Getenv("ENVIRONMENT")
	log.Printf("Translate api running on %s environment\n", env)

	translator = translate.YandexTranslator{APIKey: os.Getenv("YANDEX_API_KEY")}
}

func HandleRequest(
	_ context.Context,
	req events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
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

	resInBytes, _ := json.Marshal(res)
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
	lambda.Start(HandleRequest)
}
