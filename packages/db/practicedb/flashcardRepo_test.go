package practicedb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoTestURL = "mongodb://localhost:27017"
	mongoTestDB  = "blinder-test"
)

func TestFlashCardsRepo(t *testing.T) {
	repo := GetTestRepo(t)

	var (
		userID       = primitive.NewObjectID()
		collectionID = primitive.NewObjectID()
	)

	card := &FlashCard{
		ID:           primitive.NewObjectID(),
		UserID:       userID,
		CollectionID: collectionID,
		FrontText:    "Front of the card",
		BackText:     "Back of the card",
	}

	insertedCard, err := repo.InsertFlashCard(card)
	require.Nil(t, err)

	foundCard, err := repo.GetFlashCardByID(insertedCard.ID)
	require.Nil(t, err)
	require.Equal(t, *insertedCard, *foundCard)

	foundWithUserID, err := repo.GetFlashCardByUserID(userID)
	require.Nil(t, err)
	require.Contains(t, foundWithUserID, *insertedCard)

	foundWithCollectionID, err := repo.GetFlashCardsByCollectionID(collectionID)
	require.Nil(t, err)
	require.Contains(t, foundWithCollectionID, *insertedCard)

	collection, err := repo.GetFlashCardCollectionByID(collectionID)
	require.Nil(t, err)
	require.Contains(t, collection.FlashCards, *insertedCard)
	require.Equal(t, collection.ID, collectionID)

	// Test UpdateFlashCard
	newCard := &FlashCard{
		ID:           insertedCard.ID,
		UserID:       insertedCard.UserID,
		CollectionID: insertedCard.CollectionID,
		FrontText:    "Updated front of the card",
		BackText:     "Updated back of the card",
	}
	err = repo.UpdateFlashCard(insertedCard.ID, newCard)
	require.Nil(t, err)

	updatedCard, err := repo.GetFlashCardByID(insertedCard.ID)
	require.Nil(t, err)
	require.Equal(t, *newCard, *updatedCard)

	err = repo.DeleteFlashCardByID(insertedCard.ID)
	require.Nil(t, err)
}

func TestGetFlashCardCollectionsByUserID(t *testing.T) {
	repo := GetTestRepo(t) // Assuming you have a function to create a new repo
	userID := primitive.NewObjectID()

	collectionsID := []primitive.ObjectID{
		primitive.NewObjectID(),
		primitive.NewObjectID(),
		primitive.NewObjectID(),
	}

	cards := []FlashCard{}

	for i := 0; i < 100; i++ {
		card := &FlashCard{
			UserID:       userID,
			CollectionID: collectionsID[i%len(collectionsID)], // This will distribute the flashcards across the collections
			FrontText:    fmt.Sprintf("sample front text %d", i),
			BackText:     fmt.Sprintf("sample back text %d", i),
		}
		insertedCard, err := repo.InsertRawFlashCard(card)
		assert.NoError(t, err)
		assert.NotNil(t, insertedCard)
		assert.NotNil(t, insertedCard.ID)
		cards = append(cards, *insertedCard)
	}

	result, err := repo.GetFlashCardCollectionsByUserID(userID)
	assert.NoError(t, err)

	// Verify that the returned collections are the ones we used
	assert.Equal(t, len(collectionsID), len(result))
	for _, collection := range result {
		assert.Contains(t, collectionsID, collection.ID)
		assert.Equal(t, userID, collection.UserID)
		for _, card := range collection.FlashCards {
			assert.Contains(t, cards, card)
		}
	}

}

func GetTestRepo(t *testing.T) *FlashCardsRepo {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoTestURL))
	require.Nil(t, err)

	db := client.Database(mongoTestDB)

	return NewFlashCardRepo(db)
}
