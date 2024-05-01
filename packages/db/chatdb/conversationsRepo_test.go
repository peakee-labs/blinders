package chatdb_test

import (
	"testing"

	"blinders/packages/db/chatdb"
	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	cclient, _ = dbutils.InitMongoClient("mongodb://localhost:27017")
	convRepo   = chatdb.NewConversationsRepo(cclient.Database("blinders"))
	usersRepo  = usersdb.NewUsersRepo(cclient.Database("blinders"))
)

func TestInsertIndividualConversationSuccess(t *testing.T) {
	user, _ := usersRepo.InsertNewRawUser(
		usersdb.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	friend, _ := usersRepo.InsertNewRawUser(
		usersdb.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	conv, err := convRepo.InsertIndividualConversation(user.ID, friend.ID)
	assert.Nil(t, err)
	assert.Equal(t, len(conv.Members), 2)
	assert.Equal(t, conv.CreatedBy, user.ID)
	assert.Equal(t, conv.Type, chatdb.IndividualConversation)
	assert.Equal(t, conv.Members[0].UserID, user.ID)
	assert.Equal(t, conv.Members[1].UserID, friend.ID)
}

func TestInsertIndividualConversationFailedWithDuplicatedConversation(t *testing.T) {
	user, _ := usersRepo.InsertNewRawUser(
		usersdb.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	friend, _ := usersRepo.InsertNewRawUser(
		usersdb.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	_, err := convRepo.InsertIndividualConversation(user.ID, friend.ID)
	assert.Nil(t, err)
	conv, err := convRepo.InsertIndividualConversation(user.ID, friend.ID)
	assert.NotNil(t, err)
	assert.Nil(t, conv)
}

func TestGetConversationWithAFriend(t *testing.T) {
	user, _ := usersRepo.InsertNewRawUser(
		usersdb.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	friend, _ := usersRepo.InsertNewRawUser(
		usersdb.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	conv, _ := convRepo.InsertIndividualConversation(user.ID, friend.ID)
	conversations, err := convRepo.GetConversationByMembers(
		[]primitive.ObjectID{user.ID, friend.ID}, chatdb.IndividualConversation,
	)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(*conversations))
	if len(*conversations) > 0 {
		assert.Equal(t, conv.ID, (*conversations)[0].ID)
	}
}
