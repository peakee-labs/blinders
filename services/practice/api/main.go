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
	FlashcardRepo *practicedb.FlashcardsRepo
}

func NewService(
	app *fiber.App,
	auth auth.Manager,
	usersRepo *usersdb.UsersRepo,
	flashcardsRepo *practicedb.FlashcardsRepo,
	transport transport.Transport,
) *Service {
	return &Service{
		App:           app,
		Auth:          auth,
		UserRepo:      usersRepo,
		FlashcardRepo: flashcardsRepo,
		Transport:     transport,
	}
}

func (s *Service) InitRoute() {
	practiceRoute := s.App.Group("/practice")
	practiceRoute.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("hello from practice service")
	})

	authorized := practiceRoute.Group("/", auth.FiberAuthMiddleware(s.Auth, s.UserRepo))
	authorized.Get("/random-review", s.HandleGetRandomReview)
	authorized.Get("/fast-review", s.HandleGetFastReviewFromExplainLog)

	flashcards := authorized.Group("/flashcards")
	flashcardCollections := flashcards.Group("/collections")

	flashcardCollections.Get("/", s.HandleGetFlashcardCollections)
	flashcardCollections.Get("/default", s.HandleGetOrCreateDefaultFlashcardCollection)
	flashcardCollections.Post("/", s.HandleCreateFlashcardCollection)

	validatedCollections := flashcardCollections.Group("/:id", CheckFlashcardCollectionOwnership(s, "id"))
	validatedCollections.Get("/", s.HandleGetFlashcardCollectionByID)
	validatedCollections.Put("/", s.HandleUpdateFlashcardCollectionByID)
	validatedCollections.Delete("/", s.HandleDeleteFlashcardCollectionByID)

	validatedCollections.Post("/", s.HandleAddFlashcardToCollection)
	validatedCollections.Put("/:flashcardId", s.HandleUpdateFlashcardInCollection)
	validatedCollections.Delete("/:flashcardId", s.HandleRemoveFlashcardFromCollection)

	flashcards.Get("/sync-explain-logs", s.HandleSyncExplainLogs)
}
