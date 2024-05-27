package exploreapi

import (
	"blinders/packages/auth"
	"blinders/packages/db/usersdb"

	"github.com/gofiber/fiber/v2"
)

type Manager struct {
	App       *fiber.App
	Auth      auth.Manager
	UsersRepo *usersdb.UsersRepo
	Service   *Service
}

func NewManager(
	app *fiber.App,
	auth auth.Manager,
	usersRepo *usersdb.UsersRepo,
	service *Service,
) *Manager {
	return &Manager{
		App:       app,
		Auth:      auth,
		UsersRepo: usersRepo,
		Service:   service,
	}
}

func (m *Manager) InitRoute() {
	exploreRoute := m.App.Group("/explore")
	exploreRoute.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	authorizedRoute := exploreRoute.Group("/", auth.FiberAuthMiddleware(m.Auth, m.UsersRepo))
	authorizedRoute.Get("/suggest", m.Service.HandleGetMatches)
	authorizedRoute.Get("/profiles/:id", m.Service.HandleGetMatchingProfile)
	authorizedRoute.Post("/profiles", m.Service.HandleAddMatchingProfile)
	authorizedRoute.Put("/profiles", m.Service.HandleUpdateMatchingProfile)
}
