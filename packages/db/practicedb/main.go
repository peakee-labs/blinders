package practicedb

import "go.mongodb.org/mongo-driver/mongo"

const (
	FlashCardColName = "flashcards"
)

type PracticeDB struct {
	mongo.Database
	FlashCardRepo *FlashCardsRepo
}

func NewPracticeDB(db *mongo.Database) *PracticeDB {
	return &PracticeDB{
		Database:      *db,
		FlashCardRepo: NewFlashCardRepo(db),
	}
}
