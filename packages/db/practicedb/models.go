package practicedb

import "go.mongodb.org/mongo-driver/bson/primitive"

type FlashCard struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	FrontText    string             `json:"frontText" bson:"frontText"`
	FrontImgUrl  string             `json:"frontImageUrl"`
	BackText     string             `json:"backText" bson:"backText"`
	BackImgUrl   string             `json:"backImageURL" bson:"backImageURL"`
	UserID       primitive.ObjectID `json:"userId" bson:"userId"`
	CollectionID primitive.ObjectID `json:"collectionId" bson:"collectionId"`
}

type CardCollection struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	UserID     primitive.ObjectID `json:"userId" bson:"userId"`
	FlashCards []FlashCard        `json:"flashcards" bson:"flashcards"`
}
