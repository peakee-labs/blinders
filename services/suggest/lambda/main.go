package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"blinders/packages/apigateway"
	"blinders/services/suggest"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

var brrc *bedrockruntime.Client

func init() {
	env := os.Getenv("ENVIRONMENT")
	log.Println("Peakee Suggest API is running on environment:", env)

	brrc = suggest.InitBedrockRuntimeClientSync(context.Background())
}

func HandleRequest(
	_ context.Context,
	req events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	phrase, ok := req.QueryStringParameters["phrase"]
	if !ok {
		return apigateway.BadRequestResponse("required phrase param"), nil
	}
	sentence, ok := req.QueryStringParameters["sentence"]
	if !ok {
		return apigateway.BadRequestResponse("required sentence param"), nil
	}

	explanation, err := suggest.ExplainPhraseInSentence(brrc, phrase, sentence)
	if err != nil {
		log.Println("error when explaining: ", err)
		message := fmt.Sprintf("cannot explain \"%s\"", phrase)
		return apigateway.BadRequestResponse(message), nil
	}

	return apigateway.OkJSONResponse(explanation), nil
}

func main() {
	lambda.Start(HandleRequest)
}
