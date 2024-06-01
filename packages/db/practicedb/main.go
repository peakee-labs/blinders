package practicedb

import "go.mongodb.org/mongo-driver/mongo"

const (
	FlashcardsColName = "flashcards"
)

type PracticeDB struct {
	mongo.Database
	FlashcardsRepo *FlashcardsRepo
}

func NewPracticeDB(db *mongo.Database) *PracticeDB {
	return &PracticeDB{
		Database:       *db,
		FlashcardsRepo: NewFlashcardsRepo(db),
	}
}
