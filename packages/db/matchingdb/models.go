package matchingdb

import (
	dbutils "blinders/packages/db/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchInfo struct {
	dbutils.RawModel `json:",inline" bson:",inline"`
	UserID           primitive.ObjectID `json:"userId"    bson:"userId"`
	Name             string             `json:"name"      bson:"name"`
	Gender           string             `json:"gender"    bson:"gender"`
	Major            string             `json:"major"     bson:"major"`
	Native           string             `json:"native"    bson:"native"`    // language code with RFC-5646 format
	Country          string             `json:"country"   bson:"country"`   // ISO-3166 format
	Learnings        []string           `json:"learnings" bson:"learnings"` // languages code with RFC-5646 format
	Interests        []string           `json:"interests" bson:"interests"`
	Age              int                `json:"age"       bson:"age"`
}
