package repo_test

import (
	"testing"

	dbutils "blinders/packages/dbutils"
	"blinders/services/chat/repo"
	usersrepo "blinders/services/users/repo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	mongoClient, _ = dbutils.InitMongoClient("mongodb://localhost:27017")
	convRepo       = repo.NewConversationsRepo(mongoClient.Database("blinders"))
	usersRepo      = usersrepo.NewUsersRepo(mongoClient.Database("blinders"))
)

func TestInsertIndividualConversationSuccess(t *testing.T) {
	user, _ := usersRepo.InsertNewRawUser(
		usersrepo.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	friend, _ := usersRepo.InsertNewRawUser(
		usersrepo.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	conv, err := convRepo.InsertIndividualConversation(user.ID, friend.ID)
	assert.Nil(t, err)
	assert.Equal(t, len(conv.Members), 2)
	assert.Equal(t, conv.CreatedBy, user.ID)
	assert.Equal(t, conv.Type, repo.IndividualConversation)
	assert.Equal(t, conv.Members[0].UserID, user.ID)
	assert.Equal(t, conv.Members[1].UserID, friend.ID)
}

func TestInsertIndividualConversationFailedWithDuplicatedConversation(t *testing.T) {
	user, _ := usersRepo.InsertNewRawUser(
		usersrepo.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	friend, _ := usersRepo.InsertNewRawUser(
		usersrepo.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	_, err := convRepo.InsertIndividualConversation(user.ID, friend.ID)
	assert.Nil(t, err)
	conv, err := convRepo.InsertIndividualConversation(user.ID, friend.ID)
	assert.NotNil(t, err)
	assert.Nil(t, conv)
}

func TestGetConversationWithAFriend(t *testing.T) {
	user, _ := usersRepo.InsertNewRawUser(
		usersrepo.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	friend, _ := usersRepo.InsertNewRawUser(
		usersrepo.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	conv, _ := convRepo.InsertIndividualConversation(user.ID, friend.ID)
	conversations, err := convRepo.GetConversationByMembers(
		[]primitive.ObjectID{user.ID, friend.ID}, repo.IndividualConversation,
	)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(*conversations))
	if len(*conversations) > 0 {
		assert.Equal(t, conv.ID, (*conversations)[0].ID)
	}
}
