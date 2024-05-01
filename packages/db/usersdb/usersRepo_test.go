package usersdb_test

import (
	"testing"

	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	uclient, _ = dbutils.InitMongoClient("mongodb://localhost:27017")
	userRepo   = usersdb.NewUsersRepo(uclient.Database("blinders"))
)

func TestInsertUser(t *testing.T) {
	user := usersdb.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	newUser, err := userRepo.InsertNewRawUser(user)
	assert.Nil(t, err)
	assert.NotEqual(t, newUser.ID, primitive.ObjectID{})
	assert.Equal(t, user.ID, primitive.ObjectID{})
}

func TestInsertUserFailedWithDuplicatedFirebaseUID(t *testing.T) {
	user := usersdb.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	_, _ = userRepo.InsertNewRawUser(user)
	_, err := userRepo.InsertNewRawUser(user)
	assert.NotNil(t, err)
}

func TestGetUserByFirebaseUID(t *testing.T) {
	user := usersdb.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	user, _ = userRepo.InsertNewRawUser(user)
	queriedUser, err := userRepo.GetUserByFirebaseUID(user.FirebaseUID)
	assert.Nil(t, err)
	assert.Equal(t, user, queriedUser)
}

func TestGetUserByID(t *testing.T) {
	user := usersdb.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	user, _ = userRepo.InsertNewRawUser(user)
	queriedUser, err := userRepo.GetUserByID(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, user, queriedUser)
}

func TestGetUserByIDNotFound(t *testing.T) {
	_, err := userRepo.GetUserByID(primitive.NewObjectID())
	assert.NotNil(t, err)
}

func TestGetUserByFirebaseUIDNotFound(t *testing.T) {
	_, err := userRepo.GetUserByFirebaseUID(primitive.NewObjectID().String())
	assert.NotNil(t, err)
}

func TestDeleteUserByUserID(t *testing.T) {
	user := usersdb.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	user, _ = userRepo.InsertNewRawUser(user)

	queriedUser, err := userRepo.GetUserByID(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, user, queriedUser)

	deleted, err := userRepo.DeleteUserByID(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, user, deleted)

	failedDelete, err := userRepo.DeleteUserByID(user.ID)
	assert.NotNil(t, err)
	assert.Equal(t, usersdb.User{}, failedDelete)
}

func TestAddFriend(t *testing.T) {
	user1, _ := userRepo.InsertNewRawUser(usersdb.User{
		FirebaseUID: primitive.NewObjectID().Hex(),
		FriendIDs:   make([]primitive.ObjectID, 0),
	})
	user2, _ := userRepo.InsertNewRawUser(usersdb.User{
		FirebaseUID: primitive.NewObjectID().Hex(),
		FriendIDs:   make([]primitive.ObjectID, 0),
	})

	err := userRepo.AddFriend(user1.ID, user2.ID)
	assert.Nil(t, err)
	err = userRepo.AddFriend(user1.ID, user2.ID)
	assert.NotNil(t, err)
}
