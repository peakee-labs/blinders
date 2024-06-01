package practiceapi

import (
	"encoding/json"
	"log"

	"blinders/packages/auth"
	"blinders/packages/db/practicedb"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s Service) HandleGetFlashcardCollections(ctx *fiber.Ctx) error {
	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Fatalln("cannot get user auth information")
	}
	userID, _ := primitive.ObjectIDFromHex(userAuth.ID)
	collections, err := s.FlashcardRepo.GetByUserID(userID)
	if err != nil {
		log.Println("cannot get flashcard collections:", err)
		if err != mongo.ErrNoDocuments {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get flashcard collections"})
		}
		collections = []*practicedb.FlashcardCollection{}
	}

	return ctx.Status(fiber.StatusOK).JSON(collections)
}

func (s Service) HandleGetOrCreateDefaultFlashcardCollection(ctx *fiber.Ctx) error {
	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Fatalln("cannot get user auth information")
	}
	userID, _ := primitive.ObjectIDFromHex(userAuth.ID)

	collection, err := s.FlashcardRepo.GetByID(userID)
	if err != nil {
		log.Println("cannot get flashcard collections:", err)

		collection = &practicedb.FlashcardCollection{
			CollectionMetadata: practicedb.CollectionMetadata{
				Type:   practicedb.DefaultFlashcard,
				Name:   "Default Collection",
				UserID: userID,
			},
			FlashCards: []*practicedb.Flashcard{},
		}
		collection.SetID(userID)
		collection.SetInitTimeByNow()
		collection, err = s.FlashcardRepo.Insert(collection)
		if err != nil {
			log.Println("cannot create default flashcard collection:", err)
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"error": "cannot create default flashcard collection"})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(collection)
}

func (s Service) HandleGetFlashcardCollectionByID(ctx *fiber.Ctx) error {
	collection, ok := ctx.Locals(CollectionKey).(*practicedb.FlashcardCollection)
	if !ok {
		log.Fatalln("cannot get collection from context")
	}

	return ctx.Status(fiber.StatusOK).JSON(collection)
}

func (s Service) HandleCreateFlashcardCollection(ctx *fiber.Ctx) error {
	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Fatalln("cannot get user auth information")
	}
	userID, _ := primitive.ObjectIDFromHex(userAuth.ID)

	collection := new(practicedb.FlashcardCollection)
	if err := json.Unmarshal(ctx.Body(), collection); err != nil {
		log.Println("cannot unmarshal request body:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot unmarshal request body"})
	}

	collection.Type = practicedb.ManualFlashcard
	collection.UserID = userID

	inserted, err := s.FlashcardRepo.InsertRaw(collection)
	if err != nil {
		log.Println("cannot insert flashcard collection:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot insert flashcard collection"})
	}

	return ctx.Status(fiber.StatusOK).JSON(inserted)
}

func (s Service) HandleUpdateFlashcardCollectionByID(ctx *fiber.Ctx) error {
	collection, ok := ctx.Locals(CollectionKey).(*practicedb.FlashcardCollection)
	if !ok {
		log.Fatalln("cannot get collection from context")
	}

	newCollection := new(practicedb.CollectionMetadata)
	if err := json.Unmarshal(ctx.Body(), newCollection); err != nil {
		log.Println("cannot unmarshal request body:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot unmarshal request body"})
	}

	err := s.FlashcardRepo.UpdateCollectionMetadata(collection.ID, newCollection)
	if err != nil {
		log.Println("cannot update collection metadata:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot update collection metadata"})
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (s Service) HandleDeleteFlashcardCollectionByID(ctx *fiber.Ctx) error {
	collection, ok := ctx.Locals(CollectionKey).(*practicedb.FlashcardCollection)
	if !ok {
		log.Fatalln("cannot get collection from context")
	}
	err := s.FlashcardRepo.DeleteByID(collection.ID)
	if err != nil {
		log.Println("cannot delete flashcard", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot delete flashcard"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// define one-time used type in the usage scope
type (
	AddFlashcardResponse struct{}
)

type AddFlashcardBody struct {
	FrontText string
	BackText  string
}

func (s Service) HandleAddFlashcardToCollection(ctx *fiber.Ctx) error {
	cardBody := new(AddFlashcardBody)
	if err := json.Unmarshal(ctx.Body(), cardBody); err != nil {
		log.Println("invalid request body", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	collection, ok := ctx.Locals(CollectionKey).(*practicedb.FlashcardCollection)
	if !ok {
		log.Fatalln("cannot get collection from context")
	}

	practiceFlashcard := &practicedb.Flashcard{
		FrontText: cardBody.FrontText,
		BackText:  cardBody.BackText,
	}

	practiceFlashcard, err := s.FlashcardRepo.AddFlashcardToCollection(collection.ID, practiceFlashcard)
	if err != nil {
		log.Println("cannot add flashcard to collection", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot add flashcard to collection"})
	}

	return ctx.Status(fiber.StatusOK).JSON(practiceFlashcard)
}

func (s Service) HandleUpdateFlashcardInCollection(ctx *fiber.Ctx) error {
	cardBody := new(AddFlashcardBody)
	if err := json.Unmarshal(ctx.Body(), cardBody); err != nil {
		log.Println("invalid request body", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	collection, ok := ctx.Locals(CollectionKey).(*practicedb.FlashcardCollection)
	if !ok {
		log.Fatalln("cannot get collection from context")
	}

	flashcardParamID := ctx.Params("flashcardId")
	flashcardID, err := primitive.ObjectIDFromHex(flashcardParamID)
	if err != nil {
		log.Println("flashcardID is invalid", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "flashcardID is invalid"})
	}

	flashcard, err := s.FlashcardRepo.GetFlashcardByID(collection.ID, flashcardID)
	if err != nil {
		log.Println("cannot get flashcard", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get flashcard"})
	}
	flashcard.FrontText = cardBody.FrontText
	flashcard.BackText = cardBody.BackText

	err = s.FlashcardRepo.UpdateFlashCard(collection.ID, *flashcard)
	if err != nil {
		log.Println("cannot add flashcard to collection", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot add flashcard to collection"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (s Service) HandleRemoveFlashcardFromCollection(ctx *fiber.Ctx) error {
	collection, ok := ctx.Locals(CollectionKey).(*practicedb.FlashcardCollection)
	if !ok {
		log.Fatalln("cannot get collection from context")
	}

	flashcardParamID := ctx.Params("flashcardId")
	flashcardID, err := primitive.ObjectIDFromHex(flashcardParamID)
	if err != nil {
		log.Println("flashcardID is invalid", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "flashcardID is invalid"})
	}

	if err := s.FlashcardRepo.DeleteFlashCard(collection.ID, flashcardID); err != nil {
		log.Println("cannot delete flashcard", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot delete flashcard"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
