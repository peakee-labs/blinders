package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"blinders/packages/auth"
	"blinders/packages/db/matchingdb"
	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"
	"blinders/packages/explore"
	"blinders/packages/transport"
	"blinders/packages/utils"
	exploreapi "blinders/services/explore/api"

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
	api         *exploreapi.Manager
	fiberLambda *fiberadapter.FiberLambda
	err         error
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	redisClient := utils.NewRedisClientFromEnv(ctx)

	var usersDB *mongo.Database
	var matchingDB *mongo.Database
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		usersDB, err = dbutils.InitMongoDatabaseFromEnv("USERS")
		if err != nil {
			log.Fatal("failed to init users db:", err)
		}
	}()
	go func() {
		defer wg.Done()
		matchingDB, err = dbutils.InitMongoDatabaseFromEnv("MATCHING")
		if err != nil {
			log.Fatal("failed to init matching db:", err)
		}
	}()
	wg.Wait()

	matchingRepo := matchingdb.NewMatchingRepo(matchingDB)
	usersRepo := usersdb.NewUsersRepo(usersDB)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Panicln("failed to load aws config:", err)
	}
	consumerMap := transport.ConsumerMap{
		transport.Embed: os.Getenv("EMBEDDER_FUNCTION_NAME"),
	}
	transporter := transport.NewLambdaTransportWithConsumers(cfg, consumerMap)

	core := explore.NewExplorer(matchingRepo, usersRepo, redisClient)
	service := exploreapi.NewService(core, redisClient, transporter)

	adminJSON, _ := utils.GetFile("firebase.admin.json")
	auth, err := auth.NewFirebaseManager(adminJSON)
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	api = exploreapi.NewManager(app, auth, usersRepo, service)
	api.App.Use(logger.New(logger.Config{Format: utils.DefaultGinLoggerFormat}))
	api.App.Use(cors.New(cors.Config{
		AllowOrigins: utils.GetOriginsFromEnv(),
		AllowMethods: "GET,OPTIONS",
		AllowHeaders: "*",
	}))
	api.InitRoute()

	fiberLambda = fiberadapter.New(api.App)
}

func HandleRequest(ctx context.Context, req any) (any, error) {
	internalReq, err := utils.JSONConvert[transport.Request](req)
	if err != nil {
		log.Fatal("can not parse http proxy request:", err)
	}

	switch internalReq.Type {
	case transport.AddUserMatchInfo:
		req, err := utils.JSONConvert[transport.AddUserMatchInfoRequest](req)
		if err != nil {
			log.Println("can't parse match info from request: ", err)
			return nil, fmt.Errorf("can not parse match info: %v", err)
		}

		err = api.Service.AddUserMatch(&req.Payload)
		if err != nil {
			return nil, fmt.Errorf("can not add user match: %v", err)
		}

		return nil, nil

	default:
		req, err := utils.JSONConvert[events.APIGatewayV2HTTPRequest](req)
		if err != nil {
			log.Println("can not parse http proxy request:", err)
			return nil, fmt.Errorf("can not parse http proxy request")
		}

		return fiberLambda.ProxyWithContextV2(ctx, *req)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
