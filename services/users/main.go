package users

import (
	"blinders/packages/auth"

	"blinders/services/users/repo"

	"github.com/gofiber/fiber/v2"
)

type Manager struct {
	App       *fiber.App
	Auth      auth.Manager
	UsersRepo *repo.UsersRepo
}

func NewManager(
	app *fiber.App,
	auth auth.Manager,
) *Manager {
	return &Manager{
		App:  app,
		Auth: auth,
	}
}
