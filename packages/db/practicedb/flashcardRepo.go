package practicedb

import (
	"context"

	dbutils "blinders/packages/db/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FlashcardsRepo struct {
	dbutils.SingleCollectionRepo[*FlashcardCollection]
}

func NewFlashcardsRepo(db *mongo.Database) *FlashcardsRepo {
	col := db.Collection(FlashcardsColName)
	return &FlashcardsRepo{
		SingleCollectionRepo: dbutils.SingleCollectionRepo[*FlashcardCollection]{Collection: col},
	}
}

func (r *FlashcardsRepo) AddFlashcardToCollection(
	collectionID primitive.ObjectID,
	flashcard Flashcard,
) error {
	_, err := r.Collection.UpdateByID(context.Background(), collectionID, bson.M{
		"$push": bson.M{"flashcards": flashcard},
	})
	if err != nil {
		return err
	}

	return nil
}
