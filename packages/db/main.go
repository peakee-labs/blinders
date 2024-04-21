package db

import (
	"log"

	"blinders/packages/db/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

// username:password@host:port/database
const MongoURLTemplate = "mongodb://%s:%s@%s:%s/%s"

const (
	UserCollection          = "users"
	ConversationCollection  = "conversations"
	MessageCollection       = "messages"
	MatchCollection         = "matches"
	FriendRequestCollection = "friendrequests"
	FeedbackCollection      = "feedbacks"
)

type MongoManager struct {
	Client         *mongo.Client
	Database       string
	Users          *repo.UsersRepo
	Conversations  *repo.ConversationsRepo
	Messages       *repo.MessagesRepo
	Matches        *repo.MatchesRepo
	FriendRequests *repo.FriendRequestsRepo
	Feedbacks      *repo.FeedbacksRepo
}

func NewMongoManager(url string, name string) *MongoManager {
	client, err := InitMongoClient(url)
	if err != nil {
		log.Println("cannot init mongo client", err)
		return nil
	}

	return &MongoManager{
		Client:   client,
		Database: name,
		Users:    repo.NewUsersRepo(client.Database(name).Collection(UserCollection)),
		Conversations: repo.NewConversationsRepo(
			client.Database(name).Collection(ConversationCollection),
		),
		Messages: repo.NewMessagesRepo(client.Database(name).Collection(MessageCollection)),
		Matches:  repo.NewMatchesRepo(client.Database(name).Collection(MatchCollection)),
		FriendRequests: repo.NewFriendRequestsRepo(
			client.Database(name).Collection(FriendRequestCollection),
		),
		Feedbacks: repo.NewFeedbackRepo(client.Database(name).Collection(FeedbackCollection)),
	}
}
