package collecting

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
	EventTypeExplain             EventType = "EXPLAIN"
	LogCollection                          = "logs"
)

type EventCollector struct {
	// TODO: maybe use another db
	Col *mongo.Collection // Log collection
}

func NewEventCollector(db *mongo.Database) *EventCollector {
	return &EventCollector{
		Col: db.Collection(LogCollection),
	}
}

func (l EventCollector) AddTranslateLog(log *TranslateEventLog) (*TranslateEventLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := l.Col.InsertOne(ctx, log)
	return log, err
}

func (l EventCollector) AddRawTranslateLog(log *TranslateEventLog) (*TranslateEventLog, error) {
	log.ID = primitive.NewObjectID()
	log.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	return l.AddTranslateLog(log)
}

func (l EventCollector) GetTranslateLogByID(logID primitive.ObjectID) (*TranslateEventLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": logID}
	translateLog := new(TranslateEventLog)
	err := l.Col.FindOne(ctx, filter).Decode(translateLog)

	return translateLog, err
}

func (l EventCollector) GetTranslateLogByUserID(userID primitive.ObjectID) ([]TranslateEventLog, error) {
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

func (l EventCollector) AddSuggestPracticeLog(log *SuggestPracticeUnitEventLog) (*SuggestPracticeUnitEventLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := l.Col.InsertOne(ctx, log)
	return log, err
}

func (l EventCollector) AddRawSuggestPracticeUnitLog(log *SuggestPracticeUnitEventLog) (*SuggestPracticeUnitEventLog, error) {
	log.ID = primitive.NewObjectID()
	log.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	return l.AddSuggestPracticeLog(log)
}

func (l EventCollector) GetSuggestPracticeUnitLogByID(logID primitive.ObjectID) (*SuggestPracticeUnitEventLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": logID}
	translateLog := new(SuggestPracticeUnitEventLog)
	err := l.Col.FindOne(ctx, filter).Decode(translateLog)

	return translateLog, err
}

func (l EventCollector) GetSuggestPracticeUnitLogByUserID(userID primitive.ObjectID) ([]SuggestPracticeUnitEventLog, error) {
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

func (l EventCollector) AddExplainLog(log *ExplainEventLog) (*ExplainEventLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := l.Col.InsertOne(ctx, log)
	return log, err
}

func (l EventCollector) AddRawExplainLog(log *ExplainEventLog) (*ExplainEventLog, error) {
	log.ID = primitive.NewObjectID()
	log.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	return l.AddExplainLog(log)
}

func (l EventCollector) GetExplainLogByID(logID primitive.ObjectID) (*ExplainEventLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": logID}
	translateLog := new(ExplainEventLog)
	err := l.Col.FindOne(ctx, filter).Decode(translateLog)

	return translateLog, err
}

func (l EventCollector) GetExplainLogByUserID(userID primitive.ObjectID) ([]ExplainEventLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	filter := bson.M{"userID": userID}

	cur, err := l.Col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var res []ExplainEventLog
	if err := cur.All(ctx, res); err != nil {
		return nil, err
	}
	return res, nil
}
