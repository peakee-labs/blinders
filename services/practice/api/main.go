package practiceapi

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
	practiceRoute := s.App.Group("/practice")
	practiceRoute.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("hello from practice service")
	})
	practiceRoute.Get("/unit/random", s.HandleGetRandomLanguageUnit)

	authorized := practiceRoute.Group("/", auth.FiberAuthMiddleware(s.Auth, s.Db.Users))
	authorized.Get("/unit", s.HandleSuggestLanguageUnit)
}
