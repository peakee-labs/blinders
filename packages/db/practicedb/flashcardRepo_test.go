package practicedb_test

import (
	"blinders/packages/db/practicedb"
	dbutils "blinders/packages/db/utils"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mongoTestURL                  = "mongodb://localhost:27017"
	mongoTestDBName               = "blinders-test"
	client          *mongo.Client = nil
)

func TestInsertFlashcardCollection(t *testing.T) {
	t.Parallel()
	r := GetFlashcardTestRepo(t)
	defer CleanRepo(t, r)

	collection := practicedb.FlashcardCollection{
		CollectionMetadata: practicedb.CollectionMetadata{
			UserID: primitive.NewObjectID(),
			Name:   "test collection",
			Type:   "DefaultFlashcard",
		},
		FlashCards: []*practicedb.Flashcard{},
	}

	insertedCollection, err := r.InsertRaw(&collection)
	assert.NoError(t, err)
	assert.NotNil(t, insertedCollection)

	assert.Equal(t, len(collection.FlashCards), len(insertedCollection.FlashCards))
	assert.NotNil(t, insertedCollection.ID)
	assert.False(t, insertedCollection.ID.IsZero())
	assert.Equal(t, collection.UserID, insertedCollection.UserID)
	assert.Equal(t, collection.Name, insertedCollection.Name)
	assert.Equal(t, collection.Type, insertedCollection.Type)

	gotCollection, err := r.GetByID(insertedCollection.ID)
	assert.Nil(t, err)
	assert.NotNil(t, gotCollection)

	assert.Equal(t, len(collection.FlashCards), len(gotCollection.FlashCards))
	assert.Equal(t, insertedCollection.UserID, gotCollection.UserID)
	assert.Equal(t, insertedCollection.Name, gotCollection.Name)
	assert.Equal(t, insertedCollection.Type, gotCollection.Type)
}

func TestGetByUserID(t *testing.T) {
	t.Parallel()
	r := GetFlashcardTestRepo(t)
	defer CleanRepo(t, r)

	collection := practicedb.FlashcardCollection{
		CollectionMetadata: practicedb.CollectionMetadata{
			UserID: primitive.NewObjectID(),
			Name:   "test collection",
			Type:   "DefaultFlashcard",
		},
		FlashCards: []*practicedb.Flashcard{},
	}

	insertedCollection, err := r.InsertRaw(&collection)
	assert.NoError(t, err)
	assert.NotNil(t, insertedCollection)

	assert.Equal(t, len(collection.FlashCards), len(insertedCollection.FlashCards))
	assert.NotNil(t, insertedCollection.ID)
	assert.False(t, insertedCollection.ID.IsZero())
	assert.Equal(t, collection.UserID, insertedCollection.UserID)
	assert.Equal(t, collection.Name, insertedCollection.Name)
	assert.Equal(t, collection.Type, insertedCollection.Type)

	collections, err := r.GetByUserID(insertedCollection.UserID)
	assert.Nil(t, err)
	assert.NotNil(t, collection)
	assert.Equal(t, 1, len(collections))
	gotCollection := collections[0]

	assert.Equal(t, len(collection.FlashCards), len(gotCollection.FlashCards))
	assert.Equal(t, insertedCollection.UserID, gotCollection.UserID)
	assert.Equal(t, insertedCollection.Name, gotCollection.Name)
	assert.Equal(t, insertedCollection.Type, gotCollection.Type)

	invalidUserID := primitive.NewObjectID()
	invalidCollections, err := r.GetByUserID(invalidUserID)
	assert.NotNil(t, err)
	assert.Nil(t, invalidCollections)
}

func TestGetCollectionMetadataByUserID(t *testing.T) {
	t.Parallel()
	r := GetFlashcardTestRepo(t)
	defer CleanRepo(t, r)

	userID := primitive.NewObjectID()

	collections := []*practicedb.FlashcardCollection{
		{
			CollectionMetadata: practicedb.CollectionMetadata{
				UserID: userID,
				Name:   "test collection1",
				Type:   "DefaultFlashcard",
			},
			FlashCards: []*practicedb.Flashcard{},
		},
		{
			CollectionMetadata: practicedb.CollectionMetadata{
				UserID: userID,
				Name:   "test collection2",
				Type:   "DefaultFlashcard",
			},
			FlashCards: []*practicedb.Flashcard{},
		},
	}

	for idx, collection := range collections {
		insertedCollection, err := r.InsertRaw(collection)
		assert.NoError(t, err)
		assert.NotNil(t, insertedCollection)
		collections[idx] = insertedCollection
	}

	metadatas, err := r.GetCollectionsMetadataByUserID(userID)
	assert.NoError(t, err)
	assert.NotNil(t, metadatas)

	for _, metadata := range metadatas {
		assert.Equal(t, userID, metadata.UserID)
		for _, collection := range collections {
			if collection.ID == metadata.ID {
				assert.Equal(t, collection.Name, metadata.Name)
				assert.Equal(t, collection.Type, metadata.Type)
				break
			}
		}
	}
}

