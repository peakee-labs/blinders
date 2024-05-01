package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/db/chatdb"
	"blinders/packages/db/matchingdb"
	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"
	"blinders/packages/transport"
	"blinders/packages/utils"
	restapi "blinders/services/rest/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

var apiManager restapi.Manager

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

	log.Println("rest api running on environment:", os.Getenv("ENVIRONMENT"))

	app := fiber.New()

	dbName := os.Getenv("MONGO_DATABASE")
	url := os.Getenv("MONGO_DATABASE_URL")
	client, err := dbutils.InitMongoClient(url)
	if err != nil {
		log.Fatal(err)
	}
	usersDB := usersdb.NewUsersDB(client.Database(dbName))
	chatDB := chatdb.NewChatDB(client.Database(dbName))
	matchingRepo := matchingdb.NewMatchingRepo(client.Database(dbName))

	adminJSON, _ := utils.GetFile("firebase.admin.json")
	auth, _ := auth.NewFirebaseManager(adminJSON)

	transporter := transport.NewLocalTransport()
	consumerMap := transport.ConsumerMap{
		transport.Notification: "notification_service_id",
		transport.Explore:      "explore_service_id",
	}

	apiManager = *restapi.NewManager(
		app,
		auth,
		usersDB,
		chatDB,
		matchingRepo,
		transporter,
		consumerMap,
	)

	apiManager.App.Use(logger.New())
	_ = apiManager.InitRoute()
}

func main() {
	port := os.Getenv("REST_API_PORT")
	err := apiManager.App.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Println("launch chat service error", err)
	}
}
