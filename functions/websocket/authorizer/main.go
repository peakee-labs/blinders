package main

import (
	"context"
	"encoding/json"
	"log"

	"blinders/packages/auth"
	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"
	"blinders/packages/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type APIGatewayWebsocketProxyRequest struct {
	events.APIGatewayWebsocketProxyRequest `       json:",inline"`
	MethodArn                              string `json:"methodArn"` // ??? refs: https://gist.github.com/praveen001/1b045d1c31cd9c72e4e6638e9f883f83
}

var (
	authManager auth.Manager
	userRepo    *usersdb.UsersRepo
)

func init() {
	mongoInfo := dbutils.GetMongoInfoFromEnv()
	client, err := dbutils.InitMongoClient(mongoInfo.URL)
	if err != nil {
		log.Fatal(err)
	}
	userRepo = usersdb.NewUsersRepo(client.Database(mongoInfo.DBName))

	adminConfig, err := utils.GetFile("firebase.admin.json")
	if err != nil {
		log.Fatal(err)
	}

	authManager, err = auth.NewFirebaseManager(adminConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func HandleRequest(
	_ context.Context,
	request APIGatewayWebsocketProxyRequest,
) (events.APIGatewayCustomAuthorizerResponse, error) {
	jwt := request.QueryStringParameters["token"]
	authUser, err := authManager.Verify(jwt)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}

	user, err := userRepo.GetUserByFirebaseUID(authUser.AuthID)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}

	// Is it secure to log the id out to cloudwatch?
	// how to log the request tracking efficient and secure
	log.Println("[authorizer] issued user's policy of", user.ID.Hex())

	userBytes, _ := json.Marshal(authUser)
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: user.ID.Hex(),
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Allow",
					Resource: []string{request.MethodArn},
				},
			},
		},
		Context: map[string]interface{}{
			"user": string(userBytes),
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
