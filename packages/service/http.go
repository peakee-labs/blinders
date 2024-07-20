package service

import (
	"log"

	"blinders/packages/auth"
	"blinders/packages/dbutils"
	"blinders/services/users/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

func HTTPServerCommonSetup() (*auth.Manager, *mongo.Database) {
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
