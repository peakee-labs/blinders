package main

import (
	"context"
	"log"
	"os"
	"time"

	"blinders/packages/transport"

	"github.com/aws/aws-sdk-go-v2/config"
)

func InitTransport(ctx context.Context) chan transport.Transport {
	ch := make(chan transport.Transport)
	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Println("failed to load aws config:", err)
			ch <- nil
		}
		consumerMap := transport.ConsumerMap{
			transport.CollectingPush: os.Getenv("COLLECTING_PUSH_FUNCTION_NAME"),
		}
		transporter := transport.NewLambdaTransportWithConsumers(cfg, consumerMap)

		ch <- transporter
	}()

	return ch
}
