package practiceapi

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/db/practicedb"
	"blinders/packages/transport"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Service) HandleGetPracticeUnitFromAnalyzeExplainLog(ctx *fiber.Ctx) error {
	authUser := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if authUser == nil {
		return fmt.Errorf("cannot get user auth information")
	}

	req := transport.GetCollectingLogRequest{
		Request: transport.Request{Type: transport.GetExplainLog},
		Payload: transport.GetCollectingLogPayload{UserID: authUser.ID},
	}

	reqBytes, _ := json.Marshal(req)
	response, err := s.Transport.Request(
		ctx.Context(),
		s.ConsumerMap[transport.CollectingGet],
		reqBytes,
	)
	if err != nil {
		log.Println("cannot get explain log:", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "cannot get explain log"})
	}

	var jsonResponse map[string]any
	_ = json.Unmarshal(response, &jsonResponse)

	return ctx.Status(http.StatusOK).JSON(jsonResponse)
}

func (s Service) HandleGetRandomLanguageUnit(ctx *fiber.Ctx) error {
	unitType := ctx.Query("type", "DEFAULT")

	switch strings.ToUpper(unitType) {
	case "DEFAULT":
		idx := rand.Intn(len(DefaultSimplePracticeUnits))
		unit := DefaultSimplePracticeUnits[idx]
		return ctx.Status(fiber.StatusOK).JSON(unit)

	case "EXPLAIN":
		idx := rand.Intn(len(ExplainLogSamples))
		unit := ExplainLogSamples[idx]
		return ctx.Status(fiber.StatusOK).JSON(unit)

	default:
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid unit type"})
	}
}

