package service

import (
	"log"

	"blinders/packages/auth"
	"blinders/packages/dbutils"
	"blinders/packages/utils"
	"blinders/services/users/repo"

	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

func LambdaCommonSetup() (*auth.Manager, *mongo.Database) {
	mongoDB, err := dbutils.InitMongoDatabaseFromEnv()
	if err != nil {
		log.Fatal("failed to init database:", err)
	}

	usersRepo := repo.NewUsersRepo(mongoDB)
	auth, err := auth.NewFirebaseManagerFromFile("firebase.admin.json", usersRepo)
	if err != nil {
		log.Fatal("failed to init auth manager:", err)
	}

	return auth, mongoDB
}

func NewFiberLambdaAdapter(app *fiber.App) *fiberadapter.FiberLambda {
	app.Use(logger.New(logger.Config{Format: utils.DefaultFiberLoggerFormat}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: utils.GetOriginsFromEnv(),
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "*",
	}))

	return fiberadapter.New(app)
}
