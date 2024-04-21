package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/collecting"
	"blinders/packages/db"
	"blinders/packages/utils"
	collectingapi "blinders/services/collecting/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var service collectingapi.Service

func init() {
	env := os.Getenv("ENVIRONMENT")
	envFile := ".env"
	if env != "" {
		envFile = strings.ToLower(".env." + env)
	}
	log.Println("init service in environment", env, "loading env at", envFile)
	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("failed to load env", err)
	}

	app := fiber.New()
	url := os.Getenv("MONGO_DATABASE_URL")
	dbName := os.Getenv("MONGO_DATABASE")

	mongoManager := db.NewMongoManager(url, dbName)

	adminConfig, err := utils.GetFile("firebase.admin.json")
	if err != nil {
		log.Fatal(err)
	}

	authManager, err := auth.NewFirebaseManager(adminConfig)
	if err != nil {
		log.Fatal(err)
	}
	service = *collectingapi.NewCollectingService(
		app,
		collecting.NewEventCollector(mongoManager.Client.Database(dbName)),
		authManager)
	service.App.Use(cors.New())
	_ = service.InitRoute()
}

func main() {
	fmt.Println("hello world from collecting service")
}
