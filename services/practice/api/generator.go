package practiceapi

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"blinders/packages/auth"
	"blinders/packages/db/collectingdb"
	"blinders/packages/db/practicedb"
	"blinders/packages/transport"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s Service) HandleSyncExplainLogs(ctx *fiber.Ctx) error {
	limit := ctx.QueryInt("limit", 15)
	if limit == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid limit"})
	}

	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Fatalln("cannot get user auth information")
	}
	userID, _ := primitive.ObjectIDFromHex(userAuth.ID)
	var snapshotTime time.Time

	snapshot, err := s.SnapshotRepo.GetSnapshotOfUserByType(
		userID,
		practicedb.ExplainLogToFlashcardSnapshotType,
	)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Println("cannot get snapshot:", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get snapshot"})
		}

		snapshot = &practicedb.PracticeSnapshot{
			Type:    practicedb.ExplainLogToFlashcardSnapshotType,
			UserID:  userID,
			Current: primitive.NewDateTimeFromTime(time.Time{}),
		}

		snapshot, err = s.SnapshotRepo.InsertRaw(snapshot)
		if err != nil {
			log.Println("cannot insert snapshot:", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot insert snapshot"})
		}
	}

	snapshotTime = snapshot.Current.Time()

	request := transport.GetCollectingLogRequest{
		Request: transport.Request{Type: transport.GetExplainLogBatch},
		Payload: transport.GetCollectingLogPayload{
			UserID: userAuth.ID,
			PagintionInfo: &collectingdb.Pagination{
				From:  snapshotTime,
				To:    time.Now(),
				Limit: limit,
			},
		},
	}

	reqBytes, _ := json.Marshal(request)
	response, err := s.Transport.Request(
		ctx.Context(),
		s.Transport.ConsumerID(transport.CollectingGet),
		reqBytes,
	)
	if err != nil {
		log.Println("cannot get explain log batch:", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "cannot sync explain logs"})
	}

	logsResponse, err := utils.ParseJSON[transport.GetExplainLogBatchResponse](response)
	if err != nil {
		log.Println("cannot parse explain log batch:", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "cannot sync explain logs"})
	}

	if len(logsResponse.Logs) == 0 {
		return ctx.Status(http.StatusOK).JSON(
			practicedb.FlashcardCollection{
				Type:        practicedb.FromExplainLogCollectionType,
				Name:        "From Explain Log",
				Description: "Flashcards generated from explain logs",
				UserID:      userID,
				FlashCards:  &[]*practicedb.Flashcard{},
			},
		)
	}

	// generate flashcard collection from logs
	flashcards := make([]*practicedb.Flashcard, len(logsResponse.Logs))
	for i, log := range logsResponse.Logs {
		flashcards[i] = &practicedb.Flashcard{
			Type:      practicedb.ExplainLogToFlashcardType,
			FrontText: log.Request.Text,
			BackText:  log.Response.Translate,
			Metadata: &practicedb.ExplainLogFlashcardMetadata{
				ExplainLogID: log.ID,
			},
		}
	}

	collection := &practicedb.FlashcardCollection{
		Type:        practicedb.FromExplainLogCollectionType,
		Name:        "From Explain Log",
		Description: "Flashcards generated from explain logs",
		UserID:      userID,
		FlashCards:  &flashcards,
	}

	insertedCollection, err := s.FlashcardRepo.InsertRaw(collection)
	if err != nil {
		log.Println("cannot insert flashcard collection:", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "cannot insert flashcard collection"})
	}

	go func() {
		snapshot.Current = primitive.NewDateTimeFromTime(logsResponse.PagintionInfo.To)

		_, err := s.SnapshotRepo.UpdateSnapshot(snapshot)
		if err != nil {
			log.Println("cannot insert snapshot:", err)
		}
	}()

	return ctx.Status(http.StatusOK).JSON(insertedCollection)
}
