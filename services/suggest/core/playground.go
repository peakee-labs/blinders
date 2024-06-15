package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"blinders/packages/db/collectingdb"
	"blinders/packages/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

func ListModels(cfg aws.Config) {
	brc := bedrock.NewFromConfig(cfg)
	result, err := brc.ListFoundationModels(
		context.TODO(),
		&bedrock.ListFoundationModelsInput{},
	)
	if err != nil {
		fmt.Printf("Couldn't list foundation models. Here's why: %v\n", err)
		return
	}
	if len(result.ModelSummaries) == 0 {
		fmt.Println("There are no foundation models.")
	}
	for _, modelSummary := range result.ModelSummaries {
		fmt.Println(*modelSummary.ModelId)
	}
}

func RunPlayground(cfg aws.Config) {
	brrc := bedrockruntime.NewFromConfig(cfg)
	phrase := "technical sharing"
	sentence := "Hey guys, this channel is used for technical sharing, anyone can share link/docs or even your short blog"
	req := LlamaRequest{
		Prompt:      fmt.Sprintf(ExplainPhraseInSentencePrompt, phrase, sentence),
		MaxGenLen:   512,
		Temperature: 0.5,
		TopP:        0.9,
	}
	reqBytes, _ := json.Marshal(&req)
	result, err := brrc.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String("meta.llama3-70b-instruct-v1:0"),
		ContentType: aws.String("application/json"),
		Body:        reqBytes,
	})
	if err != nil {
		log.Println(err)
		return
	}
	res, err := utils.ParseJSON[LlamaResponse](result.Body)
	if err != nil {
		log.Println("can not parse response:", err)
	}
	log.Println("result:", res.Generation)
	log.Println("generation token count:", res.GenerationTokenCount)
	log.Println("prompt token count:", res.PromptTokenCount)
	log.Println("stop reason:", res.StopReason)

	explainResult, err := utils.ParseJSON[collectingdb.ExplainResponse]([]byte(res.Generation))
	log.Println("parsed:", explainResult, err)
}
