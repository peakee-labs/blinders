package repo_test

import (
	"context"
	"testing"
	"time"

	dbutils "blinders/packages/dbutils"
	"blinders/services/practice/repo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInsertSnapshot(t *testing.T) {
	snapshotRepo := GetSnapshotTestRepo(t)
	defer CleanSnapshotRepo(t, snapshotRepo)

	var snapshotTestType repo.SnapshotType = "test-type"

	snapshot := &repo.PracticeSnapshot{
		UserID:  primitive.NewObjectID(),
		Type:    snapshotTestType,
		Current: primitive.NewDateTimeFromTime(time.Now()),
	}

	insertedSnapshot, err := snapshotRepo.InsertRaw(snapshot)
	assert.NoError(t, err)

	assert.NotNil(t, insertedSnapshot)
	assert.Equal(t, snapshot.UserID, insertedSnapshot.UserID)
	assert.Equal(t, snapshot.Type, insertedSnapshot.Type)
	assert.Equal(t, snapshot.Current, insertedSnapshot.Current)
}

func TestGetSnapshot(t *testing.T) {
	snapshotRepo := GetSnapshotTestRepo(t)
	defer CleanSnapshotRepo(t, snapshotRepo)

	var snapshotTestType repo.SnapshotType = "test-type"

	snapshot := &repo.PracticeSnapshot{
		UserID:  primitive.NewObjectID(),
		Type:    snapshotTestType,
		Current: primitive.NewDateTimeFromTime(time.Now()),
	}

	insertedSnapshot, err := snapshotRepo.InsertRaw(snapshot)
	assert.NoError(t, err)
	assert.NotNil(t, insertedSnapshot)

	foundSnapshot, err := snapshotRepo.GetSnapshotOfUserByType(
		insertedSnapshot.UserID,
		insertedSnapshot.Type,
	)
	assert.NoError(t, err)
	assert.NotNil(t, foundSnapshot)

	assert.Equal(t, insertedSnapshot.RawModel, foundSnapshot.RawModel)
	assert.Equal(t, insertedSnapshot.Type, foundSnapshot.Type)
	assert.Equal(t, insertedSnapshot.UserID, foundSnapshot.UserID)
	assert.Equal(t, insertedSnapshot.Current, foundSnapshot.Current)

	invalidUserID := primitive.NewObjectID()
	notFoundSnapshot, err := snapshotRepo.GetSnapshotOfUserByType(
		invalidUserID,
		insertedSnapshot.Type,
	)
	assert.Error(t, err)
	assert.Nil(t, notFoundSnapshot)

	var invalidType repo.SnapshotType = "invalid-type"
	notFoundSnapshot, err = snapshotRepo.GetSnapshotOfUserByType(
		insertedSnapshot.UserID,
		invalidType,
	)
	assert.Error(t, err)
	assert.Nil(t, notFoundSnapshot)
}

func TestUpdateSnapshot(t *testing.T) {
	snapshotRepo := GetSnapshotTestRepo(t)
	defer CleanSnapshotRepo(t, snapshotRepo)

	var snapshotTestType repo.SnapshotType = "test-type"

	snapshot := &repo.PracticeSnapshot{
		UserID:  primitive.NewObjectID(),
		Type:    snapshotTestType,
		Current: primitive.NewDateTimeFromTime(time.Now()),
	}

	insertedSnapshot, err := snapshotRepo.InsertRaw(snapshot)
	assert.NoError(t, err)
	assert.NotNil(t, insertedSnapshot)

	update := *insertedSnapshot
	update.Current = primitive.NewDateTimeFromTime(time.Now().Add(time.Hour))

	updatedSnapshot, err := snapshotRepo.UpdateSnapshot(&update)
	assert.NoError(t, err)
	assert.NotNil(t, updatedSnapshot)

	assert.LessOrEqual(t, insertedSnapshot.UpdatedAt, updatedSnapshot.UpdatedAt)
	assert.Equal(t, update.Current, updatedSnapshot.Current)

	notExistedSnapshot := repo.PracticeSnapshot{
		UserID:  primitive.NewObjectID(),
		Type:    snapshotTestType,
		Current: primitive.NewDateTimeFromTime(time.Now()),
	}
	notExistedSnapshot.SetID(primitive.NewObjectID())

	invalidUpdate, err := snapshotRepo.UpdateSnapshot(&notExistedSnapshot)
	assert.Error(t, err)
	assert.Nil(t, invalidUpdate)
}

func GetSnapshotTestRepo(t *testing.T) *repo.SnapshotsRepo {
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

	return repo.NewSnapshotsRepo(client.Database(mongoTestDBName))
}

func CleanSnapshotRepo(t *testing.T, repo *repo.SnapshotsRepo) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := repo.Collection.Drop(ctx)
	assert.NoError(t, err)
}
