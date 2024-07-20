package users

import (
	"net/http"

	"blinders/packages/auth"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MiddlewareKey string

const (
	UserID      MiddlewareKey = "user_id"
	PublicQuery MiddlewareKey = "public_query"
)

type ValidateOptions struct {
	// initializes local ctx with key PublicQuery for checking public in next middlewares.
	// If it is not enabled, the next middlewares will never receive public flag even if it handles public query.
	// It make easily to toggle public query for any route
	PublicQuery bool
}

func ValidateUserIDParam(options ...ValidateOptions) fiber.Handler {
	if len(options) == 0 {
		options = append(options, ValidateOptions{PublicQuery: false})
	}

	allowPublicQuery := options[0].PublicQuery

	return func(ctx *fiber.Ctx) error {
		userIDParam := ctx.Params("id")
		userID := ctx.Locals(auth.UserIDKey).(primitive.ObjectID)

		if allowPublicQuery {
			public := ctx.Query("public")
			ctx.Locals(PublicQuery, public == "true")
		} else {
			ctx.Locals(PublicQuery, false)
		}
		isPublicQuery := ctx.Locals(PublicQuery).(bool)

		if !isPublicQuery && userID.Hex() != userIDParam {
			return ctx.Status(http.StatusForbidden).JSON(&fiber.Map{
				"error": "insufficient permissions",
			})
		}

		return ctx.Next()
	}
}
