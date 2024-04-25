package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/transport"
	"blinders/packages/utils"
	practiceapi "blinders/services/practice/api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	fiberLambda *fiberadapter.FiberLambda
)

func init() {
	log.Println("practice api running on environment:", os.Getenv("ENVIRONMENT"))
	app := fiber.New(fiber.Config{})

	url := fmt.Sprintf(
		db.MongoURLTemplate,
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_DATABASE"),
	)

	database := db.NewMongoManager(url, os.Getenv("MONGO_DATABASE"))
	if database == nil {
		log.Fatal("cannot create database manager")
	}

	adminConfig, err := utils.GetFile("firebase.admin.json")
	if err != nil {
		log.Fatal(err)
	}

	authManager, err := auth.NewFirebaseManager(adminConfig)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("failed to load aws config", err)
	}
	api := practiceapi.NewService(
		app,
		authManager,
		database,
		transport.NewLambdaTransport(cfg),
		transport.ConsumerMap{
			transport.CollectingPush: os.Getenv("COLLECTING_PUSH_FUNCTION_NAME"),
			transport.CollectingGet:  os.Getenv("COLLECTING_GET_FUNCTION_NAME"),
		},
	)
	api.App.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${queryParams} | ${error}\n",
	}))

	api.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET",
		AllowHeaders: "*",
	}))

	api.InitRoute()

	fiberLambda = fiberadapter.New(api.App)
}

func HandleRequest(
	ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	return fiberLambda.ProxyWithContextV2(ctx, req)
}

func main() {
	lambda.Start(HandleRequest)
}
