package exploreapi

import (
	"blinders/packages/auth"
	"blinders/packages/db"

	"github.com/gofiber/fiber/v2"
)

type Manager struct {
	App     *fiber.App
	Auth    auth.Manager
	DB      *db.MongoManager
	Service *Service
}

func NewManager(app *fiber.App, auth auth.Manager, db *db.MongoManager, service *Service) *Manager {
	return &Manager{
		App:     app,
		Auth:    auth,
		DB:      db,
		Service: service,
	}
}

func (m *Manager) InitRoute() {
	exploreRoute := m.App.Group("/explore")
	exploreRoute.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	authorizedRoute := exploreRoute.Group("/", auth.FiberAuthMiddleware(m.Auth, m.DB.Users))
	authorizedRoute.Get("/suggest", m.Service.HandleGetMatches)
}
