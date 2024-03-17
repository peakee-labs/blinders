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

func ValidateUserIDParam(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")
	userAuth := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	isPublicQuery := ctx.Locals(PublicQuery).(bool)
	if !isPublicQuery && userAuth.ID != userID {
		return ctx.Status(http.StatusForbidden).JSON(&fiber.Map{
			"error": "insufficient permissions",
		})
	}

	return ctx.Next()
}

// AllowPublicQuery middleware
// initializes local ctx with key PublicQuery for checking public in next middlewares.
// If it is not enabled, the next middlewares will never receive public flag even if it handles public query.
// It make easily to toggle public query for any route
func AllowPublicQuery(ctx *fiber.Ctx) error {
	public := ctx.Query("public")
	ctx.Locals(PublicQuery, public == "true")

	return ctx.Next()
}
