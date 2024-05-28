package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"blinders/packages/auth"
	"blinders/packages/db/matchingdb"
	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"
	"blinders/packages/explore"
	"blinders/packages/transport"
	"blinders/packages/utils"
	exploreapi "blinders/services/explore/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	service *exploreapi.Service
	manager *exploreapi.Manager
)

func init() {
	environment := os.Getenv("ENVIRONMENT")
	log.Println("explore api running on environment:", environment)
	envFile := ".env"
	if environment != "" {
		envFile = fmt.Sprintf(".env.%s", environment)
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("failed to load env", err)
	}
	app := fiber.New(fiber.Config{
		WriteTimeout:          time.Second * 5,
		ReadTimeout:           time.Second * 5,
		DisableStartupMessage: true,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	dbName := os.Getenv("MONGO_DATABASE")
	url := os.Getenv("MONGO_DATABASE_URL")

	client, err := dbutils.InitMongoClient(url)
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	adminJSON, _ := utils.GetFile("firebase.admin.json")
	authManager, err := auth.NewFirebaseManager(adminJSON)
	if err != nil {
		panic(err)
	}

	// matchingDB := matchingdb.
	usersRepo := usersdb.NewUsersRepo(client.Database(dbName))
	matchingRepo := matchingdb.NewMatchingRepo(client.Database(dbName))
	core := explore.NewExplorer(
		matchingRepo,
		usersRepo,
		redisClient,
	)

	embedderEndpoint := fmt.Sprintf("http://localhost:%s/embedd", os.Getenv("EMBEDDER_SERVICE_PORT"))
	fmt.Println("embedder endpoint: ", embedderEndpoint)

	tp := transport.NewLocalTransportWithConsumers(
		transport.ConsumerMap{
			transport.Embed: embedderEndpoint,
		},
	)

	service = exploreapi.NewService(core, redisClient, tp)

	manager = exploreapi.NewManager(app, authManager, usersRepo, service)

	manager.App.Use(logger.New(), cors.New())
	manager.InitRoute()

	// Expose for local development
	manager.App.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost",
		AllowMethods: "POST,GET,OPTIONS,PUT,DELETE",
	}))
	manager.App.All("/explore", manager.Service.InternalHandleAddUserMatch)
}

func main() {
	port := os.Getenv("EXPLORE_SERVICE_PORT")
	fmt.Println("listening on: ", port)
	log.Panic(manager.App.Listen(":" + port))
}
