package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	wschat "blinders/functions/websocket/chat/core"
	"blinders/packages/apigateway"
	"blinders/packages/db/chatdb"
	dbutils "blinders/packages/db/utils"
	"blinders/packages/session"
	"blinders/packages/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

var APIGatewayClient *apigateway.Client

func init() {
	redisClient := utils.NewRedisClientFromEnv(context.Background())
	sessionManager := session.NewManager(redisClient)

	chatDB, err := dbutils.InitMongoDatabaseFromEnv("CHAT")
	if err != nil {
		log.Fatal(err)
	}

	wschat.InitChatApp(sessionManager, chatdb.NewChatDB(chatDB))

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("failed to load aws config:", err)
	}
	cer := apigateway.CustomEndpointResolve{
		Domain:     os.Getenv("API_GATEWAY_DOMAIN"),
		PathPrefix: os.Getenv("API_GATEWAY_PATH_PREFIX"),
	}
	APIGatewayClient = apigateway.NewClient(context.Background(), cfg, cer)
}

func HandleRequest(
	ctx context.Context,
	req events.APIGatewayWebsocketProxyRequest,
) (any, error) {
	connectionID := req.RequestContext.ConnectionID
	userID := req.RequestContext.Authorizer.(map[string]interface{})["principalId"].(string)

	genericEvent, err := utils.ParseJSON[wschat.ChatEvent]([]byte(req.Body))
	if err != nil {
		log.Println("can not parse request payload, require type in payload:", err)
	}

	switch genericEvent.Type {
	case wschat.UserPing:
		data := []byte("pong")
		err = APIGatewayClient.Publish(ctx, req.RequestContext.ConnectionID, data)
		if err != nil {
			log.Println("can not publish message:", err)
		}
	case wschat.UserSendMessage:
		payload, err := utils.ParseJSON[wschat.UserSendMessagePayload]([]byte(req.Body))
		if err != nil {
			log.Println("invalid send message event:", err)
			_ = APIGatewayClient.Publish(ctx, connectionID, []byte("invalid send message event"))
			break
		}

		dCh, err := wschat.HandleSendMessage(userID, connectionID, *payload)
		if err != nil {
			log.Println("failed to send message:", err)
			_ = APIGatewayClient.Publish(
				ctx,
				connectionID,
				[]byte("invalid payload to send message"),
			)
			break
		}

		wg := sync.WaitGroup{}
		for {
			d := <-dCh
			if d == nil {
				log.Println("distribute message channel closed")
				break
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				data, err := json.Marshal(d.Payload)
				if err != nil {
					log.Println("can not marshal data:", err)
					return
				}

				err = APIGatewayClient.Publish(ctx, d.ConnectionID, data)
				if err != nil {
					log.Println("can not publish message:", err)
				}
			}()
		}

		wg.Wait()
		log.Println("message sent")
	default:
		log.Println("not support this event:", req.Body)
		_ = APIGatewayClient.Publish(ctx, connectionID, []byte("not support this event"))
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
