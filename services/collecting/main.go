package main

import (
	"fmt"
	"log"
	"os"

	"blinders/packages/db/collectingdb"
	dbutils "blinders/packages/db/utils"
	core "blinders/services/collecting/core"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var (
	m   *core.Manager
	err error
)

func init() {
	environment := os.Getenv("ENVIRONMENT")
	log.Println("collecting api running on environment:", environment)
	envFile := ".env"
	if environment != "" {
		envFile = fmt.Sprintf(".env.%s", environment)
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("failed to load env", err)
	}

	collectingDB, err := dbutils.InitMongoDatabaseFromEnv("COLLECTING")
	if err != nil {
		log.Fatal("failed to init practice db:", err)
	}

	m = core.NewManager(
		fiber.New(),
		collectingdb.NewCollectingDB(collectingDB),
	)

	m.App.Use(cors.New())
	_ = m.InitRoute()
}

func main() {
	port := os.Getenv("COLLECTING_SERVICE_PORT")
	fmt.Println("launching collecting service on port", port)
	err := m.App.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Println("launch collecting service error", err)
	}
}
