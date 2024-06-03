package practicedb

import (
	"context"
	"time"

	dbutils "blinders/packages/db/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SnapshotsRepo struct {
	dbutils.SingleCollectionRepo[*PracticeSnapshot]
}

func NewSnapshotsRepo(db *mongo.Database) *SnapshotsRepo {
	col := db.Collection(SnapshotColName)
	return &SnapshotsRepo{
		SingleCollectionRepo: dbutils.SingleCollectionRepo[*PracticeSnapshot]{Collection: col},
	}
}

func (r SnapshotsRepo) GetSnapshotOfUserByType(
	userID primitive.ObjectID,
	typ SnapshotType,
) (*PracticeSnapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{
		"userId": userID,
		"type":   typ,
	}

	var snapshot *PracticeSnapshot
	if err := r.FindOne(ctx, filter).Decode(&snapshot); err != nil {
		return nil, err
	}

	return snapshot, nil
}

func (r SnapshotsRepo) UpdateSnapshot(updateSnapshot *PracticeSnapshot) (*PracticeSnapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	updateSnapshot.SetUpdatedAtByNow()

	filter := bson.M{
		"_id": updateSnapshot.ID,
	}
	update := bson.M{
		"$set": bson.M{
			"current":   updateSnapshot.Current,
			"updatedAt": updateSnapshot.UpdatedAt,
		},
	}

	cur, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if cur.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return updateSnapshot, nil
}
