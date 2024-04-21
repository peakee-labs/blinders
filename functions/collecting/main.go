package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/collecting"
	"blinders/packages/db"
	"blinders/packages/utils"
	collectingapi "blinders/services/collecting/api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	service     *collectingapi.Service
	fiberLambda *fiberadapter.FiberLambda
)

func init() {
	log.Println("collecting api running on environment:", os.Getenv("ENVIRONMENT"))
	app := fiber.New()
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

	adminConfig, err := utils.GetFile("firebase.admin.json")
	if err != nil {
		log.Fatal(err)
	}
	authManager, err := auth.NewFirebaseManager(adminConfig)
	if err != nil {
		log.Fatal(err)
	}

	collector := collecting.NewEventCollector(client.Database(dbName))
	service = collectingapi.NewCollectingService(app, collector, authManager)
	service.App.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${queryParams} | ${error}\n",
	}))

	service.App.Use(cors.New(cors.Config{
		AllowOrigins: "https://app.peakee.co, http://localhost:3000",
		AllowMethods: "GET,POST",
		AllowHeaders: "*",
	}))

	err = service.InitRoute()
	if err != nil {
		panic(err)
	}
	fiberLambda = fiberadapter.New(service.App)
}

func ProxiHandler(
	ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (
	events.APIGatewayV2HTTPResponse,
	error,
) {
	return fiberLambda.ProxyWithContextV2(ctx, req)
}

func main() {
	lambda.Start(ProxiHandler)
}
