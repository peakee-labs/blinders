package suggest

type ExplainPhraseInSentenceResponse struct {
	Translate         string         `json:"translate"         bson:"translate"`
	IPA               string         `json:"IPA"               bson:"IPA"`
	GrammarAnalysis   ExplainGrammar `json:"grammarAnalysis"   bson:"grammarAnalysis"`
	KeyWords          []string       `json:"keyWords"          bson:"keyWords"`
	ExpandWords       []string       `json:"expandWords"       bson:"expandWords"`
	DurationInSeconds float32        `json:"durationInSeconds" bson:"durationInSeconds"`
}

type ExplainGrammar struct {
	Tense     ExplainGrammarTense     `json:"tense"     bson:"tense"`
	Structure ExplainGrammarStructure `json:"structure" bson:"structure"`
}

type ExplainGrammarTense struct {
	Type       string `json:"type"       bson:"type"`
	Identifier string `json:"identifier" bson:"identifier"`
}

type ExplainGrammarStructure struct {
	Type      string `json:"type"      bson:"type"`
	Structure string `json:"structure" bson:"structure"`
	For       string `json:"for"       bson:"for"`
}
