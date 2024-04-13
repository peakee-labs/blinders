package logging

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	EventTypeSuggestPracticeUnit EventType = "SUGGEST_LANGUAGE_UNIT"
	EventTypeTranslate           EventType = "TRANSLATE"
	LogCollection                          = "logs"
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

// GetTranslateLogByUserID returns list of translate log belongs to user with ID userID.
func (l EventLogger) GetTranslateLogByUserID(userID primitive.ObjectID) ([]TranslateEventLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	filter := bson.M{"userID": userID}

	cur, err := l.Col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var res []TranslateEventLog
	if err := cur.All(ctx, res); err != nil {
		return nil, err
	}
	return res, nil
}

// AddSuggestPracticeLog adds suggest practice unit log into db.
func (l EventLogger) AddSuggestPracticeLog(log *SuggestPracticeUnitEventLog) (*SuggestPracticeUnitEventLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := l.Col.InsertOne(ctx, log)
	return log, err
}

// AddRawSuggestPracticeUnitLog assigns unique primitive.ObjectID and primitive.DateTime to log then pass
// to EventLogger.AddSuggestPracticeLog method
func (l EventLogger) AddRawSuggestPracticeUnitLog(log *SuggestPracticeUnitEventLog) (*SuggestPracticeUnitEventLog, error) {
	log.ID = primitive.NewObjectID()
	log.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	return l.AddSuggestPracticeLog(log)
}

// GetSuggestPracticeUnitLogByID returns suggest practice unit with SuggestPracticeUnitEventLog.ID equal to params.logID in db.
func (l EventLogger) GetSuggestPracticeUnitLogByID(logID primitive.ObjectID) (*SuggestPracticeUnitEventLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": logID}
	translateLog := new(SuggestPracticeUnitEventLog)
	err := l.Col.FindOne(ctx, filter).Decode(translateLog)

	return translateLog, err
}

// GetSuggestPracticeUnitEventLogByUserID returns suggest practice unit logs belong to user with userID
func (l EventLogger) GetSuggestPracticeUnitEventLogByUserID(userID primitive.ObjectID) ([]SuggestPracticeUnitEventLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	filter := bson.M{"userId": userID}

	cur, err := l.Col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	res := make([]SuggestPracticeUnitEventLog, 0)
	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}
