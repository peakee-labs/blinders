package core

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

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
