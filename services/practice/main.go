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
	practiceapi "blinders/services/practice/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var service *practiceapi.Service

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
	adminJSON, _ := utils.GetFile("firebase.admin.json")
	url := os.Getenv("MONGO_DATABASE_URL")
	dbName := os.Getenv("MONGO_DATABASE")

	mongoManager := db.NewMongoManager(url, dbName)
	authManager, _ := auth.NewFirebaseManager(adminJSON)
	service = practiceapi.NewService(
		app,
		authManager,
		mongoManager,
		transport.NewLocalTransport(),
		transport.ConsumerMap{
			transport.Suggest: fmt.Sprintf(
				"http://localhost:%s/",
				os.Getenv("PYSUGGEST_SERVICE_PORT"),
			), // python suggest service
			transport.Collecting: fmt.Sprintf(
				"http://localhost:%s/",
				os.Getenv("COLLECTING_SERVICE_PORT"),
			),
		})
	service.App.Use(cors.New())
	service.InitRoute()
}

func main() {
	port := os.Getenv("PRACTICE_SERVICE_PORT")
	err := service.App.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Println("launch practice service error", err)
	}
}
