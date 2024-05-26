package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"blinders/packages/db/collectingdb"
	"blinders/packages/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// https://docs.aws.amazon.com/bedrock/latest/userguide/model-parameters-meta.html
type LlamaRequest struct {
	// Prompt format: https://llama.meta.com/docs/model-cards-and-prompt-formats/meta-llama-3/
	Prompt      string  `json:"prompt"`
	Temperature float32 `json:"temperature"`
	TopP        float32 `json:"top_p"`
	MaxGenLen   int     `json:"max_gen_len"`
}

type LlamaResponse struct {
	Generation           string `json:"generation"` // "\n\n<response>"
	PromptTokenCount     int    `json:"prompt_token_count"`
	GenerationTokenCount int    `json:"generation_token_count"`
	StopReason           string `json:"stop_reason"` // "stop" || "length"
}

func ExplainPhraseInSentence(
	brrc bedrockruntime.Client,
	phrase string,
	sentence string,
) (*collectingdb.ExplainResponse, error) {
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
		log.Println("can not request to get explanation:", err)
		return nil, fmt.Errorf("can not request to get explanation")
	}
	res, err := utils.ParseJSON[LlamaResponse](result.Body)
	if err != nil {
		log.Println("can not parse response:", err)
	}

	explainResult, err := utils.ParseJSON[collectingdb.ExplainResponse]([]byte(res.Generation))
	if err != nil {
		log.Println("can not parse explanation response, something went wrong:", err)
		return nil, fmt.Errorf("can not parse explanation response, something went wrong")
	}

	return explainResult, nil
}
