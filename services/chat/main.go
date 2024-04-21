package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/utils"
	chatapi "blinders/services/chat/api"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var service chatapi.Service

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
	adminJSON, _ := utils.GetFile("firebase.admin.json")

	authManager, _ := auth.NewFirebaseManager(adminJSON)
	service = chatapi.Service{App: app, Auth: authManager}
	service.InitRoute()
}

func main() {
	port := os.Getenv("CHAT_SERVICE_PORT")
	err := service.App.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Println("launch chat service error", err)
	}
}
