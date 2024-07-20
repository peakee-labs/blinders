package practice

import (
	"encoding/json"
	"log"

	"blinders/packages/auth"
	"blinders/packages/utils"
	"blinders/services/practice/repo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Service) HandleGetFlashcardCollections(ctx *fiber.Ctx) error {
	userID := ctx.Locals(auth.UserIDKey).(primitive.ObjectID)

	var err error
	var collections []*repo.FlashcardCollection

	if ctx.Query("preview") == "true" {
		collections, err = s.FlashcardRepo.GetCollectionsMetadataByUserID(userID)
	} else {
		collections, err = s.FlashcardRepo.GetByUserID(userID)
	}

	if err != nil {
		log.Println("cannot get metadatas", err)
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot get collection preview"})
	}

	return ctx.Status(fiber.StatusOK).JSON(collections)
}

func (s Service) HandleGetFlashcardCollectionByID(ctx *fiber.Ctx) error {
	collection, ok := ctx.Locals(CollectionKey).(*repo.FlashcardCollection)
	if !ok {
		log.Fatalln("cannot get collection from context")
	}

	return ctx.Status(fiber.StatusOK).JSON(collection)
}

func (s Service) HandleCreateFlashcardCollection(ctx *fiber.Ctx) error {
	userID := ctx.Locals(auth.UserIDKey).(primitive.ObjectID)

	collection := new(repo.FlashcardCollection)
	if err := json.Unmarshal(ctx.Body(), collection); err != nil {
		log.Println("cannot unmarshal request body:", err)
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot unmarshal request body"})
	}

	collection.Type = repo.ManualCollectionType
	collection.UserID = userID

	inserted, err := s.FlashcardRepo.InsertRaw(collection)
	if err != nil {
		log.Println("cannot insert flashcard collection:", err)
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot insert flashcard collection"})
	}

	return ctx.Status(fiber.StatusOK).JSON(inserted)
}

type UpdateFlashcardBody struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Metadata    map[string]any `json:",inline,omitempty"`
}

func (s Service) HandleUpdateFlashcardCollectionByID(ctx *fiber.Ctx) error {
	collection, ok := ctx.Locals(CollectionKey).(*repo.FlashcardCollection)
	if !ok {
		log.Fatalln("cannot get collection from context")
	}

	newCollection, err := utils.ParseJSON[UpdateFlashcardBody](ctx.Body())
	if err != nil {
		log.Println("cannot unmarshal request body:", err)
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot unmarshal request body"})
	}
	updateCollection := &repo.FlashcardCollection{
		Name:        newCollection.Name,
		Description: newCollection.Description,
		Metadata:    newCollection.Metadata,
	}

	err = s.FlashcardRepo.UpdateCollectionMetadata(collection.ID, updateCollection)
	if err != nil {
		log.Println("cannot update collection metadata:", err)
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot update collection metadata"})
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (s Service) HandleDeleteFlashcardCollectionByID(ctx *fiber.Ctx) error {
	collection, ok := ctx.Locals(CollectionKey).(*repo.FlashcardCollection)
	if !ok {
		log.Fatalln("cannot get collection from context")
	}
	err := s.FlashcardRepo.DeleteByID(collection.ID)
	if err != nil {
		log.Println("cannot delete flashcard", err)
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot delete flashcard"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// define one-time used type in the usage scope
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

	collection, ok := ctx.Locals(CollectionKey).(*repo.FlashcardCollection)
	if !ok {
		log.Fatalln("cannot get collection from context")
	}

	practiceFlashcard := &repo.Flashcard{
		FrontText: cardBody.FrontText,
		BackText:  cardBody.BackText,
	}

	practiceFlashcard, err := s.FlashcardRepo.AddFlashcardToCollection(
		collection.ID,
		practiceFlashcard,
	)
	if err != nil {
		log.Println("cannot add flashcard to collection", err)
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot add flashcard to collection"})
	}

	return ctx.Status(fiber.StatusOK).JSON(practiceFlashcard)
}

func (s Service) HandleUpdateFlashcardInCollection(ctx *fiber.Ctx) error {
	cardBody := new(AddFlashcardBody)
	if err := json.Unmarshal(ctx.Body(), cardBody); err != nil {
		log.Println("invalid request body", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	collection, ok := ctx.Locals(CollectionKey).(*repo.FlashcardCollection)
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
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot add flashcard to collection"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (s Service) HandleRemoveFlashcardFromCollection(ctx *fiber.Ctx) error {
	collection, ok := ctx.Locals(CollectionKey).(*repo.FlashcardCollection)
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
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot delete flashcard"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (s Service) HandleGetCollectionsPreview(ctx *fiber.Ctx) error {
	userID := ctx.Locals(auth.UserIDKey).(primitive.ObjectID)

	metadatas, err := s.FlashcardRepo.GetCollectionsMetadataByUserID(userID)
	if err != nil {
		log.Println("cannot get metadatas", err)
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot get collection preview"})
	}

	return ctx.Status(fiber.StatusOK).JSON(metadatas)
}

func (s Service) HandleUpdateFlashcardViewStatus(ctx *fiber.Ctx) error {
	collection, ok := ctx.Locals(CollectionKey).(*repo.FlashcardCollection)
	if !ok {
		log.Fatalln("cannot get collection from context")
	}

	paramID := ctx.Params("flashcardId")
	cardID, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		log.Println("flashcardID is invalid", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "flashcardID is invalid"})
	}
	viewStatus := ctx.QueryBool("viewed", true)

	err = s.FlashcardRepo.UpdateFlashcardViewStatus(collection.ID, cardID, viewStatus)
	if err != nil {
		log.Println("cannot update flashcard view status", err)
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot update flashcard view status"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
