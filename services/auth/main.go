package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"blinders/packages/auth"
	dbutils "blinders/packages/dbutils"
	"blinders/packages/utils"
	usersrepo "blinders/services/users/repo"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	authManager *auth.FirebaseManager
	usersRepo   *usersrepo.UsersRepo
)

func init() {
	env := os.Getenv("ENVIRONMENT")
	log.Printf("Authentication api running on %s environment\n", env)

	usersDB, err := dbutils.InitMongoDatabaseFromEnv("USERS")
	if err != nil {
		log.Fatal(err)
	}
	usersRepo = usersrepo.NewUsersRepo(usersDB)

	adminConfig, err := utils.GetFile("firebase.admin.json")
	if err != nil {
		log.Fatal(err)
	}
	authManager, err = auth.NewFirebaseManager(adminConfig)
	if err != nil {
		log.Fatal(err)
	}
}

type authRequest struct {
	Token string `json:"token"` // bearer token
}

func handler(_ context.Context, req authRequest) (auth.UserAuth, error) {
	authToken := req.Token
	fmt.Println("received", req)
	if !strings.HasPrefix(authToken, "Bearer ") {
		log.Println("invalid jwt, missing bearer token")
		return auth.UserAuth{}, fmt.Errorf("missing bearer token")
	}

	jwt := strings.Split(authToken, " ")[1]
	userAuth, err := authManager.Verify(jwt)
	if err != nil {
		log.Println("failed to verify jwt", err)
		return auth.UserAuth{}, fmt.Errorf("failed to verify jwt, err: %v", err)
	}

	// currently, user.AuthID is firebaseUID
	user, err := usersRepo.GetUserByFirebaseUID(userAuth.AuthID)
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
