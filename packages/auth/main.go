package auth

import (
	"context"
	"fmt"
	"strings"

	"blinders/packages/apigateway"
	"blinders/packages/lambda"
	"blinders/packages/utils"
	usersrepo "blinders/services/users/repo"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/aws/aws-lambda-go/events"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
)

type UserAuth struct {
	Email  string
	Name   string
	AuthID string
}

type Manager struct {
	App    *firebase.App
	Client *auth.Client

	// TODO: decouple auth and users repo
	UsersRepo *usersrepo.UsersRepo
}

type key string

const (
	UserKey     key = "user_key"
	UserIDKey   key = "user_id_key"
	UserAuthKey key = "user_auth_key"
)

func NewFirebaseManager(adminConfig []byte, usersRepo *usersrepo.UsersRepo) (*Manager, error) {
	manager := Manager{UsersRepo: usersRepo}

	opt := option.WithCredentialsJSON(adminConfig)
	newApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	manager.App = newApp

	newClient, err := manager.App.Auth(context.Background())
	if err != nil {
		return nil, err
	}
	manager.Client = newClient

	return &manager, nil
}

func NewFirebaseManagerFromFile(filename string, usersRepo *usersrepo.UsersRepo) (*Manager, error) {
	adminConfig, err := utils.GetFile(filename)
	if err != nil {
		return nil, fmt.Errorf("can not load firebase config file: %v", err)
	}

	m, err := NewFirebaseManager(adminConfig, usersRepo)
	if err != nil {
		return nil, fmt.Errorf("can not init firebase manager: %v", err)
	}

	return m, nil
}

func (m Manager) Verify(jwt string) (*UserAuth, error) {
	authToken, err := m.Client.VerifyIDToken(context.Background(), jwt)
	if err != nil {
		return nil, err
	}

	firebaseUID := authToken.UID
	email := authToken.Claims["email"].(string)
	name := authToken.Claims["name"].(string)

	userAuth := UserAuth{
		Email:  email,
		Name:   name,
		AuthID: firebaseUID,
	}

	return &userAuth, nil
}

type Config struct {
	WithUser bool
}

func (m Manager) FiberAuthMiddleware(cfg ...Config) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := ctx.Get("Authorization")
		if auth == "" {
			auth = ctx.Get("authorization")
		}
		if !strings.HasPrefix(auth, "Bearer ") {
			return ctx.Status(fiber.StatusUnauthorized).
				SendString("invalid jwt, missing bearer token")
		}

		jwt := strings.Split(auth, " ")[1]
		userAuth, err := m.Verify(jwt)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}

		ctx.Locals(UserAuthKey, userAuth)

		if len(cfg) > 0 && cfg[0].WithUser {
			user, err := m.UsersRepo.GetUserByFirebaseUID(userAuth.AuthID)
			if err != nil {
				return ctx.Status(fiber.StatusUnauthorized).SendString(err.Error())
			}

			ctx.Locals(UserKey, user)
			ctx.Locals(UserIDKey, user.ID)
		}

		return ctx.Next()
	}
}

func (m Manager) LambdaAuthMiddleware(cfg ...Config) lambda.Middleware {
	return func(next lambda.Handler) lambda.Handler {
		return func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			auth, ok := event.Headers["authorization"]
			if !ok {
				return apigateway.UnauthorizedResponse("missing authorization header"), nil
			}
			if !strings.HasPrefix(auth, "Bearer ") {
				return apigateway.UnauthorizedResponse("missing bearer token"), nil
			}

			jwt := strings.Split(auth, " ")[1]
			userAuth, err := m.Verify(jwt)
			if err != nil {
				return apigateway.UnauthorizedResponse("can not verify JWT"), nil
			}

			ctx = context.WithValue(ctx, UserAuthKey, userAuth)

			if len(cfg) > 0 && cfg[0].WithUser {
				user, err := m.UsersRepo.GetUserByFirebaseUID(userAuth.AuthID)
				if err != nil {
					return apigateway.UnauthorizedResponse(err.Error()), nil
				}

				ctx = context.WithValue(ctx, UserKey, user)
				ctx = context.WithValue(ctx, UserIDKey, user.ID)
			}

			return next(ctx, event)
		}
	}
}
