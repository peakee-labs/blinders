package main

import (
	"context"
	"log"
	"os"
	"sync"

	"blinders/packages/auth"
	"blinders/packages/db/chatdb"
	"blinders/packages/db/matchingdb"
	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"
	"blinders/packages/transport"
	"blinders/packages/utils"
	restapi "blinders/services/rest/api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	fiberLambda *fiberadapter.FiberLambda
	err         error
)

func init() {
	log.Println("rest api running on environment:", os.Getenv("ENVIRONMENT"))

	var usersDB *mongo.Database
	var chatDB *mongo.Database
	var matchingDB *mongo.Database
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		usersDB, err = dbutils.InitMongoDatabaseFromEnv("USERS")
		if err != nil {
			log.Fatal("failed to init users db", err)
		}
	}()
	go func() {
		defer wg.Done()
		chatDB, err = dbutils.InitMongoDatabaseFromEnv("CHAT")
		if err != nil {
			log.Fatal("failed to init chat db", err)
		}
	}()
	go func() {
		defer wg.Done()
		matchingDB, err = dbutils.InitMongoDatabaseFromEnv("MATCHING")
		if err != nil {
			log.Fatal("failed to init matching db", err)
		}
	}()
	wg.Wait()

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

	app := fiber.New()
	api := restapi.NewManager(
		app, authManager,
		usersdb.NewUsersDB(usersDB),
		chatdb.NewChatDB(chatDB),
		matchingdb.NewMatchingRepo(matchingDB),
		transport.NewLambdaTransport(cfg),
		transport.ConsumerMap{
			transport.Notification: os.Getenv("NOTIFICATION_FUNCTION_NAME"),
			transport.Explore:      os.Getenv("EXPLORE_FUNCTION_NAME"),
		},
	)

	api.App.Use(logger.New(logger.Config{Format: utils.DefaultGinLoggerFormat}))
	api.App.Use(cors.New(cors.Config{
		AllowOrigins: utils.GetOriginsFromEnv(),
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "*",
	}))

	err = api.InitRoute()
	if err != nil {
		panic(err)
	}

	fiberLambda = fiberadapter.New(api.App)
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
