package apigateway

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func BadRequestResponse(message string) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusBadRequest,
		Body:       message,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
	}
}

func OkJSONResponse(body any) events.APIGatewayV2HTTPResponse {
	resInBytes, _ := json.Marshal(body)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(resInBytes),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Content-Type":                "application/json",
		},
	}
}
