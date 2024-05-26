package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
}

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

const ExplainPhraseInSentencePrompt = `
<|begin_of_text|>

<|start_header_id|>system<|end_header_id|> 
You are a helpful AI assistant for learning English for Vietnamese learner. Only response to user a JSON object, the result must be short, clear
<|eot_id|>

<|start_header_id|>user<|end_header_id|>
I want to understand:
- word/phrase: "%v"
- sentence: "%v"
<|eot_id|>

<|start_header_id|>system<|end_header_id|> 
Only respond JSON format: 
    {
        "translate": translate the word/phrase to Vietnamese,
        "IPA": IPA English pronunciation of word/phrase,
        "grammar_analysis": {
            "tense": {"type": type of tense of the whole sentence, "identifier": how user can identify the tense},
            "structure": {"type": structure type of the whole sentence,
                "structure": show the grammar structure of the sentence as form example 'I know that + S + has been + V_ed +',
                "for": how the structure is used for
            }
        },
        "key_words": get 1, or 2, or 3 main words in the sentence (it can be noun/verb/adjective),
        "expand_words": give 3 words might be relevant but not in the sentence
    }
<|eot_id|>

<|start_header_id|>assistant<|end_header_id|>
`