func TestUpdateCollectionMetadata(t *testing.T) {
	t.Parallel()
	r := GetFlashcardTestRepo(t)
	defer CleanRepo(t, r)

	userID := primitive.NewObjectID()

	collections := []*practicedb.FlashcardCollection{
		{
			CollectionMetadata: practicedb.CollectionMetadata{
				UserID: userID,
				Name:   "test collection",
				Type:   "DefaultFlashcard",
			},
			FlashCards: []*practicedb.Flashcard{},
		},
		{
			CollectionMetadata: practicedb.CollectionMetadata{
				UserID: userID,
				Name:   "test collection",
				Type:   "DefaultFlashcard",
			},
			FlashCards: []*practicedb.Flashcard{},
		},
	}

	for idx, collection := range collections {
		insertedCollection, err := r.InsertRaw(collection)
		assert.NoError(t, err)
		assert.NotNil(t, insertedCollection)
		collections[idx] = insertedCollection
	}

	updateCollectionID := collections[0].ID
	update := collections[0].CollectionMetadata

	update.Name = "updated collection"
	err := r.UpdateCollectionMetadata(updateCollectionID, &update)
	assert.NoError(t, err)

	updatedCollection, err := r.GetCollectionsMetadataByID(updateCollectionID)
	assert.NoError(t, err)
	assert.Equal(t, update.Name, updatedCollection.Name)
	assert.Equal(t, update.Type, updatedCollection.Type)
	assert.Equal(t, update.UserID, updatedCollection.UserID)
	assert.LessOrEqual(t, update.UpdatedAt, updatedCollection.UpdatedAt)
}

func TestAddFlashcardToCollection(t *testing.T) {
	t.Parallel()
	r := GetFlashcardTestRepo(t)
	defer CleanRepo(t, r)

	collection := practicedb.FlashcardCollection{
		CollectionMetadata: practicedb.CollectionMetadata{
			UserID: primitive.NewObjectID(),
			Name:   "test collection",
			Type:   "DefaultFlashcard",
		},
		FlashCards: []*practicedb.Flashcard{},
	}

	insertedCollection, err := r.InsertRaw(&collection)
	assert.NoError(t, err)
	assert.NotNil(t, insertedCollection)

	flashcards := []*practicedb.Flashcard{
		{
			FrontText: "front text",
			BackText:  "back text",
		},
		{
			FrontText: "front text",
			BackText:  "back text",
		},
	}

	for _, flashcard := range flashcards {
		insertedFlashcard, err := r.AddFlashcardToCollection(insertedCollection.ID, flashcard)
		assert.Nil(t, err)
		assert.NotNil(t, insertedFlashcard)
	}

	updatedCollection, err := r.GetByID(insertedCollection.ID)
	assert.Nil(t, err)
	assert.NotNil(t, updatedCollection)

	for _, card := range updatedCollection.FlashCards {
		for _, flashcard := range flashcards {
			if card.ID == flashcard.ID {
				assert.Equal(t, flashcard.FrontText, card.FrontText)
				assert.Equal(t, flashcard.BackText, card.BackText)
				break
			}
		}
	}
}

func TestGetFlashCardByID(t *testing.T) {
	t.Parallel()
	r := GetFlashcardTestRepo(t)
	defer CleanRepo(t, r)

	collection := practicedb.FlashcardCollection{
		CollectionMetadata: practicedb.CollectionMetadata{
			UserID: primitive.NewObjectID(),
			Name:   "test collection",
			Type:   "DefaultFlashcard",
		},
		FlashCards: []*practicedb.Flashcard{},
	}

	insertedCollection, err := r.InsertRaw(&collection)
	assert.NoError(t, err)
	assert.NotNil(t, insertedCollection)

	flashcard := practicedb.Flashcard{
		FrontText: "front text",
		BackText:  "back text",
	}

	insertedFlashcard, err := r.AddFlashcardToCollection(insertedCollection.ID, &flashcard)
	assert.Nil(t, err)
	assert.NotNil(t, insertedFlashcard)

	gotFlashcard, err := r.GetFlashcardByID(insertedCollection.ID, flashcard.ID)
	assert.NoError(t, err)
	assert.NotNil(t, insertedCollection)

	assert.Equal(t, insertedFlashcard.FrontText, gotFlashcard.FrontText)
	assert.Equal(t, insertedFlashcard.BackText, gotFlashcard.BackText)
}

