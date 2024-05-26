package collectingdb

import (
	dbutils "blinders/packages/db/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	TranslateLog struct {
		dbutils.RawModel `json:",inline" bson:",inline"`
		UserID           primitive.ObjectID `json:"userId"   bson:"userId"`
		Request          TranslateRequest   `json:"request"  bson:"request"`
		Response         TranslateResponse  `json:"response" bson:"response"`
		GetCount         int                `json:"getCount" bson:"getCount"`
	}
	TranslateRequest struct {
		Text string `json:"text" bson:"text"`
	}
	TranslateResponse struct {
		Translate string `json:"translate" bson:"translate"`
	}

	ExplainLog struct {
		dbutils.RawModel `json:",inline" bson:",inline"`
		UserID           primitive.ObjectID `json:"userId"   bson:"userId"`
		Request          ExplainRequest     `json:"request"  bson:"request"`
		Response         ExplainResponse    `json:"response" bson:"response"`
		GetCount         int                `json:"getCount" bson:"getCount"`
	}
	ExplainRequest struct {
		// TODO: migrate this field to "phrase"
		Text     string `json:"text"     bson:"text"`
		Sentence string `json:"sentence" bson:"sentence"`
	}
	LegacyExplainResponse struct {
		Translate         string         `json:"translate"         bson:"translate"`
		IPA               string         `json:"IPA"               bson:"IPA"`
		GrammarAnalysis   map[string]any `json:"grammarAnalysis"   bson:"grammarAnalysis"`
		ExpandWords       []string       `json:"expandWords"       bson:"expandWords"`
		KeyWords          []string       `json:"keyWords"          bson:"keyWords"`
		DurationInSeconds float32        `json:"durationInSeconds" bson:"durationInSeconds"`
	}

	ExplainResponse struct {
		Translate         string         `json:"translate"         bson:"translate"`
		IPA               string         `json:"IPA"               bson:"IPA"`
		GrammarAnalysis   ExplainGrammar `json:"grammarAnalysis"   bson:"grammarAnalysis"`
		KeyWords          []string       `json:"keyWords"          bson:"keyWords"`
		ExpandWords       []string       `json:"expandWords"       bson:"expandWords"`
		DurationInSeconds float32        `json:"durationInSeconds" bson:"durationInSeconds"`
	}

	ExplainGrammar struct {
		Tense     ExplainGrammarTense
		Structure ExplainGrammarStructure
	}

	ExplainGrammarTense struct {
		Type       string
		Identifier string
	}
	ExplainGrammarStructure struct {
		Type      string
		Structure string
		For       string
	}
)
