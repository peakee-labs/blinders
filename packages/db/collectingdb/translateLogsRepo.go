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

type TranslateLogsRepo struct {
	dbutils.SingleCollectionRepo[*TranslateLog]
}

func NewTranslateLogsRepo(db *mongo.Database) *TranslateLogsRepo {
	return &TranslateLogsRepo{
		dbutils.SingleCollectionRepo[*TranslateLog]{
			Collection: db.Collection(TranslateLogsCollection),
		},
	}
}

func (r TranslateLogsRepo) GetLogWithSmallestGetCountByUserID(
	userID primitive.ObjectID,
) (*TranslateLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var tlog TranslateLog
	after := options.After
	err := r.FindOneAndUpdate(ctx, bson.M{
		"userId": userID,
	}, bson.M{
		"$inc": bson.M{"getCount": 1},
	}, &options.FindOneAndUpdateOptions{
		Sort:           bson.M{"getCount": 1},
		ReturnDocument: &after,
	}).Decode(&tlog)
	if err != nil {
		log.Println("can not get translate log:", err)
		return nil, fmt.Errorf("can not get translate log")
	}

	return &tlog, nil
}
