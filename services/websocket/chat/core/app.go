package wschat

import (
	"blinders/packages/session"

	chatrepo "blinders/services/chat/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

var app *App

type App struct {
	Session      *session.Manager
	MessagesRepo *chatrepo.MessagesRepo
	ConvsRepo    *chatrepo.ConversationsRepo
}

// init app construct an app instance for internal use
// is that violate stateless of functional design? app instance is used in a func
func InitChatApp(sm *session.Manager, mongoDB *mongo.Database) *App {
	app = &App{
		Session:      sm,
		MessagesRepo: chatrepo.NewMessagesRepo(mongoDB),
		ConvsRepo:    chatrepo.NewConversationsRepo(mongoDB),
	}

	return app
}
