package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"blinders/packages/auth"
	"blinders/packages/db/matchingdb"
	"blinders/packages/db/usersdb"
	"blinders/packages/explore"
	"blinders/packages/transport"
	"blinders/packages/utils"

	dbutils "blinders/packages/db/utils"

	exploreapi "blinders/services/explore/api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
			log.Fatal("failed to init users db", err)
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

	matchingRepo := matchingdb.NewMatchingRepo(matchingDB)
	usersRepo := usersdb.NewUsersRepo(usersDB)

	core := explore.NewExplorer(matchingRepo, usersRepo, redisClient)
	service := exploreapi.NewService(core, redisClient, os.Getenv("EMBEDDER_ENDPOINT"))

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

func HandleRequest(ctx context.Context, payload any) (any, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("can not marshal payload:", err)
		return nil, fmt.Errorf("can not marshal payload")
	}

	request, err := utils.ParseJSON[transport.AddUserMatchInfoRequest](bytes)
	if err != nil || request.Type != transport.AddUserMatchInfo {
		log.Println("might be http request from client app", err)
		req, err := utils.ParseJSON[events.APIGatewayV2HTTPRequest](bytes)
		if err != nil {
			log.Println("can not parse http proxy request:", err)
			return nil, fmt.Errorf("can not parse http proxy request")
		}

		return fiberLambda.ProxyWithContextV2(ctx, *req)
	} else if request == nil {
		return nil, fmt.Errorf("nil request")
	}

	err = api.Service.AddUserMatch(request.Data)
	if err != nil {
		return nil, fmt.Errorf("can not add user match: %v", err)
	}

	return nil, nil
}

func main() {
	lambda.Start(HandleRequest)
}
