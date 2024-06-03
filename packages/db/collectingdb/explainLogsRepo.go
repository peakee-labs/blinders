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
	opt *Pagination,
) ([]*ExplainLog, *Pagination, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if opt == nil {
		opt = &Pagination{
			From:  time.Time{},
			To:    time.Now(),
			Limit: 0,
		}
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"userId": userID,
				"createdAt": bson.M{
					"$gt":  primitive.NewDateTimeFromTime(opt.From),
					"$lte": primitive.NewDateTimeFromTime(opt.To),
				},
			},
		},
		{
			"$sort": bson.M{"createdAt": 1},
		},
	}

	if opt.Limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": opt.Limit})
	}

	cur, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, nil, fmt.Errorf("can not get explain log")
	}
	defer cur.Close(ctx)

	logs := make([]*ExplainLog, 0)
	if err = cur.All(ctx, &logs); err != nil {
		return nil, nil, fmt.Errorf("can not get explain log")
	}

	// maybe user already fetched all logs
	if len(logs) == 0 {
		return logs, &Pagination{
			From:  opt.From,
			To:    opt.To,
			Limit: opt.Limit,
		}, nil
	}

	return logs, &Pagination{
		From:  logs[0].CreatedAt.Time(),
		To:    logs[len(logs)-1].CreatedAt.Time(),
		Limit: len(logs),
	}, nil
}

func (r ExplainLogsRepo) GetNumberOfExplainLog(
	userID primitive.ObjectID,
	from time.Time,
	to time.Time,
) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"userId": userID,
				"createdAt": bson.M{
					"$gt":  primitive.NewDateTimeFromTime(from),
					"$lte": primitive.NewDateTimeFromTime(to),
				},
			},
		},
		{
			"$count": "count",
		},
	}

	type Count struct {
		Number int `bson:"count"`
	}

	cur, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("can not get explain log number:", err)
		return 0, fmt.Errorf("can not get explain log number")
	}
	defer cur.Close(ctx)

	count := make([]Count, 0)
	if err := cur.All(ctx, &count); err != nil {
		log.Println("can not parse explain log number:", err)
		return 0, fmt.Errorf("can not parse explain log number")
	}

	if len(count) == 0 {
		log.Println("no explain log found")
		return 0, nil
	}

	return count[0].Number, nil

}
