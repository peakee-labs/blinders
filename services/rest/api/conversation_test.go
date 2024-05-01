package restapi_test

import (
	"log"
	"testing"

	"blinders/packages/db/chatdb"
	"blinders/packages/db/usersdb"
	dbutils "blinders/packages/db/utils"
	restapi "blinders/services/rest/api"

	"github.com/test-go/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var convService restapi.ConversationsService

func init() {
	client, _ := dbutils.InitMongoClient("mongodb://localhost:27017")
	chatDB := chatdb.NewChatDB(client.Database("blinders"))
	usersDB := usersdb.NewUsersDB(client.Database("blinders"))
	convService = *restapi.NewConversationsService(chatDB.ConversationsRepo, chatDB.MessagesRepo, usersDB.UsersRepo)
}

func TestCheckFriendshipFailedWithNoFriendship(t *testing.T) {
	user1, _ := convService.UsersRepo.InsertNewRawUser(usersdb.User{
		FirebaseUID: primitive.NewObjectID().Hex(),
	})

	err := convService.CheckFriendRelationship(user1.ID, user1.ID)
	assert.NotNil(t, err)
}

func TestCheckFriendshipFailedWithFriendNotFound(t *testing.T) {
	friendID := primitive.NewObjectID()
	user1, _ := convService.UsersRepo.InsertNewRawUser(usersdb.User{
		FriendIDs:   []primitive.ObjectID{friendID},
		FirebaseUID: primitive.NewObjectID().Hex(),
	})
	err := convService.CheckFriendRelationship(user1.ID, friendID)
	assert.NotNil(t, err)
}

func TestCheckFriendshipSuccess(t *testing.T) {
	user1ID := primitive.NewObjectID()
	user2ID := primitive.NewObjectID()
	log.Println(user1ID, user2ID)
	user1, _ := convService.UsersRepo.InsertNewUser(usersdb.User{
		ID:          user1ID,
		FriendIDs:   []primitive.ObjectID{user2ID},
		FirebaseUID: primitive.NewObjectID().Hex(),
	})
	user2, _ := convService.UsersRepo.InsertNewUser(usersdb.User{
		ID:          user2ID,
		FriendIDs:   []primitive.ObjectID{user1ID},
		FirebaseUID: primitive.NewObjectID().Hex(),
	})

	err := convService.CheckFriendRelationship(user1.ID, user2.ID)
	assert.Nil(t, err)

	err = convService.CheckFriendRelationship(user2.ID, user1.ID)
	assert.Nil(t, err)
}
