package main

import (
	"blinders/functions/embedder/core"
	"blinders/packages/transport"
	"blinders/packages/utils"
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	embedder *core.Embedder
)

func init() {
	st := time.Now()
	env := os.Getenv("ENVIRONMENT")
	log.Printf("Embbeder api running on %s environment\n", env)

	brrc := core.InitBedrockRuntimeClientFromConfig()
	modelID := os.Getenv("EMBEDDER_MODEL_ID")
	if modelID == "" {
		modelID = "cohere.embed-english-v3"
	}

	embedder = core.NewEmbbeder(brrc, modelID)

	log.Printf("Cold-start duration: %v ms", time.Since(st).Milliseconds())
}

func HandleRequest(
	ctx context.Context,
	req any,
) (any, error) {
	internalReq, err := utils.JSONConvert[transport.Request](req)
	if err != nil {
		log.Fatalln("cannot parse proxy request", err)
	}
	if internalReq.Type != transport.Embedding {
		log.Fatalln("invalid request type", internalReq.Type)
	}

	embeddingReq, err := utils.JSONConvert[transport.EmbeddingRequest](req)
	if err != nil {
		log.Fatalln("cannot parse embedding request", err)
	}

	return HandleInternalEmbeddingRequest(ctx, *embeddingReq)
}

func HandleInternalEmbeddingRequest(
	ctx context.Context,
	req transport.EmbeddingRequest,
) (transport.EmbeddingResponse, error) {
	embeddVector, err := embedder.Embedding(req.Payload)
	if err != nil {
		log.Fatalln("cannot get embedding", err)
	}
	res := transport.EmbeddingResponse{
		Embedded: embeddVector,
	}
	return res, nil
}

func main() {
	lambda.Start(HandleRequest)
}
