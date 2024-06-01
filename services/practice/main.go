package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"blinders/packages/auth"
	"blinders/packages/db/practicedb"
	"blinders/packages/db/usersdb"
	"blinders/packages/transport"
	"blinders/packages/utils"
	practiceapi "blinders/services/practice/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	service *practiceapi.Service
	err     error
)

func init() {
	environment := os.Getenv("ENVIRONMENT")
	log.Println("practice api running on environment:", environment)
	envFile := ".env"
	if environment != "" {
		envFile = fmt.Sprintf(".env.%s", environment)
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("failed to load env", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mongoURL := os.Getenv("MONGO_DATABASE_URL")
	mongoDBName := os.Getenv("MONGO_DATABASE")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatalln("failed to connect to mongo:", err)
	}
	db := client.Database(mongoDBName)

	usersRepo := usersdb.NewUsersRepo(db)

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
