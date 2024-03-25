/*
This package provides transport layer for all services, for both local development and production on AWS
*/
package transport

import "context"

type Transport interface {
	Request(ctx context.Context, id string, payload []byte) (response []byte, err error)
	Push(ctx context.Context, id string, payload []byte) error
}

type Key string

const (
	Notification Key = "notification"
)

type ConsumerMap map[Key]string
