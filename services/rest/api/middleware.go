package restapi

import (
	"net/http"

	"blinders/packages/auth"

	"github.com/gofiber/fiber/v2"
)

type MiddlewareKey string

const (
	PublicQuery MiddlewareKey = "public_query"
)

type ValidateUserOptions struct {
	// initializes local ctx with key PublicQuery for checking public in next middlewares.
	// If it is not enabled, the next middlewares will never receive public flag even if it handles public query.
	// It make easily to toggle public query for any route
	allowPublicQuery bool
}

func ValidateUserIDParam(options ...ValidateUserOptions) fiber.Handler {
	if len(options) == 0 {
		options = append(options, ValidateUserOptions{
			allowPublicQuery: false,
		})
	}

	allowPublicQuery := options[0].allowPublicQuery

	return func(ctx *fiber.Ctx) error {
		userID := ctx.Params("id")
		userAuth := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)

		if allowPublicQuery {
			public := ctx.Query("public")
			ctx.Locals(PublicQuery, public == "true")
		} else {
			ctx.Locals(PublicQuery, false)
		}
		isPublicQuery := ctx.Locals(PublicQuery).(bool)

		if !isPublicQuery && userAuth.ID != userID {
			return ctx.Status(http.StatusForbidden).JSON(&fiber.Map{
				"error": "insufficient permissions",
			})
		}

		return ctx.Next()
	}
}
