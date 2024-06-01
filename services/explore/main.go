package main

import (
	"fmt"
	"log"
	"os"

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
	"github.com/redis/go-redis/v9"
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

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisUserName := os.Getenv("REDIS_USERNAME")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Username: redisUserName,
		Password: redisPassword,
	})

	mongoURL := os.Getenv("MONGO_DATABASE_URL")
	mongoDBName := os.Getenv("MONGO_DATABASE")
	client, err := dbutils.InitMongoClient(mongoURL)
	if err != nil {
		log.Fatalln("failed to connect to mongo:", err)
	}
	db := client.Database(mongoDBName)

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
