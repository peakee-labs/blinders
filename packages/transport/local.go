package transport

import (
	"context"
	"log"
)

type LocalTransport struct{}

func NewLocalTransport() *LocalTransport {
	return &LocalTransport{}
}

func (t LocalTransport) Request(
	_ context.Context,
	id string,
	_ []byte,
) (response []byte, err error) {
	log.Println("[local transport] request to", id)
	return nil, nil
}

func (t LocalTransport) Push(_ context.Context, id string, _ []byte) error {
	log.Println("[local transport] push to", id)
	return nil
}
