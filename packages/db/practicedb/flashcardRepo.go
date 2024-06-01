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

type FlashcardsRepo struct {
	dbutils.SingleCollectionRepo[*FlashcardCollection]
}

func NewFlashcardsRepo(db *mongo.Database) *FlashcardsRepo {
	col := db.Collection(FlashcardsColName)
	return &FlashcardsRepo{
		SingleCollectionRepo: dbutils.SingleCollectionRepo[*FlashcardCollection]{Collection: col},
	}
}

func (r *FlashcardsRepo) InsertRaw(collection *FlashcardCollection) (*FlashcardCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	collection.SetID(primitive.NewObjectID())
	collection.SetInitTimeByNow()
	for idx, card := range collection.FlashCards {
		card.SetID(primitive.NewObjectID())
		card.SetInitTimeByNow()
		collection.FlashCards[idx] = card
	}
	_, err := r.InsertOne(ctx, collection)
	return collection, err
}

func (r *FlashcardsRepo) GetByUserID(
	userID primitive.ObjectID,
) ([]*FlashcardCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// find and sort by field "updatedAt"
	filter := bson.M{"userId": userID}
	cur, err := r.Find(ctx, filter, options.Find().SetSort(bson.M{"updatedAt": -1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	collections := make([]*FlashcardCollection, 0)
	if err := cur.All(ctx, &collections); err != nil {
		return nil, err
	}

	if collections == nil || len(collections) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return collections, nil
}

func (r *FlashcardsRepo) UpdateLastView(
	collectionID,
	flashcardID primitive.ObjectID,
) error {
	filter := bson.M{"_id": collectionID, "flashcards._id": flashcardID}
	update := bson.M{"$set": bson.M{
		"lastViewed": flashcardID,
		"updatedAt":  primitive.NewDateTimeFromTime(time.Now()),
	}}

	cur, err := r.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}

	if cur.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *FlashcardsRepo) GetCollectionsMetadataByID(collectionID primitive.ObjectID) (*CollectionMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pipeline := []bson.M{
		{"$match": bson.M{"_id": collectionID}},
		{"$project": bson.M{"flashcards": 0}},
		{"$sort": bson.M{"updatedAt": -1}},
	}

	cur, err := r.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	collections := make([]*CollectionMetadata, 0)

	if err := cur.All(ctx, &collections); err != nil {
		return nil, err
	}

	if collections == nil || len(collections) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return collections[0], nil
}

func (r *FlashcardsRepo) GetCollectionsMetadataByUserID(userID primitive.ObjectID) ([]*CollectionMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pipeline := []bson.M{
		{"$match": bson.M{"userId": userID}},
		{"$project": bson.M{"flashcards": 0}},
		{"$sort": bson.M{"updatedAt": -1}},
	}

	cur, err := r.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	collections := make([]*CollectionMetadata, 0)

	if err := cur.All(ctx, &collections); err != nil {
		return nil, err
	}

	return collections, nil
}

func (r *FlashcardsRepo) UpdateCollectionMetadata(
	collectionID primitive.ObjectID,
	metadata *CollectionMetadata,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": collectionID}
	update := bson.M{"$set": bson.M{
		"name":        metadata.Name,
		"description": metadata.Description,
		"updatedAt":   primitive.NewDateTimeFromTime(time.Now()),
	}}

	cur, err := r.UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}

	if cur.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *FlashcardsRepo) AddFlashcardToCollection(
	collectionID primitive.ObjectID,
	flashcard *Flashcard,
) (*Flashcard, error) {
	flashcard.SetID(primitive.NewObjectID())
	flashcard.SetInitTimeByNow()

	cur, err := r.Collection.UpdateByID(context.Background(), collectionID, bson.M{
		"$push": bson.M{"flashcards": flashcard},
	})
	if err != nil {
		return nil, err
	}

	if cur.MatchedCount == 0 || cur.ModifiedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return flashcard, nil
}

func (r *FlashcardsRepo) GetFlashcardByID(collectionID primitive.ObjectID, cardID primitive.ObjectID) (*Flashcard, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pipeline := []bson.M{
		{"$match": bson.M{"_id": collectionID}},
		{"$unwind": "$flashcards"},
		{"$replaceRoot": bson.M{"newRoot": "$flashcards"}},
		{"$match": bson.M{"_id": cardID}},
	}

	cur, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	flashcards := make([]*Flashcard, 0)
	if err := cur.All(ctx, &flashcards); err != nil {
		return nil, err
	}

	if len(flashcards) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return flashcards[0], nil
}

func (r *FlashcardsRepo) UpdateFlashCard(collectionID primitive.ObjectID, card Flashcard) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	card.SetUpdatedAtByNow()

	filter := bson.M{"_id": collectionID, "flashcards._id": card.ID}
	// explain the below update query?
	update := bson.M{"$set": bson.M{
		"flashcards.$.frontText": card.FrontText,
		"flashcards.$.backText":  card.BackText,
		"flashcards.$.updatedAt": card.UpdatedAt,
	}}

	cur, err := r.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if cur.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *FlashcardsRepo) DeleteFlashCard(collectionID primitive.ObjectID, cardId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	update := bson.M{"$pull": bson.M{"flashcards": bson.M{"_id": cardId}}}
	cur, err := r.UpdateByID(ctx, collectionID, update)
	if err != nil {
		return err
	}

	if cur.MatchedCount == 0 || cur.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
