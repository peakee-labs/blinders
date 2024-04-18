package transport_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"blinders/packages/transport"

	"github.com/test-go/testify/require"
)

var (
	successEndpoint = "/success"
	successResponse = []byte("success-get")
	failedEndpoint  = "/failed"
	failedResponse  = []byte("failed-get")
	postBody        = []byte("body")
	postBodyNil     = make([]byte, 0)
)

func TestLocalTransportPost(t *testing.T) {
	prefix := "/post"
	s := InitMockServer(prefix)
	s.Start()
	defer s.Close()
	type Testcase struct {
		name            string
		transportClient *http.Client
		endpoint        string
		expectedError   bool
		body            []byte
		rsp             []byte
	}
	testcases := []Testcase{
		{
			name:            "DoSuccess",
			transportClient: s.Client(),
			endpoint:        fmt.Sprintf("%s%s%s", s.URL, prefix, successEndpoint),
			expectedError:   false,
			body:            postBody,
			rsp:             postBody,
		},
		{
			name:            "DoSuccessWithEmptyResponse",
			transportClient: s.Client(),
			endpoint:        fmt.Sprintf("%s%s%s", s.URL, prefix, successEndpoint),
			expectedError:   false,
			body:            postBodyNil,
			rsp:             nil,
		},
		{
			name:            "DoFailed",
			transportClient: s.Client(),
			endpoint:        fmt.Sprintf("%s%s%s", s.URL, prefix, failedEndpoint),
			expectedError:   true,
			body:            postBody,
			rsp:             nil,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			tp := transport.NewLocalTransport(s.Client())
			rsp, err := tp.Request(context.TODO(), tc.endpoint, tc.body, transport.RequestConfig{
				Method: "POST",
			})
			if tc.expectedError {
				require.NotNil(t, err)
			}
			require.Equal(t, tc.rsp, rsp)
		})
	}
}

func TestLocalTransportRequest(t *testing.T) {
	prefix := "/request"
	s := InitMockServer(prefix)
	s.Start()
	defer s.Close()
	type Testcase struct {
		name            string
		transportClient *http.Client
		endpoint        string
		expectedError   bool
		rsp             []byte
	}
	testcases := []Testcase{
		{
			name:            "RequestSuccess",
			transportClient: s.Client(),
			endpoint:        fmt.Sprintf("%s%s%s", s.URL, prefix, successEndpoint),
			expectedError:   false,
			rsp:             successResponse,
		},
		{
			name:            "RequestFailed",
			transportClient: s.Client(),
			endpoint:        fmt.Sprintf("%s%s%s", s.URL, prefix, failedEndpoint),
			expectedError:   true,
			rsp:             failedResponse,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			tp := transport.NewLocalTransport(s.Client())
			rsp, err := tp.Request(context.TODO(), tc.endpoint, nil)
			if tc.expectedError {
				require.NotNil(t, err)
			}

			require.Equal(t, tc.rsp, rsp)
		})
	}
}

var ServeSuccessGetEndpoint http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	switch request.Method {
	case "POST":
		_, _ = io.Copy(writer, request.Body)
	default:
		_, _ = writer.Write(successResponse)
	}
}

var ServeFailedGetEndpoint http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusBadRequest)
	switch request.Method {
	case "POST":
		return
	default:
		_, _ = writer.Write(failedResponse)
	}
}

func InitMockServer(prefix string) *httptest.Server {
	http.Handle(fmt.Sprintf("%s%s", prefix, successEndpoint), ServeSuccessGetEndpoint)
	http.Handle(fmt.Sprintf("%s%s", prefix, failedEndpoint), ServeFailedGetEndpoint)
	server := httptest.NewUnstartedServer(http.DefaultServeMux)
	return server
}
