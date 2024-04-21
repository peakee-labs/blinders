package restapi_test

import (
	"log"
	"os"
	"testing"

	"blinders/packages/db"
	"blinders/packages/db/models"
	restapi "blinders/services/rest/api"

	"github.com/joho/godotenv"
	"github.com/test-go/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var convService restapi.ConversationsService

func init() {
	_ = godotenv.Load("../../../.env")
	dbName := os.Getenv("MONGO_DATABASE")
	url := os.Getenv("MONGO_DATABASE_URL")
	log.Println("database url: ", url)
	dbManager := db.NewMongoManager(url, dbName)
	convService = *restapi.NewConversationsService(dbManager.Conversations, dbManager.Users, nil)
}

func TestCheckFriendshipFailedWithNoFriendship(t *testing.T) {
	user1, _ := convService.UserRepo.InsertNewRawUser(models.User{
		FirebaseUID: primitive.NewObjectID().Hex(),
	})

	err := convService.CheckFriendRelationship(user1.ID, user1.ID)
	assert.NotNil(t, err)
}

func TestCheckFriendshipFailedWithFriendNotFound(t *testing.T) {
	friendID := primitive.NewObjectID()
	user1, _ := convService.UserRepo.InsertNewRawUser(models.User{
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
	user1, _ := convService.UserRepo.InsertNewUser(models.User{
		ID:          user1ID,
		FriendIDs:   []primitive.ObjectID{user2ID},
		FirebaseUID: primitive.NewObjectID().Hex(),
	})
	user2, _ := convService.UserRepo.InsertNewUser(models.User{
		ID:          user2ID,
		FriendIDs:   []primitive.ObjectID{user1ID},
		FirebaseUID: primitive.NewObjectID().Hex(),
	})

	err := convService.CheckFriendRelationship(user1.ID, user2.ID)
	assert.Nil(t, err)

	err = convService.CheckFriendRelationship(user2.ID, user1.ID)
	assert.Nil(t, err)
}
