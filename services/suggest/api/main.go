package suggestapi

import (
	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/suggest"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	App       *fiber.App
	Auth      auth.Manager
	Db        *db.MongoManager
	Suggester suggest.Suggester
}

func (s *Service) InitRoute() {
	chatRoute := s.App.Group("/suggest")
	chatRoute.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("hello from suggest service")
	})

	authorized := chatRoute.Group("/", auth.FiberAuthMiddleware(s.Auth, s.Db.Users))
	authorized.Post("/text", s.HandleChatSuggestion())
	authorized.Post("/chat", s.HandleTextSuggestion())
}
