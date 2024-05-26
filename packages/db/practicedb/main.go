package practicedb

import "go.mongodb.org/mongo-driver/mongo"

const (
	FlashCardColName         = "flashcards"
	FlashCardMetadataColName = "flashcard-metadata"
)

type PracticeDB struct {
	mongo.Database
	FlashCardRepo          *FlashCardsRepo
	CollectionMetadataRepo *CollectionMetadatasRepo
}

func NewPracticeDB(db *mongo.Database) *PracticeDB {
	return &PracticeDB{
		Database:               *db,
		FlashCardRepo:          NewFlashCardRepo(db),
		CollectionMetadataRepo: NewCollectionMetadataRepo(db),
	}
}
