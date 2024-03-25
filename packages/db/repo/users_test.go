package repo_test

import (
	"testing"

	"blinders/packages/db"
	"blinders/packages/db/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var manager = db.NewMongoManager("mongodb://localhost:27017", "blinders")

func TestInsertUser(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	newUser, err := manager.Users.InsertNewRawUser(user)
	assert.Nil(t, err)
	assert.NotEqual(t, newUser.ID, primitive.ObjectID{})
	assert.Equal(t, user.ID, primitive.ObjectID{})
}

func TestInsertUserFailedWithDuplicatedFirebaseUID(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	_, _ = manager.Users.InsertNewRawUser(user)
	_, err := manager.Users.InsertNewRawUser(user)
	assert.NotNil(t, err)
}

func TestGetUserByFirebaseUID(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	user, _ = manager.Users.InsertNewRawUser(user)
	queriedUser, err := manager.Users.GetUserByFirebaseUID(user.FirebaseUID)
	assert.Nil(t, err)
	assert.Equal(t, user, queriedUser)
}

func TestGetUserByID(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	user, _ = manager.Users.InsertNewRawUser(user)
	queriedUser, err := manager.Users.GetUserByID(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, user, queriedUser)
}

func TestGetUserByIDNotFound(t *testing.T) {
	_, err := manager.Users.GetUserByID(primitive.NewObjectID())
	assert.NotNil(t, err)
}

func TestGetUserByFirebaseUIDNotFound(t *testing.T) {
	_, err := manager.Users.GetUserByFirebaseUID(primitive.NewObjectID().String())
	assert.NotNil(t, err)
}

func TestDeleteUserByUserID(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	user, _ = manager.Users.InsertNewRawUser(user)

	queriedUser, err := manager.Users.GetUserByID(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, user, queriedUser)

	deleted, err := manager.Users.DeleteUserByID(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, user, deleted)

	failedDelete, err := manager.Users.DeleteUserByID(user.ID)
	assert.NotNil(t, err)
	assert.Equal(t, models.User{}, failedDelete)
}

func TestAddFriend(t *testing.T) {
	user1, _ := manager.Users.InsertNewRawUser(models.User{
		FirebaseUID: primitive.NewObjectID().Hex(),
		FriendIDs:   make([]primitive.ObjectID, 0),
	})
	user2, _ := manager.Users.InsertNewRawUser(models.User{
		FirebaseUID: primitive.NewObjectID().Hex(),
		FriendIDs:   make([]primitive.ObjectID, 0),
	})

	err := manager.Users.AddFriend(user1.ID, user2.ID)
	assert.Nil(t, err)
	err = manager.Users.AddFriend(user1.ID, user2.ID)
	assert.NotNil(t, err)
}
