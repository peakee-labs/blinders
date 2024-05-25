package practiceapi

import (
	"blinders/packages/auth"
	"blinders/packages/db/practicedb"
	"blinders/packages/db/usersdb"
	"blinders/packages/transport"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	App           *fiber.App
	Auth          auth.Manager
	UserRepo      *usersdb.UsersRepo
	Transport     transport.Transport
	ConsumerMap   transport.ConsumerMap
	FlashCardRepo *practicedb.FlashCardsRepo
}

func NewService(
	app *fiber.App,
	auth auth.Manager,
	usersRepo *usersdb.UsersRepo,
	transport transport.Transport,
	consumerMap transport.ConsumerMap,
) *Service {
	return &Service{
		App:         app,
		Auth:        auth,
		UserRepo:    usersRepo,
		Transport:   transport,
		ConsumerMap: consumerMap,
	}
}

func (s *Service) InitRoute() {
	practiceRoute := s.App.Group("/practice")
	practiceRoute.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("hello from practice service")
	})
	practiceRoute.Get("/public/unit", s.HandleGetRandomLanguageUnit)
	authorized := practiceRoute.Group("/", auth.FiberAuthMiddleware(s.Auth, s.UserRepo))
	authorized.Get("/unit", s.HandleGetPracticeUnitFromAnalyzeExplainLog)

	authorized.Get("/flashcards", s.HandleGetFlashCards)
	authorized.Post("/flashcards", s.HandleGetFlashCards)

	authorized.Get("/flashcards/:id", s.HandleGetFlashCardByID)
	authorized.Put("/flashcards/:id", s.HandleUpdateFlashCard)
	authorized.Delete("/flashcards/:id", s.HandleDeleteFlashCard)

	authorized.Get("/flashcard/collections", s.HandleGetFlashCardCollections)
	authorized.Get("/flashcard/collections/:id", s.HandleGetFlashCardCollectionByID)
	authorized.Delete("/flashcard/collections/:id", s.handleDeleteFlashCardCollection)
}
