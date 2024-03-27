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

type InitOptions struct {
	Prefix string
}

func (m *Manager) InitRoute(options InitOptions) {
	if options.Prefix == "" {
		options.Prefix = "/"
	}

	m.App.Get(options.Prefix+"/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	routes := m.App.Group(options.Prefix, auth.FiberAuthMiddleware(m.Auth, m.DB.Users))
	routes.Get("/suggest", m.Service.HandleGetMatches)
}
