package practicedb

import (
	dbutils "blinders/packages/db/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CollectionType string
	SnapshotType   string
	FlashcardType  string
)

const (
	ManualCollectionType         CollectionType = "ManualCollection"
	FromExplainLogCollectionType CollectionType = "FromExplainLogCollection"
	DefaultCollectionType        CollectionType = "DefaultCollection"

	ExplainLogToFlashcardSnapshotType SnapshotType = "ExplainLogToFlashcardSnapshot"

	ExplainLogToFlashcardType FlashcardType = "ExplainLogFlashcard"
	ManualFlashcardType       FlashcardType = "ManualFlashcard"
	DefaultFlashcardType      FlashcardType = "ManualFlashcard"
)

type FlashcardCollection struct {
	dbutils.RawModel `                   json:",inline"                bson:",inline"`
	Type             CollectionType       `json:"type"                 bson:"type"`
	Name             string               `json:"name"                 bson:"name"`
	Description      string               `json:"description"          bson:"description"`
	Viewed           []primitive.ObjectID `json:"viewed"               bson:"viewed"`
	Total            []primitive.ObjectID `json:"total"                bson:"total"`
	UserID           primitive.ObjectID   `json:"userId"               bson:"userId"`
	FlashCards       *[]*Flashcard        `json:"flashcards,omitempty" bson:"flashcards"`
	Metadata         map[string]any       `json:"metadata,omitempty"   bson:"metadata,omitempty"`
}

type Flashcard struct {
	dbutils.RawModel `               json:",inline"           bson:",inline"`
	Type             FlashcardType `json:"type"               bson:"type"`
	FrontText        string        `json:"frontText"          bson:"frontText"`
	BackText         string        `json:"backText"           bson:"backText"`
	Metadata         any           `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

type ExplainLogFlashcardMetadata struct {
	ExplainLogID primitive.ObjectID `json:"explainLogId" bson:"explain_log_id"`
}

type PracticeSnapshot struct {
	dbutils.RawModel `json:",inline" bson:",inline"`
	Type             SnapshotType       `json:"type" bson:"type"`
	UserID           primitive.ObjectID `json:"userId" bson:"userId"`
	Current          primitive.DateTime `json:"current" bson:"current"`
}
