package practicedb_test

import (
	"blinders/packages/db/practicedb"
	dbutils "blinders/packages/db/utils"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInsertSnapshot(t *testing.T) {
	repo := GetSnapshotTestRepo(t)
	defer CleanSnapshotRepo(t, repo)

	var snapshotTestType practicedb.SnapshotType = "test-type"

	snapshot := &practicedb.PracticeSnapshot{
		UserID:  primitive.NewObjectID(),
		Type:    snapshotTestType,
		Current: primitive.NewDateTimeFromTime(time.Now()),
	}

	insertedSnapshot, err := repo.InsertRaw(snapshot)
	assert.NoError(t, err)

	assert.NotNil(t, insertedSnapshot)
	assert.Equal(t, snapshot.UserID, insertedSnapshot.UserID)
	assert.Equal(t, snapshot.Type, insertedSnapshot.Type)
	assert.Equal(t, snapshot.Current, insertedSnapshot.Current)
}

func TestGetSnapshot(t *testing.T) {
	repo := GetSnapshotTestRepo(t)
	defer CleanSnapshotRepo(t, repo)

	var snapshotTestType practicedb.SnapshotType = "test-type"

	snapshot := &practicedb.PracticeSnapshot{
		UserID:  primitive.NewObjectID(),
		Type:    snapshotTestType,
		Current: primitive.NewDateTimeFromTime(time.Now()),
	}

	insertedSnapshot, err := repo.InsertRaw(snapshot)
	assert.NoError(t, err)
	assert.NotNil(t, insertedSnapshot)

	foundSnapshot, err := repo.GetSnapshotOfUserByType(insertedSnapshot.UserID, insertedSnapshot.Type)
	assert.NoError(t, err)
	assert.NotNil(t, foundSnapshot)

	assert.Equal(t, insertedSnapshot.RawModel, foundSnapshot.RawModel)
	assert.Equal(t, insertedSnapshot.Type, foundSnapshot.Type)
	assert.Equal(t, insertedSnapshot.UserID, foundSnapshot.UserID)
	assert.Equal(t, insertedSnapshot.Current, foundSnapshot.Current)

	invalidUserID := primitive.NewObjectID()
	notFoundSnapshot, err := repo.GetSnapshotOfUserByType(invalidUserID, insertedSnapshot.Type)
	assert.Error(t, err)
	assert.Nil(t, notFoundSnapshot)

	var invalidType practicedb.SnapshotType = "invalid-type"
	notFoundSnapshot, err = repo.GetSnapshotOfUserByType(insertedSnapshot.UserID, invalidType)
	assert.Error(t, err)
	assert.Nil(t, notFoundSnapshot)
}

func TestUpdateSnapshot(t *testing.T) {
	repo := GetSnapshotTestRepo(t)
	defer CleanSnapshotRepo(t, repo)

	var snapshotTestType practicedb.SnapshotType = "test-type"

	snapshot := &practicedb.PracticeSnapshot{
		UserID:  primitive.NewObjectID(),
		Type:    snapshotTestType,
		Current: primitive.NewDateTimeFromTime(time.Now()),
	}

	insertedSnapshot, err := repo.InsertRaw(snapshot)
	assert.NoError(t, err)
	assert.NotNil(t, insertedSnapshot)

	update := *insertedSnapshot
	update.Current = primitive.NewDateTimeFromTime(time.Now().Add(time.Hour))

	updatedSnapshot, err := repo.UpdateSnapshot(&update)
	assert.NoError(t, err)
	assert.NotNil(t, updatedSnapshot)

	assert.LessOrEqual(t, insertedSnapshot.UpdatedAt, updatedSnapshot.UpdatedAt)
	assert.Equal(t, update.Current, updatedSnapshot.Current)

	notExistedSnapshot := practicedb.PracticeSnapshot{
		UserID:  primitive.NewObjectID(),
		Type:    snapshotTestType,
		Current: primitive.NewDateTimeFromTime(time.Now()),
	}
	notExistedSnapshot.SetID(primitive.NewObjectID())

	invalidUpdate, err := repo.UpdateSnapshot(&notExistedSnapshot)
	assert.Error(t, err)
	assert.Nil(t, invalidUpdate)
}

func GetSnapshotTestRepo(t *testing.T) *practicedb.SnapshotsRepo {
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

	return practicedb.NewSnapshotsRepo(client.Database(mongoTestDBName))
}

func CleanSnapshotRepo(t *testing.T, repo *practicedb.SnapshotsRepo) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := repo.Collection.Drop(ctx)
	assert.NoError(t, err)
}
