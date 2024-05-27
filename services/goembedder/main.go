package main

import (
	"blinders/services/embedder/api"
	"blinders/services/embedder/core"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var service *api.Service

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

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("ap-south-1"),
		config.WithSharedConfigProfile("admin.peakee"),
	)

	brrc, err := core.InitBedrockRuntimeClientFromConfig(cfg)
	if err != nil {
		log.Fatal("failed to load bedrock runtime client", err)
	}
	embedder := core.NewEmbbeder(brrc, "cohere.embed-english-v3")

	service = api.NewService(app, embedder)

	service.App.Use(cors.New(
		cors.Config{
			// cors origin to allow localhost only
			AllowOrigins: "*",
			AllowMethods: "GET",
		},
	))
	service.InitRoute()
}

func main() {
	port := os.Getenv("EMBEDDER_SERVICE_PORT")
	err := service.App.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Println("launch embedder service error", err)
	}
}
