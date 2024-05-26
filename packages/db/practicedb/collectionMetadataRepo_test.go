package practicedb_test

import (
	"blinders/packages/db/practicedb"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMetadataCollection(t *testing.T) {
	repo := GetTestCollectionMetadatasRepo(t)
	userID := primitive.NewObjectID()

	metadata := &practicedb.FlashCardCollectionMetadata{
		UserID:      userID,
		Name:        "Test Collection",
		Description: "This is a test collection",
		Viewed:      0,
		Total:       10,
	}
	insertedMetadata, err := repo.InsertRaw(metadata)
	assert.Nil(t, err)
	assert.Equal(t, metadata.UserID, insertedMetadata.UserID)
	assert.Equal(t, metadata.Name, insertedMetadata.Name)
	assert.Equal(t, metadata.Description, insertedMetadata.Description)
	assert.Equal(t, metadata.Viewed, insertedMetadata.Viewed)
	assert.Equal(t, metadata.Total, insertedMetadata.Total)

	foundMetadata, err := repo.GetByUserID(userID)
	assert.Nil(t, err)
	assert.Contains(t, foundMetadata, *insertedMetadata)
	defer CleanRepo(t, repo.Collection)
}

func GetTestCollectionMetadatasRepo(t *testing.T) *practicedb.CollectionMetadatasRepo {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoTestURL))
	assert.Nil(t, err)

	db := client.Database(mongoTestDB)

	return practicedb.NewCollectionMetadataRepo(db)
}

func CleanRepo(t *testing.T, repo *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	assert.Nil(t, repo.Drop(ctx))
}
