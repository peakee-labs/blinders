package practiceapi

import (
	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/logging"
	"blinders/packages/transport"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	App         *fiber.App
	Auth        auth.Manager
	Db          *db.MongoManager
	Transport   transport.Transport
	ConsumerMap transport.ConsumerMap
	Logger      logging.EventLogger // Temporarily use, logging should run in separate service
	// Suggester   suggest.Suggester // this field is deprecated
}

func NewService(
	app *fiber.App,
	auth auth.Manager,
	db *db.MongoManager,
	logger *logging.EventLogger,
	transport transport.Transport,
	consumerMap transport.ConsumerMap,
) *Service {
	return &Service{
		App:         app,
		Auth:        auth,
		Db:          db,
		Logger:      *logger,
		Transport:   transport,
		ConsumerMap: consumerMap,
	}
}

func (s *Service) InitRoute() {
	practiceRoute := s.App.Group("/practice")
	practiceRoute.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("hello from practice service")
	})

	authorized := practiceRoute.Group("/", auth.FiberAuthMiddleware(s.Auth, s.Db.Users))
	authorized.Get("/unit", s.HandleSuggestLanguageUnit)
}
