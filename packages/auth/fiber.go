package auth

import (
	"log"
	"strings"

	usersrepo "blinders/services/users/repo"

	"github.com/gofiber/fiber/v2"
)

type key string

const UserAuthKey key = "user_auth_key"

type MiddlewareOptions struct {
	// permit not checking if user exists, should only use to initialize user
	CheckUser bool
}

func FiberAuthMiddleware(
	m Manager,
	userRepo *usersrepo.UsersRepo,
	options ...MiddlewareOptions,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := ctx.Get("Authorization")
		if auth == "" {
			auth = ctx.Get("authorization")
		}
		if !strings.HasPrefix(auth, "Bearer ") {
			log.Println("invalid jwt, missing bearer token")
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid jwt, missing bearer token",
			})
		}

		jwt := strings.Split(auth, " ")[1]
		userAuth, err := m.Verify(jwt)
		if err != nil {
			log.Println("failed to verify jwt:", err)
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "failed to verify jwt",
			})
		}

		if len(options) == 0 || options[0].CheckUser {
			// currently, user.AuthID is firebaseUID
			user, err := userRepo.GetUserByFirebaseUID(userAuth.AuthID)
			if err != nil {
				log.Println("failed to get user:", err)
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "failed to get user",
				})
			}

			userAuth.ID = user.ID.Hex()
		}

		ctx.Locals(UserAuthKey, userAuth)

		return ctx.Next()
	}
}
