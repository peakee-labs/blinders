package main

import (
	"context"
	"log"
	"os"

	dbutils "blinders/packages/dbutils"
	"blinders/packages/utils"
	"blinders/services/chat"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var fiberLambda *fiberadapter.FiberLambda

func init() {
	log.Println("rest api running on environment:", os.Getenv("ENVIRONMENT"))

	app := fiber.New()

	chatDB, err := dbutils.InitMongoDatabaseFromEnv("CHAT")
	if err != nil {
		log.Fatal("failed to init chat db:", err)
	}
	chatService := chat.NewService(chatDB)
	chatService.InitFiberRoutes(app.Group("/chat"))

	app.Use(logger.New(logger.Config{Format: utils.DefaultFiberLoggerFormat}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: utils.GetOriginsFromEnv(),
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "*",
	}))

	fiberLambda = fiberadapter.New(app)
}

func Handler(
	ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	return fiberLambda.ProxyWithContextV2(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
