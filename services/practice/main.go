package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/db/practicedb"
	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"
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

	dbName := os.Getenv("MONGO_DATABASE")
	url := os.Getenv("MONGO_DATABASE_URL")
	client, err := dbutils.InitMongoClient(url)
	if err != nil {
		log.Fatal(err)
	}
	usersRepo := usersdb.NewUsersRepo(client.Database(dbName))

	adminJSON, _ := utils.GetFile("firebase.admin.json")
	authManager, _ := auth.NewFirebaseManager(adminJSON)

	flashcardsRepo := practicedb.NewFlashCardRepo(client.Database(dbName))
	service = practiceapi.NewService(
		app,
		authManager,
		usersRepo,
		flashcardsRepo,
		transport.NewLocalTransport(),
		transport.ConsumerMap{
			transport.Suggest: fmt.Sprintf(
				"http://localhost:%s/",
				os.Getenv("PYSUGGEST_SERVICE_PORT"),
			), // python suggest service
			transport.CollectingPush: fmt.Sprintf(
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
