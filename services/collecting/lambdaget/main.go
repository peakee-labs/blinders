package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"blinders/packages/db/collectingdb"
	"blinders/packages/transport"

	dbutils "blinders/packages/db/utils"

	collecting "blinders/services/collecting/core"

	"github.com/aws/aws-lambda-go/lambda"
)

var service *collecting.Service

func init() {
	env := os.Getenv("ENVIRONMENT")
	log.Println("collecting api running on environment:", env)

	mongoDB, err := dbutils.InitMongoDatabaseFromEnv("COLLECTING")
	if err != nil {
		log.Fatal(err)
	}

	collectingDB := collectingdb.NewCollectingDB(mongoDB)
	service = collecting.NewService(collectingDB.ExplainLogsRepo, collectingDB.TranslateLogsRepo)
}

func LambdaHandler(_ context.Context, request transport.Request) (any, error) {
	res, err := service.HandleGetRequest(request)
	if err != nil {
		log.Println("can not handle request:", err)
		return nil, fmt.Errorf("can not handle request")
	}

	return res, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
