package practiceapi

import (
	"log"

	"blinders/packages/auth"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CollectionKey = "collection"

// TODO: if you want to check if a collection is own by a user or not, use this middleware after main handler instead
func (s Service) CheckFlashcardCollectionOwnership(collectionParam string) fiber.Handler {
	log.Println("applying middleware to check flashcard collection ownership with param", collectionParam)
	return func(ctx *fiber.Ctx) error {
		paramID := ctx.Params(collectionParam)
		collectionID, err := primitive.ObjectIDFromHex(paramID)
		if err != nil {
			log.Println("cannot parse collection id:", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse collection id"})
		}

		collection, err := s.FlashcardRepo.GetCollectionByID(collectionID)
		if err != nil {
			log.Println("cannot get flashcard collection:", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get flashcard collection"})
		}

		userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
		if !ok {
			log.Fatalln("cannot get user auth information")
		}
		userID, _ := primitive.ObjectIDFromHex(userAuth.ID)

		if collection.UserID != userID {
			log.Println("user does not have permission to access this collection", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user does not have permission to access this collection"})
		}
		ctx.Locals(CollectionKey, collection)
		return ctx.Next()
	}
}
