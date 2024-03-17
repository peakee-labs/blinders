package restapi

import (
	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/transport"

	"github.com/gofiber/fiber/v2"
)

type Manager struct {
	App           *fiber.App
	Auth          auth.Manager
	DB            *db.MongoManager
	Users         *UsersService
	Conversations *ConversationsService
	Messages      *MessagesService
}

func NewManager(
	app *fiber.App,
	auth auth.Manager,
	db *db.MongoManager,
	transporter transport.Transport,
	consumerMap transport.ConsumerMap,
) *Manager {
	return &Manager{
		App:           app,
		Auth:          auth,
		DB:            db,
		Users:         NewUsersService(db.Users, db.FriendRequests, transporter, consumerMap),
		Conversations: NewConversationsService(db.Conversations, db.Users),
		Messages:      NewMessagesService(db.Messages),
	}
}

type InitOptions struct {
	Prefix string
}

func (m Manager) InitRoute(options InitOptions) error {
	if options.Prefix == "" {
		options.Prefix = "/"
	}

	rootRoute := m.App.Group(options.Prefix)
	rootRoute.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello from Peakee Rest API")
	})

	authorizedWithoutUser := rootRoute.Group(
		"/users/self",
		auth.FiberAuthMiddleware(m.Auth, m.DB.Users,
			auth.MiddlewareOptions{
				CheckUser: false,
			}),
	)
	authorizedWithoutUser.Get("/", m.Users.GetSelfFromAuth)
	authorizedWithoutUser.Post("/", m.Users.CreateNewUserBySelf)

	authorized := rootRoute.Group("/", auth.FiberAuthMiddleware(m.Auth, m.DB.Users))

	users := authorized.Group("/users")
	users.Get("/:id", AllowPublicQuery, ValidateUserIDParam, m.Users.GetUserByID)
	users.Get("/:id/friend-requests", ValidateUserIDParam, m.Users.GetPendingFriendRequests)
	users.Post("/:id/friend-requests", ValidateUserIDParam, m.Users.CreateAddFriendRequest)
	users.Put("/:id/friend-requests/:requestId", ValidateUserIDParam, m.Users.RespondFriendRequest)

	conversations := authorized.Group("/conversations")
	conversations.Get("/:id", m.Conversations.GetConversationByID)
	conversations.Post("/", m.Conversations.CreateNewIndividualConversation)

	authorized.Get("/messages/:id", m.Messages.GetMessageByID)

	return nil
}
