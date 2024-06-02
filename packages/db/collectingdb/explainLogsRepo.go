package collectingdb

import (
	"context"
	"fmt"
	"log"
	"time"

	dbutils "blinders/packages/db/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExplainLogsRepo struct {
	dbutils.SingleCollectionRepo[*ExplainLog]
}

func NewExplainLogsRepo(db *mongo.Database) *ExplainLogsRepo {
	return &ExplainLogsRepo{
		dbutils.SingleCollectionRepo[*ExplainLog]{
			Collection: db.Collection(ExplainLogsCollection),
		},
	}
}

func (r ExplainLogsRepo) GetLogWithSmallestGetCountByUserID(
	userID primitive.ObjectID,
) (*ExplainLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var elog ExplainLog
	after := options.After
	err := r.FindOneAndUpdate(ctx, bson.M{
		"userId": userID,
	}, bson.M{
		"$inc": bson.M{"getCount": 1},
	}, &options.FindOneAndUpdateOptions{
		Sort:           bson.M{"getCount": 1},
		ReturnDocument: &after,
	}).Decode(&elog)
	if err != nil {
		log.Println("can not get explain log:", err)
		return nil, fmt.Errorf("can not get explain log")
	}

	return &elog, nil
}

type Pagination struct {
	From  time.Time `json:"from_time"`
	To    time.Time `json:"to_time"`
	Limit int       `json:"limit"`
}

func (r ExplainLogsRepo) GetLogWithPagination(
	userID primitive.ObjectID,
	opt ...Pagination,
) ([]*ExplainLog, *Pagination, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var pipeline []bson.M

	if len(opt) > 0 {
		option := opt[0]
		pipeline = []bson.M{
			{"$match": bson.M{
				"userId": userID,
				"createdAt": bson.M{
					"$gte": primitive.NewDateTimeFromTime(option.From),
					"$lte": primitive.NewDateTimeFromTime(option.To)},
			},
			},
			{"$sort": bson.M{"createdAt": 1}},
			{"$limit": option.Limit},
		}
	} else {
		pipeline = []bson.M{
			{"$match": bson.M{"userId": userID}},
			{"$sort": bson.M{"createdAt": 1}},
		}
	}

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, nil, fmt.Errorf("can not get explain log")
	}

	var logs []*ExplainLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, nil, fmt.Errorf("can not get explain log")
	}

	return logs, &Pagination{
		From:  logs[0].CreatedAt.Time(),
		To:    logs[len(logs)-1].CreatedAt.Time(),
		Limit: len(logs),
	}, nil
}
