package practice

import (
	"log"

	"blinders/packages/auth"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CollectionKey = "collection"

// TODO: if you want to check if a collection is own by a user or not, use this middleware after main handler instead
func (s Service) ValidateOwnership(idParam string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		paramID := ctx.Params(idParam)
		collectionID, err := primitive.ObjectIDFromHex(paramID)
		if err != nil {
			log.Println("cannot parse collection id:", err)
			return ctx.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "cannot parse collection id"})
		}

		collection, err := s.FlashcardRepo.GetCollectionByID(collectionID)
		if err != nil {
			log.Println("cannot get flashcard collection:", err)
			return ctx.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "cannot get flashcard collection"})
		}

		userID := ctx.Locals(auth.UserIDKey).(primitive.ObjectID)

		if collection.UserID != userID {
			log.Println("user does not have permission to access this collection", err)
			return ctx.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "user does not have permission to access this collection"})
		}
		ctx.Locals(CollectionKey, collection)

		return ctx.Next()
	}
}
