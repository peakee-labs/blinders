package practiceapi

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/db/collectingdb"
	"blinders/packages/db/practicedb"
	"blinders/packages/transport"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	metadata, err := s.CollectionMetadatasRepo.GetByID(collectionOID)
	if err != nil {
		if err != mongo.ErrNoDocuments || collectionOID != userOID {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot add flashcard"})
		}

		// insert missing metadata if this is default collection
		metadata = CreateDefaultCollectionMetadata(userOID)
		_, err = s.CollectionMetadatasRepo.Insert(metadata)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot add flashcard"})
		}
	}

	flashcard, err := s.FlashCardRepo.InsertRaw(&rawFlashCard)
	if err != nil {
		log.Println("cannot insert card", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot add flashcard"})
	}

	err = s.CollectionMetadatasRepo.AddFlashCardInformation(collectionOID, flashcard.ID)
	if err != nil {
		log.Println("cannot update metadata:", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(flashcard)
}

// this functions get practice unit from collection service and create a new flashcard to review that practice unit
// this flashcard will be added to default collection
func (s Service) HandleGetPracticeFlashCard(ctx *fiber.Ctx) error {
	authUser, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panic("cannot get user auth information")
	}
	userOID, _ := primitive.ObjectIDFromHex(authUser.ID)

	// random select a practice unit type from collection service
	types := []transport.RequestType{transport.GetExplainLog, transport.GetTranslateLog}
	requestType := types[rand.Intn(len(types))]

	req := transport.GetCollectingLogRequest{
		Request: transport.Request{Type: requestType},
		Payload: transport.GetCollectingLogPayload{UserID: authUser.ID},
	}
	reqBytes, _ := json.Marshal(req)
	response, err := s.Transport.Request(
		ctx.Context(),
		s.ConsumerMap[transport.CollectingGet],
		reqBytes,
	)
	if err != nil {
		log.Printf("cannot get %v, error: %v\n", requestType, err)
		// we could return 1 flashcard from user collection
		goto returnRandomFlashCard
	}

	{
		var practiceFlashcard *practicedb.FlashCard

		switch requestType {
		case transport.GetExplainLog:
			explainLog, err := utils.ParseJSON[collectingdb.ExplainLog](response)
			if err != nil {
				log.Println("cannot parse explain log:", err)
				goto returnRandomFlashCard
			}
			// create a new flashcard from explain log, and add it to default collection
			practiceFlashcard = &practicedb.FlashCard{
				UserID:       userOID,
				FrontText:    explainLog.Request.Text,
				BackText:     explainLog.Response.Translate + "\n" + explainLog.Response.IPA,
				CollectionID: userOID,
			}
			practiceFlashcard.SetID(explainLog.ID)
			practiceFlashcard.SetInitTimeByNow()

		case transport.GetTranslateLog:
			translateLog, err := utils.ParseJSON[collectingdb.TranslateLog](response)
			if err != nil {
				log.Println("cannot parse translate log:", err)
				goto returnRandomFlashCard
			}

			// create a new flashcard from translate log, and add it to default collection
			practiceFlashcard = &practicedb.FlashCard{
				UserID:       userOID,
				FrontText:    translateLog.Request.Text,
				BackText:     translateLog.Response.Translate,
				CollectionID: userOID,
			}
			practiceFlashcard.SetID(translateLog.ID)
			practiceFlashcard.SetInitTimeByNow()
		}

		// get default collection metadata
		metadata, err := s.CollectionMetadatasRepo.GetByID(practiceFlashcard.CollectionID)
		if err != nil {
			// if default metadata is not created, create a new one
			if err == mongo.ErrNoDocuments {
				metadata = CreateDefaultCollectionMetadata(userOID)
				metadata, err = s.CollectionMetadatasRepo.Insert(metadata)
				if err != nil {
					log.Println("cannot insert metadata:", err)
				}
			}
		}

		practiceFlashcard, err = s.FlashCardRepo.Insert(practiceFlashcard)
		// check if this flashcard is already existed, if so, get latest version from db
		if mongo.IsDuplicateKeyError(err) {
			practiceFlashcard, err = s.FlashCardRepo.GetByID(practiceFlashcard.ID)
			if err != nil {
				goto returnRandomFlashCard
			}
		}

		// update metadata
		err = s.CollectionMetadatasRepo.AddFlashCardInformation(practiceFlashcard.CollectionID, practiceFlashcard.ID)
		if err != nil {
			log.Println("cannot add practice flashcard to collection metadata:", err)
		}

		return ctx.Status(http.StatusOK).JSON(practiceFlashcard)
	}

returnRandomFlashCard:
	cards, err := s.FlashCardRepo.GetFlashCardByUserID(userOID)
	if err != nil {
		log.Println("cannot get flashcard:", err)
		return ctx.Status(http.StatusOK).JSON(DefaultFlashCard)
	}
	if len(cards) == 0 {
		log.Println("user has no flashcard")
		return ctx.Status(http.StatusOK).JSON(DefaultFlashCard)
	}
	return ctx.Status(http.StatusOK).JSON(cards[rand.Intn(len(cards))])
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

	// if this flashcard is move to another collection (if collection id is updated), update metadata of old and new collection
	if updatedCard.CollectionID != card.CollectionID {
		err := s.CollectionMetadatasRepo.RemoveFlashCardInformation(card.CollectionID, cardOID)
		if err != nil {
			log.Println("cannot remove flash card information from old collection:", err)
		}

		err = s.CollectionMetadatasRepo.AddFlashCardInformation(updatedCard.CollectionID, cardOID)
		if err != nil {
			log.Println("cannot add flash card information to new collection:", err)
		}
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (s Service) HandleDeleteFlashCard(ctx *fiber.Ctx) error {
	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panic("cannot get user auth information")
	}
	userOID, _ := primitive.ObjectIDFromHex(userAuth.ID)

	cardID := ctx.Params("id")
	cardOID, err := primitive.ObjectIDFromHex(cardID)
	if err != nil {
		log.Println("invalid card id:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid card id"})
	}

	card, err := s.FlashCardRepo.GetByID(cardOID)
	if err != nil {
		log.Println("cannot get flash card:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get flash card"})
	}

	if card.UserID != userOID {
		log.Println("inefficent permission")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "inefficent permission"})
	}

	err = s.FlashCardRepo.DeleteFlashCardByID(cardOID)
	if err != nil {
		log.Println("cannot delete flash card:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot delete flash card"})
	}

	err = s.CollectionMetadatasRepo.RemoveFlashCardInformation(card.CollectionID, cardOID)
	if err != nil {
		log.Println("cannot remove flash card information from collection:", err)
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

	cards := make([]*practicedb.FlashCard, len(bodyCollection.Cards))
	for idx, card := range bodyCollection.Cards {
		cards[idx] = &practicedb.FlashCard{
			FrontText:   card.FrontText,
			FrontImgURL: card.FrontImgURL,
			BackText:    card.BackText,
			BackImgURL:  card.BackImgURL,
			UserID:      userOID,
		}
	}

	collection, err := s.FlashCardRepo.AddFlashCardsOfCollection(cards)
	if err != nil {
		log.Println("cannot add flash cards of collection:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot add flash cards of collection"})
	}

	total := make([]primitive.ObjectID, len(cards))
	for i, card := range collection.FlashCards {
		total[i] = card.ID
	}
	metadata := &practicedb.CardCollectionMetadata{
		UserID:      userOID,
		Name:        bodyCollection.Name,
		Description: bodyCollection.Description,
		Viewed:      make([]primitive.ObjectID, 0),
		Total:       total,
	}
	metadata.SetID(collection.ID)
	metadata.SetInitTimeByNow()

	metadata, err = s.CollectionMetadatasRepo.Insert(metadata)
	if err != nil {
		// TODO: find a way to rollback the collection if metadata cannot be added
		log.Println("cannot add flash cards metadata:", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(metadata)
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

	metadatas, err := s.CollectionMetadatasRepo.GetByUserID(userOID)
	if err != nil {
		log.Println("cannot get flash card collections metadata:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get flash card collections metadata"})
	}

	lookup := make(map[primitive.ObjectID]practicedb.CardCollectionMetadata)
	for _, metadata := range metadatas {
		lookup[metadata.ID] = metadata
	}

	responseCollections := make([]ResponseFlashCardCollection, len(collections))
	for idx, collection := range collections {
		responseCollections[idx] = ResponseFlashCardCollection{
			Metadata:   lookup[collection.ID],
			FlashCards: collection.FlashCards,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(responseCollections)
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

	metadata, err := s.CollectionMetadatasRepo.GetByID(collectionOID)
	if err != nil {
		log.Print("cannot get flash card collection metadata:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get flash card collection metadata"})
	}

	responseCollection := ResponseFlashCardCollection{
		Metadata:   *metadata,
		FlashCards: collection.FlashCards,
	}
	return ctx.Status(fiber.StatusOK).JSON(responseCollection)
}

func (s Service) HandleGetDefaultFlashcardCollection(ctx *fiber.Ctx) error {
	authUser, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panic("cannot get user auth information")
	}

	// default collection id will be user id
	userOID, _ := primitive.ObjectIDFromHex(authUser.ID)
	collection, err := s.FlashCardRepo.GetFlashCardCollectionByID(userOID)
	if err != nil {
		log.Println("cannot get default flash card collection:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user has no default collection"})
	}

	metadata, err := s.CollectionMetadatasRepo.GetByID(userOID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// metadata of default collection is not created
			// update default collection metadata
			total := make([]primitive.ObjectID, len(collection.FlashCards))
			for i, card := range collection.FlashCards {
				total[i] = card.ID
			}
			metadata = CreateDefaultCollectionMetadata(userOID)
			metadata, err = s.CollectionMetadatasRepo.Insert(metadata)
			if err != nil {
				log.Println("cannot insert default collection metadata:", err)
			}
		}
	}

	responseCollection := ResponseFlashCardCollection{
		Metadata:   *metadata,
		FlashCards: collection.FlashCards,
	}
	return ctx.Status(fiber.StatusOK).JSON(responseCollection)
}

func (s Service) HandleGetFlashCardCollectionsPreview(ctx *fiber.Ctx) error {
	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panicln("cannot get user auth information")
	}
	userOID, _ := primitive.ObjectIDFromHex(userAuth.ID)
	metadatas, err := s.CollectionMetadatasRepo.GetByUserID(userOID)
	if err != nil {
		log.Println("cannot get flash card collections metadata:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get flash card collections metadata"})
	}
	return ctx.Status(fiber.StatusOK).JSON(metadatas)
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

	err = s.CollectionMetadatasRepo.DeleteByID(collectionOID)
	if err != nil {
		log.Println("cannot delete flash card collection metadata:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot delete flash card collection metadata"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func CreateDefaultCollectionMetadata(userID primitive.ObjectID) *practicedb.CardCollectionMetadata {
	metadata := &practicedb.CardCollectionMetadata{
		UserID:      userID,
		Name:        "Default collection",
		Description: "This collection is automatically created for you",
		Total:       make([]primitive.ObjectID, 0),
		Viewed:      make([]primitive.ObjectID, 0),
	}
	// default collection will have ID equal to user ID
	metadata.SetID(userID)
	metadata.SetInitTimeByNow()
	return metadata
}
