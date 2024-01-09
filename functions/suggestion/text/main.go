package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"blinders/functions/utils"
	"blinders/packages/suggestion"
	"blinders/packages/user"
	cutils "blinders/packages/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sashabaranov/go-openai"
)

var (
	suggester suggestion.Suggester
	apiKey    = os.Getenv("OPENAI_API_KEY")
)

func init() {
	if suggester == nil {
		var err error
		client := openai.NewClient(apiKey)
		suggester, err = suggestion.NewGPTSuggester(client)
		if err != nil {
			panic(err)
		}
	}
}

type SuggestionPayload struct {
	Text string `json:"text"`
}

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	token, ok := event.Headers["authorization"]
	if !ok {
		return utils.APIGatewayProxyResponseWithJSON(400, map[string]any{
			"error": "function: Token in authorization header not found",
		})
	}

	usr, err := cutils.VerifyFireStoreToken(token)
	if err != nil {
		return utils.APIGatewayProxyResponseWithJSON(400, map[string]any{
			"error": fmt.Sprintf("function: Cannot verify given token, err: %s", err.Error()),
		})
	}

	suggestionRequest := new(SuggestionPayload)
	if err := json.Unmarshal([]byte(event.Body), suggestionRequest); err != nil {
		return utils.APIGatewayProxyResponseWithJSON(400, map[string]any{
			"error": fmt.Sprintf("functions: cannot unmarshal struct from json, err: (%s)", err.Error()),
		})
	}

	userData, err := user.GetUserData(usr.ID)
	if err != nil {
		return utils.APIGatewayProxyResponseWithJSON(400, map[string]any{
			"error": fmt.Sprintf("functions: cannot get data of user, err: (%s)", err.Error()),
		})
	}

	suggestion, err := suggester.TextCompletion(ctx, userData, suggestionRequest.Text)
	if err != nil {
		return utils.APIGatewayProxyResponseWithJSON(400, map[string]any{
			"error": fmt.Sprintf("functions: cannot get suggestions, err: (%s)", err.Error()),
		})
	}

	return utils.APIGatewayProxyResponseWithJSON(200, map[string]any{
		"suggestions": suggestion,
	})
}

func main() {
	lambda.Start(HandleRequest)
}
