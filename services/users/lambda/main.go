package main

import (
	"context"
	"log"
	"os"

	"blinders/packages/service"
	"blinders/services/users"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
)

var fiberLambda *fiberadapter.FiberLambda

func init() {
	env := os.Getenv("ENVIRONMENT")
	log.Println("Peakee Users API is running on environment:", env)

	auth, mongoDB := service.LambdaCommonSetup()
	usersService := users.NewService(auth, mongoDB)

	app := fiber.New()
	usersService.InitFiberRoutes(app.Group("/users"))

	fiberLambda = service.NewFiberLambdaAdapter(app)
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
