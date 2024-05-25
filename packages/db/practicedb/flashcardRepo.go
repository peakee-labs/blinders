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

type FlashCardsRepo struct {
	dbutils.SingleCollectionRepo[*FlashCard]
}

func NewFlashCardRepo(db *mongo.Database) *FlashCardsRepo {
	col := db.Collection(FlashCardColName)
	return &FlashCardsRepo{SingleCollectionRepo: dbutils.SingleCollectionRepo[*FlashCard]{Collection: col}}
}

func (r *FlashCardsRepo) GetFlashCardByUserID(userID primitive.ObjectID) ([]FlashCard, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"userId": userID}

	cur, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	cards := make([]FlashCard, 0)
	err = cur.All(ctx, &cards)
	return cards, err
}

func (r *FlashCardsRepo) GetFlashCardsByCollectionID(collectionID primitive.ObjectID) ([]FlashCard, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"collectionId": collectionID}
	cur, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	cards := make([]FlashCard, 0)
	err = cur.All(ctx, &cards)
	return cards, err
}

func (r *FlashCardsRepo) GetFlashCardCollectionByID(collectionID primitive.ObjectID) (*CardCollection, error) {
	cards, err := r.GetFlashCardsByCollectionID(collectionID)
	if err != nil {
		return nil, err
	}
	if len(cards) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &CardCollection{
		ID:         collectionID,
		FlashCards: cards,
	}, err
}

func (r *FlashCardsRepo) UpdateFlashCard(cardID primitive.ObjectID, newCard *FlashCard) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": cardID}
	cur, err := r.ReplaceOne(ctx, filter, newCard, options.Replace().SetUpsert(false))
	if err != nil {
		return err
	}

	if cur.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *FlashCardsRepo) DeleteFlashCardByID(cardID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": cardID}
	cur, err := r.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if cur.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *FlashCardsRepo) GetFlashCardCollectionsByUserID(userID primitive.ObjectID) ([]*CardCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pipeline := []bson.M{
		{
			"$match": bson.M{"userId": userID},
		},
		{
			"$group": bson.M{
				"_id":        "$collectionId",
				"flashcards": bson.M{"$push": "$$ROOT"},
			},
		},
		{
			"$addFields": bson.M{
				"userId": userID,
			},
		},
	}
	cur, err := r.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	collections := make([]*CardCollection, 0)
	err = cur.All(ctx, &collections)
	return collections, err
}

func (r *FlashCardsRepo) DeleteCardCollectionByID(collectionID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"collectionId": collectionID}
	cur, err := r.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	if cur.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
