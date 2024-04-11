package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/transport"
	"blinders/packages/utils"
	suggestapi "blinders/services/suggest/api"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var service *suggestapi.Service

func init() {
	env := os.Getenv("ENVIRONMENT")
	envFile := ".env"
	if env != "" {
		envFile = ".env." + strings.ToLower(env)
	}
	log.Println("init service in environment", env, "loading env at", envFile)
	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("failed to load env", err)
	}

	app := fiber.New()
	adminJSON, _ := utils.GetFile("firebase.admin.development.json")
	url := os.Getenv("MONGO_DATABASE_URL")
	dbName := os.Getenv("MONGO_DATABASE")

	mongoManager := db.NewMongoManager(url, dbName)
	authManager, _ := auth.NewFirebaseManager(adminJSON)
	service = suggestapi.NewService(app, authManager, mongoManager, transport.NewLocalTransport(), transport.ConsumerMap{
		transport.Logging: fmt.Sprintf("http://localhost:%s/", os.Getenv("LOGGING_SERVICE_PORT")),
		transport.Suggest: fmt.Sprintf("http://localhost:%s/", os.Getenv("PYSUGGEST_SERVICE_PORT")), // python suggest service
	})
	service.InitRoute()
}

func main() {
	port := os.Getenv("SUGGEST_SERVICE_PORT")
	err := service.App.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Println("launch suggest service error", err)
	}
}
