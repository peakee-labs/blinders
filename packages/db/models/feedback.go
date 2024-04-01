package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Feedback struct {
	UserID    primitive.ObjectID `json:"userID,omitempty" bson:"userID"`
	Comment   string             `json:"comment,omitempty" bson:"comment"`
	CreatedAt primitive.DateTime `json:"createdAt,omitempty" bson:"createdAt"`
}
