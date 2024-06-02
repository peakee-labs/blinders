package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"blinders/packages/auth"
	"blinders/packages/db/matchingdb"
	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"
	"blinders/packages/explore"
	"blinders/packages/transport"
	"blinders/packages/utils"
	exploreapi "blinders/services/explore/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var manager *exploreapi.Manager

func init() {
	env := os.Getenv("ENVIRONMENT")
	envFile := ".env"
	if env != "" {
		envFile = fmt.Sprintf(".env.%s", env)
	}
	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("failed to load env", err)
	}
	log.Println("explore api running on environment:", env)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	redisClient := utils.NewRedisClientFromEnv(ctx)

	db, err := dbutils.InitMongoDatabaseFromEnv()
	if err != nil {
		log.Fatal("failed to connect to mongo:", err)
	}

	matchingRepo := matchingdb.NewMatchingRepo(db)
	usersRepo := usersdb.NewUsersRepo(db)

	embedderEndpoint := fmt.Sprintf("http://localhost:%s/embedd", os.Getenv("EMBEDDER_SERVICE_PORT"))
	fmt.Println("embedder endpoint: ", embedderEndpoint)

	tp := transport.NewLocalTransportWithConsumers(
		transport.ConsumerMap{
			transport.Embed: embedderEndpoint,
		},
	)

	core := explore.NewExplorer(matchingRepo, usersRepo, redisClient)
	service := exploreapi.NewService(core, redisClient, tp)

	adminJSON, _ := utils.GetFile("firebase.admin.json")
	auth, err := auth.NewFirebaseManager(adminJSON)
	if err != nil {
		panic(err)
	}
	app := fiber.New()

	manager = exploreapi.NewManager(app, auth, usersRepo, service)
	manager.InitRoute()

	// Expose for local development
	manager.App.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost",
		AllowMethods: "POST,GET,OPTIONS,PUT,DELETE",
	}))

	manager.App.All("/explore", manager.Service.InternalHandleAddUserMatch)
}

func main() {
	port := os.Getenv("EXPLORE_SERVICE_PORT")
	fmt.Println("listening on: ", port)
	log.Panic(manager.App.Listen(":" + port))
}
