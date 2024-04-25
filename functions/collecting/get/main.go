package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"blinders/packages/collecting"
	"blinders/packages/db"
	"blinders/packages/transport"
	collectingapi "blinders/services/collecting/api"

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
	transport.GetEventResponse,
	error,
) {
	if eventRequest.Request.Type != transport.GetEvent {
		log.Printf("collecting: event type mismatch, type: %v\n", eventRequest.Request.Type)
		return transport.GetEventResponse{}, fmt.Errorf("event type mismatch")
	}
	userOID, err := primitive.ObjectIDFromHex(eventRequest.UserID)
	if err != nil {
		log.Println("cannot get object id from event", err)
		return transport.GetEventResponse{}, err
	}

	event, err := service.HandleGetGenericEvent(userOID, eventRequest.Type)
	if err != nil {
		log.Println("collecting: failed to get event", event)
		return transport.GetEventResponse{}, nil
	}
	response := transport.GetEventResponse{
		Data: []collecting.GenericEvent{
			event,
		},
	}
	return response, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
