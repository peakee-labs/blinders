package core

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type Embedder struct {
	Client  *bedrockruntime.Client
	ModelID string
}

// https://docs.aws.amazon.com/bedrock/latest/userguide/model-parameters-embed.html#model-parameters-embed-request-response
type CohereRequest struct {
	String    []string `json:"texts"`
	InputType string   `json:"input_type"`
}

type CohereResponse struct {
	Embedding    [][]float32 `json:"embeddings"`
	ID           string      `json:"id"`
	ResponseType string      `json:"response_type"`
	Texts        []string    `json:"texts"`
}

func NewEmbbeder(model *bedrockruntime.Client, modelID string) *Embedder {
	return &Embedder{
		Client:  model,
		ModelID: modelID,
	}
}

func (e *Embedder) Embedding(prompt string) ([]float32, error) {
	req := CohereRequest{
		String:    []string{prompt},
		InputType: "search_query",
	}

	bodyBytes, _ := json.Marshal(req)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	res, err := e.Client.InvokeModel(ctx,
		&bedrockruntime.InvokeModelInput{
			Body:        bodyBytes,
			ModelId:     aws.String(e.ModelID),
			ContentType: aws.String("application/json"),
			Accept:      aws.String("*/*"),
		})
	if err != nil {
		log.Printf("can not request to get embedding: %v", err)
		return nil, err
	}

	var embedderRes CohereResponse
	if err := json.Unmarshal(res.Body, &embedderRes); err != nil {
		log.Printf("can not parse response: %v", err)
		return nil, err
	}

	if len(embedderRes.Embedding) != 1 {
		log.Printf("unexpected error: embedding in response length != 1, received length: %d\n", len(embedderRes.Embedding))
		return nil, err
	}

	return embedderRes.Embedding[0], nil
}
