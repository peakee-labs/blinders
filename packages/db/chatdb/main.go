package chatdb

import "go.mongodb.org/mongo-driver/mongo"

var (
	ConversationsCollection = "conversations"
	MessagesCollection      = "messages"
)

type UsersDB struct {
	mongo.Database
	ConversationsRepo *ConversationsRepo
	MessagesRepo      *MessagesRepo
}

func NewUsersDB(db *mongo.Database) *UsersDB {
	return &UsersDB{
		Database:          *db,
		ConversationsRepo: NewConversationsRepo(db),
		MessagesRepo:      NewMessagesRepo(db),
	}
}
