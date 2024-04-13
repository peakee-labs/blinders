package transport

import (
	"bytes"
	"context"
	"encoding/json"
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

func (t LocalTransport) Request(
	_ context.Context,
	id string,
	payload []byte,
	config ...RequestConfig,
) (response []byte, err error) {
	method := "GET"
	header := make(http.Header)
	if len(config) == 1 {
		method = config[0].Method
		for key, vals := range config[0].Header {
			for _, val := range vals {
				header.Add(key, val)
			}
		}
	}

	req, err := http.NewRequest(method, id, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header = header

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
	if 200 < rsp.StatusCode || rsp.StatusCode > 299 {
		msg := ParseResponseMessage(rsp.Body)
		if msg != "" {
			return bodyReader.Bytes(), fmt.Errorf(msg)
		}
		return bodyReader.Bytes(), fmt.Errorf("cannot make push request to target, statuscode: %d", rsp.StatusCode)
	}

	return bodyReader.Bytes(), nil
}

func (t LocalTransport) Push(_ context.Context, id string, body []byte) error {
	log.Println("[local transport] push to", id)
	rsp, err := t.client.Post(id, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	if 200 > rsp.StatusCode || rsp.StatusCode > 299 {
		msg := ParseResponseMessage(rsp.Body)
		if msg != "" {
			return fmt.Errorf(msg)
		}
		return fmt.Errorf("cannot make push request to target, statuscode: %d", rsp.StatusCode)
	}
	return nil
}

// ParseResponseMessage tries to parse the response body from a failed request.
// The message could be passed via the 'error' field of the JSON body.
func ParseResponseMessage(body io.ReadCloser) string {
	bodyReader := new(bytes.Buffer)
	_, err := io.Copy(bodyReader, body)
	if err != nil {
		return ""
	}

	res := make(map[string]string)
	if err := json.Unmarshal(bodyReader.Bytes(), &res); err != nil {
		return ""
	}
	return res["error"]
}
