package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
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
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	manager *exploreapi.Manager
	err     error
)

func init() {
	environment := os.Getenv("ENVIRONMENT")
	log.Println("explore api running on environment:", environment)
	envFile := ".env"
	if environment != "" {
		envFile = fmt.Sprintf(".env.%s", environment)
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("failed to load env", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	redisClient := utils.NewRedisClientFromEnv(ctx)

	var usersDB *mongo.Database
	var matchingDB *mongo.Database
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		usersDB, err = dbutils.InitMongoDatabaseFromEnv("USERS")
		if err != nil {
			log.Fatal("failed to init users db:", err)
		}
	}()
	go func() {
		defer wg.Done()
		matchingDB, err = dbutils.InitMongoDatabaseFromEnv("MATCHING")
		if err != nil {
			log.Fatal("failed to init matching db:", err)
		}
	}()
	wg.Wait()

	matchingRepo := matchingdb.NewMatchingRepo(matchingDB)
	usersRepo := usersdb.NewUsersRepo(usersDB)

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
