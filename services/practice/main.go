package practice

import (
	"blinders/packages/auth"
	"blinders/services/practice/repo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Auth          *auth.Manager
	FlashcardRepo *repo.FlashcardsRepo
	SnapshotRepo  *repo.SnapshotsRepo
}

func NewService(auth *auth.Manager, db *mongo.Database) *Service {
	return &Service{
		Auth:          auth,
		FlashcardRepo: repo.NewFlashcardsRepo(db),
		SnapshotRepo:  repo.NewSnapshotsRepo(db),
	}
}

func (s *Service) InitFiberRoutes(r fiber.Router) {
	flashcards := r.Group("/flashcards", s.Auth.FiberAuthMiddleware(auth.Config{WithUser: true}))
	flashcardCollections := flashcards.Group("/collections")

	flashcardCollections.Get("/", s.HandleGetFlashcardCollections)
	flashcardCollections.Post("/", s.HandleCreateFlashcardCollection)

	validatedCollection := flashcardCollections.Group("/:id", s.ValidateOwnership("id"))
	validatedCollection.Get("/", s.HandleGetFlashcardCollectionByID)
	validatedCollection.Put("/", s.HandleUpdateFlashcardCollectionByID)
	validatedCollection.Delete("/", s.HandleDeleteFlashcardCollectionByID)
	validatedCollection.Post("/", s.HandleAddFlashcardToCollection)

	validatedCollection.Put("/:flashcardId/status", s.HandleUpdateFlashcardViewStatus)
	validatedCollection.Put("/:flashcardId", s.HandleUpdateFlashcardInCollection)
	validatedCollection.Delete("/:flashcardId", s.HandleRemoveFlashcardFromCollection)
}
