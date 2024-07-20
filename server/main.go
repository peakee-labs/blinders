/*
Run all services as a monolith server
*/
package main

import (
	"fmt"
	"log"
	"os"

	"blinders/packages/dbutils"
	"blinders/services/chat"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	env := os.Getenv("ENVIRONMENT")
	log.Println("Running on:", env)

	envFile := ".env"
	if env != "" {
		envFile = fmt.Sprintf(".env.%s", env)
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("failed to load env", err)
	}
}

func main() {
	db, err := dbutils.InitMongoDatabaseFromEnv()
	if err != nil {
		log.Fatal("failed to connect to mongo:", err)
	}

	services := []Service{{
		PathPrefix: "chat",
		Fiber:      chat.NewService(db),
	}}

	fiberApp := &fiber.App{}

	for _, service := range services {
		router := fiberApp.Group(service.PathPrefix)
		service.Fiber.InitFiberRoutes(router)
	}

	port := os.Getenv("MONOLITHIC_SERVER_PORT")
	fmt.Println("listening on: ", port)
	err = fiberApp.Listen(":" + port)
	if err != nil {
		log.Panic(err)
	}
}
