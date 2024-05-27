package core

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

func InitBedrockRuntimeClientFromConfig(awsConfig ...aws.Config) *bedrockruntime.Client {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var cfg aws.Config
	if len(awsConfig) == 1 {
		cfg = awsConfig[0]
	} else {
		var err error
		cfg, err = config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatal("can not load aws config")
		}
	}

	brrc := bedrockruntime.NewFromConfig(cfg)
	return brrc
}
