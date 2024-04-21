package collecting

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	EventType    string
	GenericEvent struct {
		Type    EventType `json:"type"`
		Payload any       `json:"event"` //  event struct
	}
	TranslateRequest struct {
		Text string `json:"text" bson:"text"`
	}
	TranslateResponse struct {
		Translate string `json:"translate" bson:"translate"`
	}
	TranslateEvent struct {
		UserID   primitive.ObjectID `json:"userId" bson:"userId"`
		Request  TranslateRequest   `json:"request" bson:"request"`
		Response TranslateResponse  `json:"response" bson:"response"`
	}
	TranslateEventLog struct {
		TranslateEvent `json:",inline" bson:",inline"`
		ID             primitive.ObjectID `json:"logId" bson:"_id"`
		CreatedAt      primitive.DateTime `json:"createdAt" bson:"createdAt"`
	}
	ExplainRequest struct {
		Text     string `json:"text" bson:"text"`
		Sentence string `json:"sentence" bson:"sentence"`
	}
	ExplainResponse struct {
		Translate       string         `json:"translate" bson:"translate"`
		GrammarAnalysis map[string]any `json:"grammarAnalysis" bson:"grammarAnalysis"`
		ExpandWords     []string       `json:"expandWords" bson:"expandWords"`
	}
	ExplainEvent struct {
		UserID   primitive.ObjectID `json:"userId" bson:"userId"`
		Request  ExplainRequest     `json:"request" bson:"request"`
		Response ExplainResponse    `json:"response" bson:"response"`
	}
	ExplainEventLog struct {
		ExplainEvent `json:",inline" bson:",inline"`
		ID           primitive.ObjectID `json:"logId" bson:"_id"`
		CreatedAt    primitive.DateTime `json:"createdAt" bson:"createdAt"`
	}
	SuggestPracticeUnitRequest struct {
		Text    string `json:"text" bson:"text"`
		Context string `json:"context" bson:"context"`
	}
	SuggestPracticeUnitResponse struct {
		Word        string   `json:"word" bson:"word"`
		Explain     string   `json:"explain" bson:"explain"`
		ExpandWords []string `json:"expandWords" bson:"expandWords"`
	}
	SuggestPracticeUnitEvent struct {
		UserID   primitive.ObjectID          `json:"userId" bson:"userId"`
		Request  SuggestPracticeUnitRequest  `json:"request" bson:"request"`
		Response SuggestPracticeUnitResponse `json:"response" bson:"response"`
	}
	SuggestPracticeUnitEventLog struct {
		ID                       primitive.ObjectID `json:"logId" bson:"_id"`
		CreatedAt                primitive.DateTime `json:"createdAt" bson:"createdAt"`
		SuggestPracticeUnitEvent `json:",inline" bson:",inline"`
	}
)

func NewGenericEvent(EventType EventType, Event any) GenericEvent {
	return GenericEvent{
		Type:    EventType,
		Payload: Event,
	}
}
