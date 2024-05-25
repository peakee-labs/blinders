package practicedb

import (
	dbutils "blinders/packages/db/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlashCard struct {
	dbutils.RawModel `json:",inline" bson:",inline"`
	FrontText        string             `json:"frontText" bson:"frontText"`
	FrontImgURL      string             `json:"frontImageUrl"`
	BackText         string             `json:"backText" bson:"backText"`
	BackImgURL       string             `json:"backImageURL" bson:"backImageURL"`
	UserID           primitive.ObjectID `json:"userId" bson:"userId"`
	CollectionID     primitive.ObjectID `json:"collectionId" bson:"collectionId"`
}

type CardCollection struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	UserID     primitive.ObjectID `json:"userId" bson:"userId"`
	FlashCards []FlashCard        `json:"flashcards" bson:"flashcards"`
}
