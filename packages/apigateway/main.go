package apigateway

import (
	"context"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	agm "github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

type Client struct {
	agm.Client
}

// NewClient creates a new API Gateway Management Client instance
// from the provided parameters. The new client will have a custom endpoint
// that resolves to the application's deployed API.
func NewClient(_ context.Context, cfg aws.Config, cer CustomEndpointResolve) *Client {
	return &Client{*agm.NewFromConfig(cfg, agm.WithEndpointResolverV2(cer))}
}

func (c Client) Publish(ctx context.Context, connectionID string, data []byte) error {
	_, err := c.PostToConnection(ctx, &agm.PostToConnectionInput{
		ConnectionId: &connectionID,
		Data:         data,
	})

	return err
}

type CustomEndpointResolve struct {
	Domain, PathPrefix string
}

func (c CustomEndpointResolve) ResolveEndpoint(_ context.Context, _ agm.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	var endpoint url.URL
	endpoint.Scheme = "https" // use https by default
	endpoint.Host = c.Domain
	endpoint.Path = c.PathPrefix

	return smithyendpoints.Endpoint{
		URI:     endpoint,
		Headers: http.Header{},
	}, nil
}
