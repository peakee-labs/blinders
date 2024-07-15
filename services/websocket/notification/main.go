/*
This service is responsible for notifying any event to users via websocket or push notification
*/
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"

	"blinders/packages/apigateway"
	"blinders/packages/session"
	"blinders/packages/transport"
	"blinders/packages/utils"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	APIGatewayClient *apigateway.Client
	SessionManager   *session.Manager
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("failed to load aws config", err)
	}
	cer := apigateway.CustomEndpointResolve{
		Domain:     os.Getenv("API_GATEWAY_DOMAIN"),
		PathPrefix: os.Getenv("API_GATEWAY_PATH_PREFIX"),
	}
	APIGatewayClient = apigateway.NewClient(context.Background(), cfg, cer)

	redisClient := utils.NewRedisClientFromEnv(context.Background())
	SessionManager = session.NewManager(redisClient)
}

func HandleRequest(ctx context.Context, event transport.Event) error {
	switch event.Type {
	case transport.AddFriend:
		event, err := utils.JSONConvert[transport.AddFriendEvent](event)
		if err != nil {
			log.Println("can not parse request payload:", err)
			return err
		}
		userConIDs, err := SessionManager.GetSessions(event.Payload.UserID)
		if err != nil {
			log.Println("can not get session:", err)
			return err
		}

		eventBytes, _ := json.Marshal(event)
		wg := sync.WaitGroup{}
		for _, conID := range userConIDs {
			wg.Add(1)
			go func(conID string, event []byte) {
				conID = strings.Split(conID, ":")[1]
				err = APIGatewayClient.Publish(ctx, conID, event)
				if err != nil {
					log.Println("failed to publish:", err)
				}
				wg.Done()
			}(conID, eventBytes)
		}
		wg.Wait()
	default:
		log.Print("does not support event type:", event.Type)
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
