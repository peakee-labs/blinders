package main

import (
	"context"
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"
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

var fiberLambda *fiberadapter.FiberLambda

func init() {
	log.Println("practice api running on environment:", os.Getenv("ENVIRONMENT"))

	usersDB, err := dbutils.InitMongoDatabaseFromEnv("USERS")
	if err != nil {
		log.Fatal(err)
	}
	usersRepo := usersdb.NewUsersRepo(usersDB)

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

	app := fiber.New(fiber.Config{})
	api := practiceapi.NewService(
		app,
		authManager,
		usersRepo,
		transport.NewLambdaTransport(cfg),
		transport.ConsumerMap{
			transport.CollectingPush: os.Getenv("COLLECTING_PUSH_FUNCTION_NAME"),
			transport.CollectingGet:  os.Getenv("COLLECTING_GET_FUNCTION_NAME"),
		},
	)
	api.App.Use(logger.New(logger.Config{Format: utils.DefaultGinLoggerFormat}))
	api.App.Use(cors.New(cors.Config{
		AllowOrigins: utils.GetOriginsFromEnv(),
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
	lambda.Start((HandleRequest))
}
