package core

import (
	"context"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/assert"
)

func TestRunBedrockClient(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("ap-south-1"),
		config.WithSharedConfigProfile("tanle.peakee.admin"),
	)
	log.Println("ok")

	assert.Nil(t, err)

	ListModels(cfg)
	RunPlayground(cfg)
}
