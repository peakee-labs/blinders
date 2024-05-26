package practiceapi

import (
	"blinders/packages/auth"
	"blinders/packages/db/practicedb"
	"blinders/packages/db/usersdb"
	"blinders/packages/transport"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	App                     *fiber.App
	Auth                    auth.Manager
	UserRepo                *usersdb.UsersRepo
	Transport               transport.Transport
	ConsumerMap             transport.ConsumerMap
	FlashCardRepo           *practicedb.FlashCardsRepo
	CollectionMetadatasRepo *practicedb.CollectionMetadatasRepo
}

func NewService(
	app *fiber.App,
	auth auth.Manager,
	usersRepo *usersdb.UsersRepo,
	flashCardsRepo *practicedb.FlashCardsRepo,
	collectionMetadatasRepo *practicedb.CollectionMetadatasRepo,
	transport transport.Transport,
	consumerMap transport.ConsumerMap,
) *Service {
	return &Service{
		App:                     app,
		Auth:                    auth,
		UserRepo:                usersRepo,
		FlashCardRepo:           flashCardsRepo,
		CollectionMetadatasRepo: collectionMetadatasRepo,
		Transport:               transport,
		ConsumerMap:             consumerMap,
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
	authorized.Get("/unit/flashcard", s.HandleGetPracticeFlashCard)

	authorized.Get("/flashcards/collections", s.HandleGetFlashCardCollections)
	authorized.Get("/flashcards/collections/default", s.HandleGetDefaultFlashcardCollection)
	authorized.Get("/flashcards/collections/preview", s.HandleGetFlashCardCollectionsPreview)
	authorized.Post("/flashcards/collections", s.HandleAddFlashCardCollection)
	authorized.Get("/flashcards/collections/:id", s.HandleGetFlashCardCollectionByID)
	// TODO: view status of collection APIs
	authorized.Delete("/flashcards/collections/:id", s.HandleDeleteFlashCardCollection)

	authorized.Get("/flashcards/:id", s.HandleGetFlashCardByID)
	authorized.Put("/flashcards/:id", s.HandleUpdateFlashCard)
	authorized.Delete("/flashcards/:id", s.HandleDeleteFlashCard)

	authorized.Get("/flashcards", s.HandleGetFlashCards)
	authorized.Post("/flashcards", s.HandleAddFlashCard)
}
