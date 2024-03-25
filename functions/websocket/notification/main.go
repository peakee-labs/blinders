/*
This service is responsible for notifying any event to users via websocket or push notification
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/redis/go-redis/v9"
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

	SessionManager = session.NewManager(redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}))
}

func HandleRequest(ctx context.Context, payload any) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("can not marshal payload:", err)
		return err
	}
	event, err := utils.ParseJSON[transport.Event](bytes)
	if err != nil {
		log.Println("can not parse request payload, require type in payload:", err)
		return err
	}

	log.Println("handle event:", event.Type)

	switch event.Type {
	case transport.AddFriend:
		event, err := utils.ParseJSON[transport.AddFriendEvent](bytes)
		if err != nil {
			log.Println("can not parse request payload:", err)
			return err
		}
		userConIDs, err := SessionManager.GetSessions(event.UserID)
		if err != nil {
			log.Println("can not get session:", err)
			return err
		}

		wg := sync.WaitGroup{}
		for _, conID := range userConIDs {
			wg.Add(1)
			go func(conID string, payload []byte) {
				conID = strings.Split(conID, ":")[1]
				err = APIGatewayClient.Publish(ctx, conID, payload)
				if err != nil {
					log.Println("failed to publish:", err)
				}
				wg.Done()
			}(conID, bytes)
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
