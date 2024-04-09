package suggestapi

import (
	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/transport"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	App         *fiber.App
	Auth        auth.Manager
	Db          *db.MongoManager
	Transport   transport.Transport
	ConsumerMap transport.ConsumerMap
	// Suggester   suggest.Suggester // this field is deprecated
}

func NewService(
	app *fiber.App,
	auth auth.Manager,
	db *db.MongoManager,
	transport transport.Transport,
	consumerMap transport.ConsumerMap,
) *Service {
	return &Service{
		App:         app,
		Auth:        auth,
		Db:          db,
		Transport:   transport,
		ConsumerMap: consumerMap,
	}
}

func (s *Service) InitRoute() {
	chatRoute := s.App.Group("/suggest")
	chatRoute.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("hello from suggest service")
	})

	authorized := chatRoute.Group("/", auth.FiberAuthMiddleware(s.Auth, s.Db.Users))

	// authorized.Post("/text", s.HandleChatSuggestion)
	// authorized.Post("/chat", s.HandleTextSuggestion)
	authorized.Get("/practice/unit", s.HandleSuggestLanguageUnit)
}
