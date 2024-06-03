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
	SnapshotRepo  *practicedb.SnapshotsRepo
}

func NewService(
	app *fiber.App,
	auth auth.Manager,
	usersRepo *usersdb.UsersRepo,
	flashcardsRepo *practicedb.FlashcardsRepo,
	snapshotRepo *practicedb.SnapshotsRepo,
	transport transport.Transport,
) *Service {
	return &Service{
		App:           app,
		Auth:          auth,
		UserRepo:      usersRepo,
		FlashcardRepo: flashcardsRepo,
		SnapshotRepo:  snapshotRepo,
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
	flashcardCollections.Get("/preview", s.HandleGetCollectionsPreview)
	flashcardCollections.Post("/", s.HandleCreateFlashcardCollection)

	validatedCollections := flashcardCollections.Group("/:id", s.CheckFlashcardCollectionOwnership("id"))
	validatedCollections.Get("/", s.HandleGetFlashcardCollectionByID)
	validatedCollections.Put("/", s.HandleUpdateFlashcardCollectionByID)
	validatedCollections.Delete("/", s.HandleDeleteFlashcardCollectionByID)

	validatedCollections.Post("/", s.HandleAddFlashcardToCollection)
	validatedCollections.Put("/:flashcardId/status", s.HandleUpdateFlashcardViewStatus)
	validatedCollections.Put("/:flashcardId", s.HandleUpdateFlashcardInCollection)
	validatedCollections.Delete("/:flashcardId", s.HandleRemoveFlashcardFromCollection)

	explainLog := authorized.Group("/explain-log")
	explainLog.Get("/", s.HandleFetchExplainMetadata)
	explainLog.Get("/flashcard", s.HandleCreateFlashcardFromExplainLog)
}
