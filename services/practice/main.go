package main

import (
	"fmt"
	"log"
	"os"

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
	log.Println("practice api running on environment:", env)
	envFile := ".env"
	if env != "" {
		envFile = fmt.Sprintf(".env.%s", env)
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("failed to load env", err)
	}

	db, err := dbutils.InitMongoDatabaseFromEnv()
	if err != nil {
		log.Fatal("failed to connect to mongo:", err)
	}

	usersRepo := usersdb.NewUsersRepo(db)
	snapshotRepo := practicedb.NewSnapshotsRepo(db)

	adminJSON, _ := utils.GetFile("firebase.admin.json")
	auth, _ := auth.NewFirebaseManager(adminJSON)
	flashcardsRepo := practicedb.NewFlashcardsRepo(db)
	transportConsumers := transport.ConsumerMap{
		transport.Suggest: fmt.Sprintf(
			"http://localhost:%s/",
			os.Getenv("PYSUGGEST_SERVICE_PORT"),
		), // python suggest service
		transport.CollectingPush: fmt.Sprintf(
			"http://localhost:%s/",
			os.Getenv("COLLECTING_SERVICE_PORT"),
		),
		transport.CollectingGet: fmt.Sprintf(
			"http://localhost:%s/",
			os.Getenv("COLLECTING_SERVICE_PORT"),
		),
	}
	transport := transport.NewLocalTransportWithConsumers(transportConsumers)
	app := fiber.New()
	service = practiceapi.NewService(
		app,
		auth,
		usersRepo,
		flashcardsRepo,
		snapshotRepo,
		transport,
	)

	service.App.Use(cors.New())
	service.InitRoute()
}

func main() {
	port := os.Getenv("PRACTICE_SERVICE_PORT")
	log.Println("listening on: ", port)
	err := service.App.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Println("launch practice service error", err)
	}
}
