package chatdb

import "go.mongodb.org/mongo-driver/mongo"

var (
	ConversationsCollection = "conversations"
	MessagesCollection      = "messages"
)

type ChatDB struct {
	mongo.Database
	ConversationsRepo *ConversationsRepo
	MessagesRepo      *MessagesRepo
}

func NewChatDB(db *mongo.Database) *ChatDB {
	return &ChatDB{
		Database:          *db,
		ConversationsRepo: NewConversationsRepo(db),
		MessagesRepo:      NewMessagesRepo(db),
	}
}
