package practiceapi

import (
	"blinders/packages/auth"
	"blinders/packages/collecting"
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
	Logger      collecting.EventCollector // Temporarily use, logging should run in separate service
	// Suggester   suggest.Suggester // this field is deprecated
}

func NewService(
	app *fiber.App,
	auth auth.Manager,
	db *db.MongoManager,
	logger *collecting.EventCollector,
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
	practiceRoute.Get("/unit/random", s.HandleGetRandomLanguageUnit)

	authorized := practiceRoute.Group("/", auth.FiberAuthMiddleware(s.Auth, s.Db.Users))
	authorized.Get("/unit", s.HandleSuggestLanguageUnit)
}
