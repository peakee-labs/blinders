package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/explore"
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

	mongoManager := db.NewMongoManager(url, dbName)

	fmt.Println("Connected to mongo url", url)

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

	core := explore.NewMongoExplorer(mongoManager, redisClient)

	embedderEndpoint := fmt.Sprintf(
		"http://localhost:%s/embed",
		os.Getenv("EMBEDDER_SERVICE_PORT"),
	)
	service = exploreapi.NewService(core, redisClient, embedderEndpoint)

	manager = exploreapi.NewManager(app, authManager, mongoManager, service)
	manager.App.Use(logger.New(), cors.New())

	// Expose for local development
	manager.App.Post("/explore", manager.Service.HandleAddUserMatch)
	manager.InitRoute()
}

func main() {
	port := os.Getenv("EXPLORE_API_PORT")
	fmt.Println("listening on: ", port)
	log.Panic(manager.App.Listen(":" + port))
}
