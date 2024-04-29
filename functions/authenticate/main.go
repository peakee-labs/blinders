package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/db/repo"
	"blinders/packages/utils"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	authManager *auth.FirebaseManager
	userRepo    *repo.UsersRepo
)

func init() {
	env := os.Getenv("ENVIRONMENT")
	log.Printf("Authentication api running on %s environment\n", env)

	url := fmt.Sprintf(
		db.MongoURLTemplate,
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_DATABASE"),
	)

	database := db.NewMongoManager(url, os.Getenv("MONGO_DATABASE"))
	if database == nil {
		log.Fatal("cannot create database manager")
	}
	adminConfig, err := utils.GetFile("firebase.admin.json")
	if err != nil {
		log.Fatal(err)
	}
	authManager, err = auth.NewFirebaseManager(adminConfig)
	if err != nil {
		log.Fatal(err)
	}
	userRepo = database.Users
}

type authRequest struct {
	Token string `json:"token"` // bearer token
}

func handler(
	ctx context.Context,
	event authRequest,
) (
	auth.UserAuth, error,
) {
	authToken := event.Token
	fmt.Println("received", event)
	if !strings.HasPrefix(authToken, "Bearer ") {
		log.Println("invalid jwt, missing bearer token")
		return auth.UserAuth{}, fmt.Errorf("missing bearer token")
	}

	jwt := strings.Split(authToken, " ")[1]
	userAuth, err := authManager.Verify(jwt)
	if err != nil {
		log.Println("failed to verify jwt", err)
		return auth.UserAuth{}, fmt.Errorf("failed to verify jwt")
	}

	// currently, user.AuthID is firebaseUID
	user, err := userRepo.GetUserByFirebaseUID(userAuth.AuthID)
	if err != nil {
		log.Println("failed to get user", err)
		return auth.UserAuth{}, fmt.Errorf("failed to get user")
	}

	userAuth.ID = user.ID.Hex()
	return *userAuth, nil
}

func main() {
	lambda.Start(handler)
}