func (s Service) HandleAddFlashCard(ctx *fiber.Ctx) error {
	authUser, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panic("cannot get user auth information")
	}
	userOID, _ := primitive.ObjectIDFromHex(authUser.ID)
	bodyCard := new(RequestFlashCardBody)
	err := json.Unmarshal(ctx.Body(), &bodyCard)
	if err != nil {
		log.Println("cannot get card in request body:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get card in request body"})
	}

	if err := bodyCard.Validate(); err != nil {
		log.Println("invalid request body:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	var collectionOID primitive.ObjectID
	if bodyCard.CollectionID == "" {
		// if collectionID is empty, this flashcard will be added to default collection
		collectionOID = userOID
	} else {
		collectionOID, err = primitive.ObjectIDFromHex(bodyCard.CollectionID)
		if err != nil {
			log.Println("invalid collection id:", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid collection id"})
		}
	}

	rawFlashCard := practicedb.FlashCard{
		FrontText:    bodyCard.FrontText,
		FrontImgURL:  bodyCard.FrontImgURL,
		BackText:     bodyCard.BackText,
		BackImgURL:   bodyCard.BackImgURL,
		UserID:       userOID,
		CollectionID: collectionOID,
	}

	flashcard, err := s.FlashCardRepo.InsertRaw(&rawFlashCard)
	if err != nil {
		log.Println("cannot insert card", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot add flashcard"})
	}

	return ctx.Status(fiber.StatusOK).JSON(flashcard)
}

func (s Service) HandleGetFlashCards(ctx *fiber.Ctx) error {
	authUser, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panic("cannot get user auth information")
	}
	userOID, err := primitive.ObjectIDFromHex(authUser.ID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	cards, err := s.FlashCardRepo.GetFlashCardByUserID(userOID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get flash cards"})
	}

	return ctx.Status(fiber.StatusOK).JSON(cards)
}

func (s Service) HandleGetFlashCardByID(ctx *fiber.Ctx) error {
	cardID := ctx.Params("id")

	cardOID, err := primitive.ObjectIDFromHex(cardID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid card id"})
	}

	card, err := s.FlashCardRepo.GetByID(cardOID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get flash card"})
	}

	return ctx.Status(fiber.StatusOK).JSON(card)
}

func (s Service) HandleUpdateFlashCard(ctx *fiber.Ctx) error {
	cardID := ctx.Params("id")
	cardOID, err := primitive.ObjectIDFromHex(cardID)
	if err != nil {
		log.Println("invalid card id:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid card id"})
	}

	card, err := s.FlashCardRepo.GetByID(cardOID)
	if err != nil {
		log.Println("cannot get flash card:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "card not existed"})
	}

	authUser, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panicln("cannot get user auth information")
	}

	userOID, _ := primitive.ObjectIDFromHex(authUser.ID)
	if card.UserID != userOID {
		log.Println("inefficent permission")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "inefficent permission"})
	}

	var bodyCard RequestFlashCardBody
	if err := json.Unmarshal(ctx.Body(), &bodyCard); err != nil {
		log.Println("cannot get card in request body:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get card in request body"})
	}

	if err := bodyCard.Validate(); err != nil {
		log.Println("invalid request body:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	var collectionOID primitive.ObjectID
	if bodyCard.CollectionID == "" {
		// if collectionID is empty, this flashcard is be from default collection
		collectionOID = userOID
	} else {
		collectionOID, err = primitive.ObjectIDFromHex(bodyCard.CollectionID)
		if err != nil {
			log.Println("invalid collection id:", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid collection id"})
		}
	}

	// init practiceDbcard from request flashcardbody
	updatedCard := practicedb.FlashCard{
		UserID:       userOID,
		FrontText:    bodyCard.FrontText,
		BackText:     bodyCard.BackText,
		FrontImgURL:  bodyCard.FrontImgURL,
		BackImgURL:   bodyCard.BackImgURL,
		CollectionID: collectionOID,
	}
	updatedCard.SetID(cardOID)
	updatedCard.SetInitTime(card.CreatedAt.Time())
	updatedCard.SetUpdatedAtByNow()

	err = s.FlashCardRepo.UpdateFlashCard(cardOID, &updatedCard)
	if err != nil {
		log.Println("cannot update flash card:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot update flash card"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (s Service) HandleDeleteFlashCard(ctx *fiber.Ctx) error {
	cardID := ctx.Params("id")
	cardOID, err := primitive.ObjectIDFromHex(cardID)
	if err != nil {
		log.Println("invalid card id:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid card id"})
	}

	err = s.FlashCardRepo.DeleteFlashCardByID(cardOID)
	if err != nil {
		log.Println("cannot delete flash card:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot delete flash card"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (s Service) HandleAddFlashCardCollection(ctx *fiber.Ctx) error {
	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panic("cannot get user auth information")
	}
	userOID, _ := primitive.ObjectIDFromHex(userAuth.ID)

	var bodyCollection RequestFlashCardCollection
	if err := json.Unmarshal(ctx.Body(), &bodyCollection); err != nil {
		log.Println("cannot get collection in request body:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get collection in request body"})
	}

	if err := bodyCollection.Validate(); err != nil {
		log.Println("invalid collection", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid collection, " + err.Error()})
	}

	cards := make([]practicedb.FlashCard, 0)
	for _, card := range bodyCollection.Cards {
		cards = append(cards, practicedb.FlashCard{
			FrontText:   card.FrontText,
			FrontImgURL: card.FrontImgURL,
			BackText:    card.BackText,
			BackImgURL:  card.BackImgURL,
			UserID:      userOID,
		})
	}
	if err := s.FlashCardRepo.AddFlashCardsOfCollection(cards); err != nil {
		log.Println("cannot add flash cards of collection:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot add flash cards of collection"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (s Service) HandleGetFlashCardCollections(ctx *fiber.Ctx) error {
	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user auth information"})
	}

	userOID, _ := primitive.ObjectIDFromHex(userAuth.ID)
	collections, err := s.FlashCardRepo.GetFlashCardCollectionsByUserID(userOID)
	if err != nil {
		log.Println("cannot get flash card collections:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get flash card collections"})
	}

	return ctx.Status(fiber.StatusOK).JSON(collections)
}

func (s Service) HandleGetFlashCardCollectionByID(ctx *fiber.Ctx) error {
	collectionID := ctx.Params("id")
	collectionOID, err := primitive.ObjectIDFromHex(collectionID)
	if err != nil {
		log.Println("invalid collection id:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid collection id"})
	}

	collection, err := s.FlashCardRepo.GetFlashCardCollectionByID(collectionOID)
	if err != nil {
		log.Println("cannot get flash card collection:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get flash card collection"})
	}

	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user auth information"})
	}

	userOID, _ := primitive.ObjectIDFromHex(userAuth.ID)
	if collection.UserID != userOID {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "inefficent permission"})
	}

	return ctx.Status(fiber.StatusOK).JSON(collection)
}

func (s Service) HandleGetDefaultFlashcardCollection(ctx *fiber.Ctx) error {
	authUser, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panic("cannot get user auth information")
	}

	// default collection id will be user id
	defaultCollectionOID, _ := primitive.ObjectIDFromHex(authUser.ID)
	collection, err := s.FlashCardRepo.GetFlashCardCollectionByID(defaultCollectionOID)
	if err != nil {
		log.Println("cannot get default flash card collection:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get default flash card collection"})
	}
	return ctx.Status(fiber.StatusOK).JSON(collection)
}

func (s Service) HandleDeleteFlashCardCollection(ctx *fiber.Ctx) error {
	collectionID := ctx.Params("id")
	collectionOID, err := primitive.ObjectIDFromHex(collectionID)
	if err != nil {
		log.Println("invalid collection id")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid collection id"})
	}

	collection, err := s.FlashCardRepo.GetFlashCardCollectionByID(collectionOID)
	if err != nil {
		log.Println("cannot get flash card collection:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get flash card collection"})
	}

	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user auth information"})
	}

	userOID, _ := primitive.ObjectIDFromHex(userAuth.ID)
	if collection.UserID != userOID {
		log.Println("inefficent permission")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "inefficent permission"})
	}

	err = s.FlashCardRepo.DeleteCardCollectionByID(collectionOID)
	if err != nil {
		log.Println("cannot delete flash card collection:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot delete flash card collection"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
