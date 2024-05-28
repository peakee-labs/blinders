package core_test

import (
	"context"
	"testing"

	"blinders/services/embedder/core"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/test-go/testify/assert"
)

func TestEmbedder(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("ap-south-1"),
		config.WithSharedConfigProfile("admin.peakee"),
	)
	assert.Nil(t, err)

	brrc, err := core.InitBedrockRuntimeClientFromConfig(cfg)
	assert.Nil(t, err)

	embedder := core.NewEmbbeder(brrc, "cohere.embed-english-v3")
	prompt := "this text is to test embedder core"
	embedded, err := embedder.Embedding(prompt)
	assert.Nil(t, err)
	assert.Equal(t, 1024, len(embedded))
}
