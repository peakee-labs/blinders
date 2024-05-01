package main

import (
	"context"
	"log"
	"os"

	"blinders/packages/db/collectingdb"
	dbutils "blinders/packages/db/utils"
	"blinders/packages/transport"
	collecting "blinders/services/collecting/core"

	"github.com/aws/aws-lambda-go/lambda"
)

var service *collecting.Service

func init() {
	env := os.Getenv("ENVIRONMENT")
	log.Println("collecting api running on environment:", env)

	mongoInfo := dbutils.GetMongoInfoFromEnv()
	client, err := dbutils.InitMongoClient(mongoInfo.URL)
	if err != nil {
		log.Fatal(err)
	}

	collectingDB := collectingdb.NewCollectingDB(client.Database(mongoInfo.DBName))
	service = collecting.NewService(collectingDB.ExplainLogsRepo, collectingDB.TranslateLogsRepo)
}

func LambdaHandler(_ context.Context, event transport.Event) error {
	return service.HandlePushEvent(event)
}

func main() {
	lambda.Start(LambdaHandler)
}
