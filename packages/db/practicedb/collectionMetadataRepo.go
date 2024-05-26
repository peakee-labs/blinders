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
	dbutils.SingleCollectionRepo[*CardCollectionMetadata]
}

func NewCollectionMetadataRepo(db *mongo.Database) *CollectionMetadatasRepo {
	col := db.Collection(FlashCardMetadataColName)
	return &CollectionMetadatasRepo{
		dbutils.SingleCollectionRepo[*CardCollectionMetadata]{
			Collection: col,
		},
	}
}

func (r *CollectionMetadatasRepo) GetByUserID(userID primitive.ObjectID) ([]CardCollectionMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"userId": userID}

	cur, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	collections := make([]CardCollectionMetadata, 0)
	err = cur.All(ctx, &collections)
	return collections, err
}

func (r *CollectionMetadatasRepo) Update(collectionID primitive.ObjectID, update *CardCollectionMetadata) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": collectionID}
	cur, err := r.ReplaceOne(ctx, filter, update, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}

	if cur.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *CollectionMetadatasRepo) MarkFlashCardAsViewed(collectionID, cardID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": collectionID}
	update := bson.M{"$addToSet": bson.M{"viewed": cardID}}
	_, err := r.UpdateOne(ctx, filter, update)
	return err
}

func (r *CollectionMetadatasRepo) RemoveFlashCardViewe(collectionID, cardID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": collectionID}
	update := bson.M{"$pull": bson.M{"viewed": cardID}}
	_, err := r.UpdateOne(ctx, filter, update)
	return err
}

// AddFlashCardInformation adds a flashcard to the total list of flashcards in the collection
func (r *CollectionMetadatasRepo) AddFlashCardInformation(collectionID primitive.ObjectID, cardID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": collectionID}
	update := bson.M{"$addToSet": bson.M{"total": cardID}}
	_, err := r.UpdateOne(ctx, filter, update)
	return err
}

// RemoveFlashCardInformation removes a flashcard from the total list of flashcards in the collection
func (r *CollectionMetadatasRepo) RemoveFlashCardInformation(collectionID primitive.ObjectID, cardID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": collectionID}
	// remove the card from the total list, and viewed list if existed
	update := bson.M{
		"$addToSet": bson.M{"total": cardID},
		"$pull":     bson.M{"viewed": cardID},
	}
	_, err := r.UpdateOne(ctx, filter, update)
	return err
}
