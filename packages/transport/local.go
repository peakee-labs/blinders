package transport

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

var _ Transport = &LocalTransport{}

type LocalTransport struct {
	client *http.Client
}

func NewLocalTransport(client ...*http.Client) *LocalTransport {
	c := http.DefaultClient
	if len(client) == 1 {
		c = client[0]
	}
	return &LocalTransport{
		client: c,
	}
}

func (t LocalTransport) Do(
	_ context.Context,
	id string,
	payload []byte,
	config RequestConfig,
) (response []byte, err error) {
	req, err := http.NewRequest(config.Method, id, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	for key, vals := range config.Header {
		for _, val := range vals {
			req.Header.Add(key, val)
		}
	}

	rsp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyReader := new(bytes.Buffer)
	written, err := io.Copy(bodyReader, rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body, err: %v", err)
	}

	if rsp.ContentLength > 0 && rsp.ContentLength != written {
		return nil, fmt.Errorf("expected %d bytes from body, readed %d", rsp.ContentLength, written)
	}
	if 200 > rsp.StatusCode || rsp.StatusCode > 299 {
		return bodyReader.Bytes(), fmt.Errorf("received unexpected status code from %s, code: %d", id, rsp.StatusCode)
	}

	return bodyReader.Bytes(), nil
}

func (t LocalTransport) Request(
	ctx context.Context,
	addr string, // with LocalTransport, id is addr of target service
	body []byte,
) (response []byte, err error) {
	return t.Do(ctx, addr, body, RequestConfig{Method: "GET"})
}

func (t LocalTransport) Push(_ context.Context, id string, body []byte) error {
	log.Println("[local transport] push to", id)
	rsp, err := t.client.Post(id, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	if 200 > rsp.StatusCode || rsp.StatusCode > 299 {
		return fmt.Errorf("received unexpected status code from %s, code: %d", id, rsp.StatusCode)
	}
	return nil
}
