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

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	service *collectingapi.Service
)

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

func ProxiHandler(
	ctx context.Context,
	eventRequest transport.CollectEventRequest,
) (
	events.APIGatewayV2HTTPResponse,
	error,
) {
	fmt.Println("received", eventRequest)
	if eventRequest.Request.Type != transport.CollectEvent {
		log.Printf("collecting: event type mismatch, type: %v\n", eventRequest.Request.Type)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       "request type mismatch",
		}, nil
	}

	eventID, err := service.HandleGenericEvent(eventRequest.Data)
	if err != nil {
		log.Printf("collecting: cannot process generic event, err: %v\n", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("failed to process event, err: %v", err),
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       eventID,
	}, nil
}

func main() {
	lambda.Start(ProxiHandler)
}
