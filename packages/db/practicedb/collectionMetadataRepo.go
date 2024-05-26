package practicedb

import (
	"context"
	"time"

	dbutils "blinders/packages/db/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionMetadatasRepo struct {
	dbutils.SingleCollectionRepo[*FlashCardCollectionMetadata]
}

func NewCollectionMetadataRepo(db *mongo.Database) *CollectionMetadatasRepo {
	col := db.Collection(FlashCardMetadataColName)
	return &CollectionMetadatasRepo{
		dbutils.SingleCollectionRepo[*FlashCardCollectionMetadata]{
			Collection: col,
		},
	}
}

func (r *CollectionMetadatasRepo) GetByUserID(userID primitive.ObjectID) ([]FlashCardCollectionMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"userId": userID}

	cur, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	collections := make([]FlashCardCollectionMetadata, 0)
	err = cur.All(ctx, &collections)
	return collections, err
}

func (r *CollectionMetadatasRepo) Update(cardID primitive.ObjectID, update *FlashCardCollectionMetadata) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": cardID}
	cur, err := r.ReplaceOne(ctx, filter, update, options.Replace().SetUpsert(false))
	if err != nil {
		return err
	}

	if cur.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *FlashCardsRepo) DeleteByID(collectionID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": collectionID}
	cur, err := r.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if cur.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