func TestUpdateFlashCardByID(t *testing.T) {
	t.Parallel()
	r := GetFlashcardTestRepo(t)
	defer CleanRepo(t, r)

	collection := practicedb.FlashcardCollection{
		CollectionMetadata: practicedb.CollectionMetadata{
			UserID: primitive.NewObjectID(),
			Name:   "test collection",
			Type:   "DefaultFlashcard",
		},
		FlashCards: []*practicedb.Flashcard{},
	}

	insertedCollection, err := r.InsertRaw(&collection)
	assert.NoError(t, err)
	assert.NotNil(t, insertedCollection)

	flashcard := &practicedb.Flashcard{
		FrontText: "front text",
		BackText:  "back text",
	}

	insertedFlashcard, err := r.AddFlashcardToCollection(insertedCollection.ID, flashcard)
	assert.Nil(t, err)
	assert.NotNil(t, insertedFlashcard)

	gotFlashcard, err := r.GetFlashcardByID(insertedCollection.ID, flashcard.ID)
	assert.NoError(t, err)
	assert.NotNil(t, insertedCollection)

	update := *gotFlashcard

	update.FrontText = "new front text"
	update.BackText = "new back text"

	err = r.UpdateFlashCard(insertedCollection.ID, update)
	assert.Nil(t, err)

	updatedFlashcard, err := r.GetFlashcardByID(insertedCollection.ID, update.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedFlashcard)

	assert.Equal(t, update.FrontText, updatedFlashcard.FrontText)
	assert.Equal(t, update.BackText, updatedFlashcard.BackText)
	assert.LessOrEqual(t, update.UpdatedAt.Time(), updatedFlashcard.UpdatedAt.Time())
}

func TestDeleteFlashcard(t *testing.T) {
	t.Parallel()
	r := GetFlashcardTestRepo(t)
	defer CleanRepo(t, r)

	collection := practicedb.FlashcardCollection{
		CollectionMetadata: practicedb.CollectionMetadata{
			UserID: primitive.NewObjectID(),
			Name:   "test collection",
			Type:   "DefaultFlashcard",
		},
		FlashCards: []*practicedb.Flashcard{},
	}

	insertedCollection, err := r.InsertRaw(&collection)
	assert.NoError(t, err)
	assert.NotNil(t, insertedCollection)

	flashcard := &practicedb.Flashcard{
		FrontText: "front text",
		BackText:  "back text",
	}

	insertedFlashcard, err := r.AddFlashcardToCollection(insertedCollection.ID, flashcard)
	assert.Nil(t, err)
	assert.NotNil(t, insertedFlashcard)

	updatedCollection, err := r.GetByID(insertedCollection.ID)
	assert.Nil(t, err)
	assert.NotNil(t, insertedFlashcard)
	assert.Equal(t, len(insertedCollection.FlashCards)+1, len(updatedCollection.FlashCards))

	err = r.DeleteFlashCard(insertedCollection.ID, insertedFlashcard.ID)
	assert.Nil(t, err)

	failed, err := r.GetFlashcardByID(insertedCollection.ID, insertedFlashcard.ID)
	assert.NotNil(t, err)
	assert.Nil(t, failed)
}

func TestUpdateLastView(t *testing.T) {
	t.Parallel()
	r := GetFlashcardTestRepo(t)
	defer CleanRepo(t, r)

	collection := practicedb.FlashcardCollection{
		CollectionMetadata: practicedb.CollectionMetadata{
			UserID: primitive.NewObjectID(),
			Name:   "test collection",
			Type:   "DefaultFlashcard",
		},
		FlashCards: []*practicedb.Flashcard{},
	}

	insertedCollection, err := r.InsertRaw(&collection)
	assert.NoError(t, err)
	assert.NotNil(t, insertedCollection)

	flashcards := []*practicedb.Flashcard{
		{
			FrontText: "front text",
			BackText:  "back text",
		},
		{
			FrontText: "another front text",
			BackText:  "another back text",
		},
	}

	for idx, flashcard := range flashcards {
		insertedFlashcard, err := r.AddFlashcardToCollection(insertedCollection.ID, flashcard)
		assert.Nil(t, err)
		assert.NotNil(t, insertedFlashcard)
		flashcards[idx] = insertedFlashcard
	}

	for _, flashcard := range flashcards {
		col, err := r.GetByID(insertedCollection.ID)
		assert.NoError(t, err)
		assert.NotEqual(t, col.LastViewed, flashcard.ID)

		err = r.UpdateLastView(insertedCollection.ID, flashcard.ID)
		assert.Nil(t, err)

		col, err = r.GetByID(insertedCollection.ID)
		assert.NoError(t, err)
		assert.Equal(t, col.LastViewed, flashcard.ID)
	}
}

func GetFlashcardTestRepo(t *testing.T) *practicedb.FlashcardsRepo {
	t.Helper()
	if client == nil {
		var err error
		client, err = dbutils.InitMongoClient(mongoTestURL)
		assert.NoError(t, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := client.Ping(ctx, nil)
	assert.NoError(t, err)

	return practicedb.NewFlashcardsRepo(client.Database(mongoTestDBName))
}

func CleanRepo(t *testing.T, repo *practicedb.FlashcardsRepo) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := repo.Collection.Drop(ctx)
	assert.NoError(t, err)
}
