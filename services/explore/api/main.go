package exploreapi

import (
	"blinders/packages/auth"
	"blinders/packages/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	m.App.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("service healthy")
	})

	matchRoute := m.App.Group("/explore", cors.New(), auth.FiberAuthMiddleware(m.Auth, m.DB.Users))
	matchRoute.Get("/suggest", m.Service.HandleGetMatches)
	// Temporarily expose this method, it must be call internal, or we will listen to user update-match-information event
	matchRoute.Post("/", m.Service.HandleAddUserMatch)
}
