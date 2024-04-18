package transport

import (
	"context"
	"fmt"
)

type MockTransport struct{}

func (m MockTransport) Push(_ context.Context, id string, payload []byte) error {
	fmt.Printf("transport: Push to %v, payload: %v\n", id, payload)
	return nil
}

func (m MockTransport) Request(_ context.Context, id string, payload []byte, config ...RequestConfig) (response []byte, err error) {
	fmt.Printf("transport: Request to %v, payload: %v\n", id, payload)
	return nil, nil
}
