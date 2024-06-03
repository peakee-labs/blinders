package practicedb

import "go.mongodb.org/mongo-driver/mongo"

const (
	FlashcardsColName = "flashcards"
	SnapshotColName   = "snapshots"
)

type PracticeDB struct {
	mongo.Database
	FlashcardsRepo *FlashcardsRepo
	SnapshotsRepo  *SnapshotsRepo
}

func NewPracticeDB(db *mongo.Database) *PracticeDB {
	return &PracticeDB{
		Database:       *db,
		FlashcardsRepo: NewFlashcardsRepo(db),
		SnapshotsRepo:  NewSnapshotsRepo(db),
	}
}
