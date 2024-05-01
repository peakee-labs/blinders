package matchingdb

import "go.mongodb.org/mongo-driver/bson/primitive"

type MatchInfo struct {
	UserID    primitive.ObjectID `json:"userId"    bson:"userId,omitempty"`
	Name      string             `json:"name"      bson:"name,omitempty"`
	Gender    string             `json:"gender"    bson:"gender,omitempty"`
	Major     string             `json:"major"     bson:"major,omitempty"`
	Native    string             `json:"native"    bson:"native,omitempty"`    // language code with RFC-5646 format
	Country   string             `json:"country"   bson:"country,omitempty"`   // ISO-3166 format
	Learnings []string           `json:"learnings" bson:"learnings,omitempty"` // languages code with RFC-5646 format
	Interests []string           `json:"interests" bson:"interests,omitempty"`
	Age       int                `json:"age"       bson:"age,omitempty"`
}
