package practicedb

import (
	dbutils "blinders/packages/db/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlashCard struct {
	// ID of flashcard if from collecting event will be set to the event id, to avoid duplicate flashcard
	dbutils.RawModel `json:",inline" bson:",inline"`
	FrontText        string             `json:"frontText" bson:"frontText"`
	FrontImgURL      string             `json:"frontImageUrl"`
	BackText         string             `json:"backText" bson:"backText"`
	BackImgURL       string             `json:"backImageURL" bson:"backImageURL"`
	UserID           primitive.ObjectID `json:"userId" bson:"userId"`
	CollectionID     primitive.ObjectID `json:"collectionId" bson:"collectionId"` // by default, the collectionId is the same as the userId as the user can have 1 default collection
	Viewed           bool               `json:"count" bson:"count"`
}

type CardCollection struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	UserID     primitive.ObjectID `json:"userId" bson:"userId"`
	FlashCards []FlashCard        `json:"flashcards" bson:"flashcards"`
}

type FlashCardCollectionMetadata struct {
	dbutils.RawModel `json:",inline" bson:",inline"`
	UserID           primitive.ObjectID `json:"userId" bson:"userId"`
	Name             string             `json:"name" bson:"name"`
	Description      string             `json:"description" bson:"description"`
	Total            int                `json:"total" bson:"total"`
	Viewed           int                `json:"viewed" bson:"viewed"`
}
