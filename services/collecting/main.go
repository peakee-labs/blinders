package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"blinders/packages/db/collectingdb"
	core "blinders/services/collecting/core"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mongoURL := os.Getenv("MONGO_DATABASE_URL")
	mongoDBName := os.Getenv("MONGO_DATABASE")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatalln("failed to connect to mongo:", err)
	}
	db := client.Database(mongoDBName)

	m = core.NewManager(
		fiber.New(),
		collectingdb.NewCollectingDB(db),
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
