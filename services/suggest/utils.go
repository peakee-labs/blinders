package suggest

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

func InitBedrockRuntimeClientSync(ctx context.Context) *bedrockruntime.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("can not load aws config")
	}
	brrc := bedrockruntime.NewFromConfig(cfg)

	return brrc
}
