package practice

import (
	"blinders/services/practice/repo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	FlashcardRepo *repo.FlashcardsRepo
	SnapshotRepo  *repo.SnapshotsRepo
}

func NewService(mongoDB *mongo.Database) *Service {
	return &Service{
		FlashcardRepo: repo.NewFlashcardsRepo(mongoDB),
		SnapshotRepo:  repo.NewSnapshotsRepo(mongoDB),
	}
}

func (s *Service) InitFiberRoutes(r fiber.Router) {
	flashcards := r.Group("/flashcards")
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
