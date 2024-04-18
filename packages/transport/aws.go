package transport

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type LambdaTransport struct {
	*lambda.Client
}

func NewLambdaTransport(cfg aws.Config) *LambdaTransport {
	return &LambdaTransport{lambda.NewFromConfig(cfg)}
}

func (t LambdaTransport) Request(
	ctx context.Context,
	id string,
	payload []byte,
) (response []byte, err error) {
	log.Println("[lambda transport] request to", id)
	res, err := t.Invoke(ctx, &lambda.InvokeInput{
		InvocationType: types.InvocationTypeRequestResponse,
		FunctionName:   &id,
		Payload:        payload,
	})
	if err != nil {
		return nil, err
	}

	return res.Payload, nil
}

func (t LambdaTransport) Push(ctx context.Context, id string, payload []byte) error {
	log.Println("[lambda transport] push to", id)
	_, err := t.Invoke(ctx, &lambda.InvokeInput{
		InvocationType: types.InvocationTypeEvent,
		FunctionName:   &id,
		Payload:        payload,
	})

	return err
}
