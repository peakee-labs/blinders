package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/explore"
	"blinders/packages/transport"
	"blinders/packages/utils"
	exploreapi "blinders/services/explore/api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
)

var (
	api         *exploreapi.Manager
	fiberLambda *fiberadapter.FiberLambda
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

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
		panic("cannot create mongo manager")
	}
	log.Println("database connected")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	core := explore.NewMongoExplorer(database, redisClient)
	service := exploreapi.NewService(core, redisClient, os.Getenv("EMBEDDER_ENDPOINT"))
	app := fiber.New()

	adminJSON, _ := utils.GetFile("firebase.admin.json")
	authManager, err := auth.NewFirebaseManager(adminJSON)
	if err != nil {
		panic(err)
	}

	api = exploreapi.NewManager(app, authManager, database, service)

	api.App.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${queryParams} | ${error}\n",
	}))

	api.App.Use(cors.New(cors.Config{
		AllowOrigins: "https://app.peakee.co, http://localhost:3000",
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
