package main

import (
	"context"
	"log"

	"blinders/packages/session"
	"blinders/packages/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var sessionManager *session.Manager

func init() {
	redisClient := utils.NewRedisClientFromEnv(context.Background())
	sessionManager = session.NewManager(redisClient)
}

func HandleRequest(
	_ context.Context,
	request events.APIGatewayWebsocketProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	connectionID := request.RequestContext.ConnectionID
	userID := request.RequestContext.Authorizer.(map[string]interface{})["principalId"].(string)

	if userID == "" {
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "user not found"}, nil
	}

	err := sessionManager.AddSession(userID, connectionID)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "failed to add session"}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "connected"}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
