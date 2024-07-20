/*
Run all services as a monolith server
*/
package main

import (
	"fmt"
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/dbutils"
	"blinders/services/chat"
	"blinders/services/practice"
	"blinders/services/users"
	"blinders/services/users/repo"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	db *mongo.Database
	am *auth.Manager
)

func init() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "default"
	}

	log.Println("Running on:", env)

	envFile := ".env"
	firebaseFile := "firebase.admin.json"
	if env != "default" {
		envFile = fmt.Sprintf(".env.%s", env)
		firebaseFile = fmt.Sprintf("firebase.admin.%v.json", env)
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("failed to load env:", err)
	}

	db, err = dbutils.InitMongoDatabaseFromEnv()
	if err != nil {
		log.Fatal("failed to connect to mongo:", err)
	}

	usersRepo := repo.NewUsersRepo(db)
	am, err = auth.NewFirebaseManagerFromFile(firebaseFile, usersRepo)
	if err != nil {
		log.Fatal("failed to init auth manager:", err)
	}
}

func main() {
	services := []Service{
		{PathPrefix: "chat", Fiber: chat.NewService(am, db)},
		{PathPrefix: "users", Fiber: users.NewService(am, db)},
		{PathPrefix: "practice", Fiber: practice.NewService(am, db)},
	}

	fiberApp := fiber.New()

	for _, service := range services {
		router := fiberApp.Group(service.PathPrefix)
		service.Fiber.InitFiberRoutes(router)
	}

	port := os.Getenv("MONOLITHIC_SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("listening on: ", port)

	err := fiberApp.Listen(":" + port)
	if err != nil {
		log.Panic(err)
	}
}
