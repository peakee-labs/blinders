package core

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// experimental: not work for now
func InitBedrockRuntimeClient(ctx context.Context) chan *bedrockruntime.Client {
	ch := make(chan *bedrockruntime.Client)
	go func() {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			ch <- nil
		}
		brrc := bedrockruntime.NewFromConfig(cfg)
		ch <- brrc
	}()

	return ch
}

func InitBedrockRuntimeClientSync(ctx context.Context) *bedrockruntime.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("can not load aws config")
	}
	brrc := bedrockruntime.NewFromConfig(cfg)

	return brrc
}
