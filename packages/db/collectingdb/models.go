package collectingdb

import (
	dbutils "blinders/packages/db/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	TranslateRequest struct {
		Text string `json:"text" bson:"text"`
	}
	TranslateResponse struct {
		Translate string `json:"translate" bson:"translate"`
	}
	TranslateLog struct {
		dbutils.RawModel `json:",inline" bson:",inline"`
		UserID           primitive.ObjectID `json:"userId"   bson:"userId"`
		Request          TranslateRequest   `json:"request"  bson:"request"`
		Response         TranslateResponse  `json:"response" bson:"response"`
		GetCount         int                `json:"getCount" bson:"getCount"`
	}

	ExplainRequest struct {
		Text     string `json:"text"     bson:"text"`
		Sentence string `json:"sentence" bson:"sentence"`
	}
	ExplainResponse struct {
		Translate       string         `json:"translate"       bson:"translate"`
		GrammarAnalysis map[string]any `json:"grammarAnalysis" bson:"grammarAnalysis"`
		ExpandWords     []string       `json:"expandWords"     bson:"expandWords"`
		KeyWords        []string       `json:"keyWords"        bson:"keyWords"`
	}
	ExplainLog struct {
		dbutils.RawModel `json:",inline" bson:",inline"`
		UserID           primitive.ObjectID `json:"userId"   bson:"userId"`
		Request          ExplainRequest     `json:"request"  bson:"request"`
		Response         ExplainResponse    `json:"response" bson:"response"`
		GetCount         int                `json:"getCount" bson:"getCount"`
	}
)
