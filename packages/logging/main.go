package logging

import (
	"context"
	"time"

	"blinders/packages/transport"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	EventTypeSuggestLanguageUnit transport.EventType = "SUGGEST_LANGUAGE_UNIT"
	LogCollection                                    = "logs"
)

type (
	LanguageUnitRequest struct {
		Text     string `json:"text" bson:"text"`
		Sentence string `json:"sentence" bson:"sentence"`
	} // TODO: we could put this struct in another pkg
	LanguageUnitResponse struct {
		Translate       string         `json:"translate" bson:"translate"`
		GrammarAnalysis map[string]any `json:"grammarAnalysis" bson:"grammarAnalysis"`
		ExpandWords     []string       `json:"expandWords" bson:"expandWords"`
	} // TODO: we could put this struct in another pkg
	SuggestLanguageUnitEvent struct {
		UserID   string               `json:"userId" bson:"userId"`
		Request  LanguageUnitRequest  `json:"request" bson:"request"`
		Response LanguageUnitResponse `json:"response" bson:"response"`
	} // TODO: we could put this struct in another pkg
	TranslateEventLog struct {
		SuggestLanguageUnitEvent `json:",inline" bson:",inline"`
		ID                       primitive.ObjectID `json:"logId" bson:"_id"`
		CreatedAt                primitive.DateTime `json:"createdAt" bson:"createdAt"`
	}
)

type EventLogger struct {
	// TODO: maybe use another db
	Col *mongo.Collection // Log collection
}

func NewEventLogger(db *mongo.Database) *EventLogger {
	return &EventLogger{
		Col: db.Collection(LogCollection),
	}
}

// AddTranslateLog adds translate log into db.
func (l EventLogger) AddTranslateLog(log *TranslateEventLog) (*TranslateEventLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := l.Col.InsertOne(ctx, log)
	return log, err
}

// AddRawTranslateLog assigns unique primitive.ObjectID and primitive.DateTime to log then pass
// to EventLogger.AddTranslateLog method
func (l EventLogger) AddRawTranslateLog(log *TranslateEventLog) (*TranslateEventLog, error) {
	log.ID = primitive.NewObjectID()
	log.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	return l.AddTranslateLog(log)
}

// GetTranslateLogByID returns translate log with TranslateEventLog.ID equal to params.logID in db.
func (l EventLogger) GetTranslateLogByID(logID primitive.ObjectID) (*TranslateEventLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": logID}
	translateLog := new(TranslateEventLog)
	err := l.Col.FindOne(ctx, filter).Decode(translateLog)

	return translateLog, err
}
