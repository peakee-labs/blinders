package auth

import (
	"context"
	"log"
	"strings"

	"blinders/packages/db/usersdb"

	"github.com/aws/aws-lambda-go/events"
)

type (
	LambdaHandler    func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
	LambdaMiddleware func(next LambdaHandler) LambdaHandler
)

func LambdaLoggingMiddleware() LambdaMiddleware {
	return func(next LambdaHandler) LambdaHandler {
		return func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			log.Println(event)
			return next(ctx, event)
		}
	}
}

func LambdaAuthMiddleware(
	m Manager,
	userRepo *usersdb.UsersRepo,
	options ...MiddlewareOptions,
) LambdaMiddleware {
	return func(next LambdaHandler) LambdaHandler {
		return func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			auth, ok := event.Headers["authorization"]
			if !ok {
				return events.APIGatewayV2HTTPResponse{
					StatusCode: 400,
					Body:       "missing authorization header",
					Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				}, nil
			}
			if !strings.HasPrefix(auth, "Bearer ") {
				log.Println("invalid jwt, missing bearer token")
				return events.APIGatewayV2HTTPResponse{
					StatusCode: 400,
					Body:       "missing bearer token",
					Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				}, nil
			}

			jwt := strings.Split(auth, " ")[1]
			userAuth, err := m.Verify(jwt)
			if err != nil {
				log.Println("failed to verify jwt:", err)
				return events.APIGatewayV2HTTPResponse{
					StatusCode: 400,
					Body:       "failed to verify jwt",
					Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				}, nil
			}

			if len(options) == 0 || options[0].CheckUser {
				// currently, user.AuthID is firebaseUID
				user, err := userRepo.GetUserByFirebaseUID(userAuth.AuthID)
				if err != nil {
					log.Println("failed to get user:", err)
					return events.APIGatewayV2HTTPResponse{
						StatusCode: 400,
						Body:       "failed to get user",
						Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
					}, nil
				}
				userAuth.ID = user.ID.Hex()
			}

			ctx = context.WithValue(ctx, UserAuthKey, userAuth)
			return next(ctx, event)
		}
	}
}

func LambdaAuthMiddlewareFromChan(
	mCh chan Manager,
	userRepoCh chan *usersdb.UsersRepo,
	options ...MiddlewareOptions,
) LambdaMiddleware {
	return func(next LambdaHandler) LambdaHandler {
		return func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			auth, ok := event.Headers["authorization"]
			if !ok {
				return events.APIGatewayV2HTTPResponse{
					StatusCode: 400,
					Body:       "missing authorization header",
					Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				}, nil
			}
			if !strings.HasPrefix(auth, "Bearer ") {
				log.Println("invalid jwt, missing bearer token")
				return events.APIGatewayV2HTTPResponse{
					StatusCode: 400,
					Body:       "missing bearer token",
					Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				}, nil
			}

			jwt := strings.Split(auth, " ")[1]
			m := <-mCh
			userAuth, err := m.Verify(jwt)
			if err != nil {
				log.Println("failed to verify jwt:", err)
				return events.APIGatewayV2HTTPResponse{
					StatusCode: 400,
					Body:       "failed to verify jwt",
					Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				}, nil
			}

			if len(options) == 0 || options[0].CheckUser {
				// currently, user.AuthID is firebaseUID
				userRepo := <-userRepoCh
				user, err := userRepo.GetUserByFirebaseUID(userAuth.AuthID)
				if err != nil {
					log.Println("failed to get user:", err)
					return events.APIGatewayV2HTTPResponse{
						StatusCode: 400,
						Body:       "failed to get user",
						Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
					}, nil
				}
				userAuth.ID = user.ID.Hex()
			}

			ctx = context.WithValue(ctx, UserAuthKey, userAuth)
			return next(ctx, event)
		}
	}
}
