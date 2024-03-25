package main

import (
	"fmt"
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/transport"
	"blinders/packages/utils"
	restapi "blinders/services/rest/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

var apiManager restapi.Manager

func init() {
	if err := godotenv.Load(".env.development"); err != nil {
		log.Fatal("failed to load env", err)
	}

	log.Println("rest api running on environment:", os.Getenv("ENVIRONMENT"))

	app := fiber.New()

	dbName := os.Getenv("MONGO_DATABASE")
	url := os.Getenv("MONGO_DATABASE_URL")
	database := db.NewMongoManager(url, dbName)
	if database == nil {
		log.Fatal("cannot create database manager")
	}

	adminJSON, _ := utils.GetFile("firebase.admin.development.json")
	authManager, _ := auth.NewFirebaseManager(adminJSON)

	apiManager = *restapi.NewManager(
		app, authManager, database,
		transport.NewLocalTransport(),
		transport.ConsumerMap{
			transport.Notification: "notification_service_id",
		})
	apiManager.App.Use(logger.New())
	_ = apiManager.InitRoute(restapi.InitOptions{})
}

func main() {
	port := os.Getenv("REST_API_PORT")
	err := apiManager.App.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Println("launch chat service error", err)
	}
}
