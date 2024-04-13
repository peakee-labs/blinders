package transport

import (
	"context"
)

type RequestConfig struct {
	Header map[string][]string
	Method string
}

type Transport interface {
	Request(ctx context.Context, id string, payload []byte, config ...RequestConfig) (response []byte, err error)
	Push(ctx context.Context, id string, payload []byte) error
	// Do(ctx context.Context, id string, payload []byte, config RequestConfig) (response []byte, err error)
}

type Key string

const (
	Notification Key = "notification"
	Explore      Key = "explore"
	Logging      Key = "logging"
	Suggest      Key = "suggest"
)

type ConsumerMap map[Key]string
