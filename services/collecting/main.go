package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"blinders/packages/db/collectingdb"
	dbutils "blinders/packages/db/utils"
	core "blinders/services/collecting/core"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var m *core.Manager

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

	client, err := dbutils.InitMongoClient(url)
	if err != nil {
		log.Fatal(err)
	}

	m = core.NewManager(
		app,
		collectingdb.NewCollectingDB(client.Database(dbName)),
	)

	m.App.Use(cors.New())
	_ = m.InitRoute()
}

func main() {
	port := os.Getenv("COLLECTING_SERVICE_PORT")
	err := m.App.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Println("launch collecting service error", err)
	}
}
