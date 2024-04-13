package main

import (
	"fmt"
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/suggest"
	"blinders/packages/utils"
	suggestapi "blinders/services/suggest/api"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var service suggestapi.Service

func init() {
	if err := godotenv.Load(".env.development"); err != nil {
		log.Fatal("failed to load env", err)
	}
	app := fiber.New()
	adminJSON, _ := utils.GetFile("firebase.admin.development.json")
	url := fmt.Sprintf(
		db.MongoURLTemplate,
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_DATABASE"),
	)

	mongoManager := db.NewMongoManager(url, os.Getenv("MONGO_DATABASE"))
	fmt.Println("Connect to mongo url", url)

	authManager, _ := auth.NewFirebaseManager(adminJSON)

	openaiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(openaiKey)
	suggester, err := suggest.NewGPTSuggester(client)
	if err != nil {
		log.Fatal("failed to init openai client", err)
	}

	service = suggestapi.Service{App: app, Auth: authManager, Suggester: suggester, Db: mongoManager}
	service.InitRoute()
}

func main() {
	port := os.Getenv("SUGGEST_SERVICE_PORT")
	err := service.App.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Println("launch suggest service error", err)
	}
}
