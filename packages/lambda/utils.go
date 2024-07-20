package lambda

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type (
	Handler    func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
	Middleware func(next Handler) Handler
)
