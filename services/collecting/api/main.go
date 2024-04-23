package collectingapi

import (
	"blinders/packages/auth"
	"blinders/packages/collecting"

	"github.com/gofiber/fiber/v2"
)

type (
	Manager struct {
		App               *fiber.App
		Auth              auth.Manager
		CollectingService Service
	}
	Service struct {
		Collector *collecting.EventCollector
	}
)

func NewManager(app *fiber.App, collector *collecting.EventCollector) *Manager {
	return &Manager{
		App:               app,
		CollectingService: Service{Collector: collector},
	}
}

func (s *Manager) InitRoute() error {
	s.App.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("service healthy")
	})

	s.App.Post("/", s.CollectingService.HandlePushEvent)
	s.App.Get("/", s.CollectingService.HandleGetEvent)
	return nil
}
