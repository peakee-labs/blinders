package main

import (
	"context"
	"log"
	"os"

	dbutils "blinders/packages/dbutils"
	"blinders/packages/utils"
	"blinders/services/practice"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var fiberLambda *fiberadapter.FiberLambda

func init() {
	log.Println("Peakee Practice API is running on environment:", os.Getenv("ENVIRONMENT"))

	fiberApp := fiber.New(fiber.Config{})

	practiceDB, err := dbutils.InitMongoDatabaseFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	practiceService := practice.NewService(practiceDB)
	practiceService.InitFiberRoutes(fiberApp.Group("practice"))

	fiberApp.Use(logger.New(logger.Config{Format: utils.DefaultFiberLoggerFormat}))
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: utils.GetOriginsFromEnv(),
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "*",
	}))

	fiberLambda = fiberadapter.New(fiberApp)
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
