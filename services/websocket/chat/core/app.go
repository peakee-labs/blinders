package wschat

import (
	"blinders/packages/db/chatdb"
	"blinders/packages/session"
)

var app *App

type App struct {
	Session *session.Manager
	ChatDB  *chatdb.ChatDB
}

// init app construct an app instance for internal use
// is that violate stateless of functional design? app instance is used in a func
func InitChatApp(sm *session.Manager, db *chatdb.ChatDB) *App {
	app = &App{
		Session: sm,
		ChatDB:  db,
	}

	return app
}
