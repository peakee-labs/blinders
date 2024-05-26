// Package transport provides transport layer for all services, for both local development and production on AWS
package transport

import (
	"context"
)

type Transport interface {
	Request(ctx context.Context, id string, payload []byte) (response []byte, err error)
	Push(ctx context.Context, id string, payload []byte) error
	ConsumerID(key Key) string
}

type Key string

const (
	Notification   Key = "notification"
	Explore        Key = "explore"
	CollectingPush Key = "collecting-push"
	CollectingGet  Key = "collecting-get"
	Suggest        Key = "suggest"
)

type ConsumerMap map[Key]string
