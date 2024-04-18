// Package transport provides transport layer for all services, for both local development and production on AWS
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
}

type Key string

const (
	Notification Key = "notification"
	Explore      Key = "explore"
	Collecting   Key = "collecting"
	Suggest      Key = "suggest"
)

type ConsumerMap map[Key]string
