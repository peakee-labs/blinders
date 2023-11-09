package main

import (
	"context"
	"fmt"
	"log"

	"blinders/packages/suggestion"
	"github.com/aws/aws-lambda-go/lambda"
)

type SuggestionPayload struct {
	Text string `json:"text"`
}

func HandleRequest(ctx context.Context, event *SuggestionPayload) (*string, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}
	log.Println("Suggestion triggered", event.Text)
	message := fmt.Sprintf("%s - from suggestion function!", suggestion.GetSuggestion(event.Text))

	return &message, nil
}

func main() {
	lambda.Start(HandleRequest)
}
