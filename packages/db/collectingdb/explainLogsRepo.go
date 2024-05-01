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
