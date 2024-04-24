package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/collecting"
	"blinders/packages/db"
	"blinders/packages/transport"
	collectingapi "blinders/services/collecting/api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var service *collectingapi.Service

func init() {
	env := os.Getenv("ENVIRONMENT")
	log.Println("collecting api running on environment:", env)
	url := fmt.Sprintf(
		db.MongoURLTemplate,
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_DATABASE"),
	)
	dbName := os.Getenv("MONGO_DATABASE")
	client, err := db.InitMongoClient(url)
	if err != nil {
		log.Fatalf("cannot init mongo client, err: %v", err)
	}

	collector := collecting.NewEventCollector(client.Database(dbName))
	service = &collectingapi.Service{Collector: collector}
}

func LambdaHandler(
	ctx context.Context,
	eventRequest transport.GetEventRequest,
) (
	events.APIGatewayV2HTTPResponse,
	error,
) {
	authUser, ok := ctx.Value(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panicln("unexpected err, authUser not included in ctx")
	}
	fmt.Println("received", eventRequest)
	if eventRequest.Request.Type != transport.GetEvent {
		log.Printf("collecting: event type mismatch, type: %v\n", eventRequest.Request.Type)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       "request type mismatch",
		}, nil
	}
	userOID, _ := primitive.ObjectIDFromHex(authUser.ID)

	event, err := service.HandleGetGenericEvent(userOID, eventRequest.Type)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("failed to get event, err: %v", err),
		}, nil
	}

	payloadBytes, _ := json.Marshal(event.Payload)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(payloadBytes),
	}, nil
}

func main() {
	// TODO: auth lambda
	lambda.Start(LambdaHandler)
}
