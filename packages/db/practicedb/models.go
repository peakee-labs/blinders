package practicedb

import (
	dbutils "blinders/packages/db/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlashcardGenerationType string

const (
	ManualFlashcard         FlashcardGenerationType = "ManualFlashcard"
	FromExplainLogFlashcard FlashcardGenerationType = "FromExplainLogFlashcard"
)

type FlashcardCollection struct {
	dbutils.RawModel `json:",inline" bson:",inline"`
	Type             FlashcardGenerationType `json:"type"        bson:"type"`
	Name             string                  `json:"name"        bson:"name"`
	Description      string                  `json:"description" bson:"description"`
	FlashCards       []*Flashcard            `json:"flashcards"  bson:"flashcards"`
	UserID           primitive.ObjectID      `json:"userId"      bson:"userId"`
	Metadata         any                     `json:"metadata"    bson:"metadata"`
}

type Flashcard struct {
	dbutils.RawModel `       json:",inline"   bson:",inline"`
	FrontText        string `json:"frontText" bson:"frontText"`
	BackText         string `json:"backText"  bson:"backText"`
}

type ExplainLogFlashcardMetadata struct {
	ExplainLogID primitive.ObjectID `json:"explainLogId" bson:"explain_log_id"`
}
