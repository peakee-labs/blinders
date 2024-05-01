package restapi

import (
	"blinders/packages/auth"
	"blinders/packages/db/chatdb"
	"blinders/packages/db/matchingdb"
	"blinders/packages/db/usersdb"
	"blinders/packages/transport"

	"github.com/gofiber/fiber/v2"
)

type Manager struct {
	App           *fiber.App
	Auth          auth.Manager
	UsersRepo     *usersdb.UsersRepo
	Users         *UsersService
	Conversations *ConversationsService
	Messages      *MessagesService
	Onboardings   *OnboardingService
	Feedbacks     *FeedbacksService
}

func NewManager(
	app *fiber.App,
	auth auth.Manager,
	usersDB *usersdb.UsersDB,
	chatDB *chatdb.ChatDB,
	matchingRepo *matchingdb.MatchingRepo,
	transporter transport.Transport,
	consumerMap transport.ConsumerMap,
) *Manager {
	return &Manager{
		App:       app,
		Auth:      auth,
		UsersRepo: usersDB.UsersRepo,
		Users: NewUsersService(
			usersDB.UsersRepo,
			usersDB.FriendRequestsRepo,
			transporter,
			consumerMap,
		),
		Conversations: NewConversationsService(
			chatDB.ConversationsRepo,
			chatDB.MessagesRepo,
			usersDB.UsersRepo,
		),
		Messages: NewMessagesService(chatDB.MessagesRepo),
		Onboardings: NewOnboardingService(
			usersDB.UsersRepo,
			matchingRepo,
			transporter,
			consumerMap,
		),
		Feedbacks: NewFeedbacksService(usersDB.FeedbackRepo),
	}
}

func (m Manager) InitRoute() error {
	rootRoute := m.App.Group("/")
	rootRoute.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello from Peakee Rest API")
	})

	authorizedWithoutUser := rootRoute.Group(
		"/users/self",

		auth.FiberAuthMiddleware(m.Auth, m.UsersRepo,
			auth.MiddlewareOptions{
				CheckUser: false,
			}),
	)
	authorizedWithoutUser.Get("/", m.Users.GetSelfFromAuth)
	authorizedWithoutUser.Post("/", m.Users.CreateNewUserBySelf)

	authorized := rootRoute.Group("/", auth.FiberAuthMiddleware(m.Auth, m.UsersRepo))

	users := authorized.Group("/users")
	users.Get("/", m.Users.GetUsers)
	users.Get(
		"/:id",
		ValidateUserIDParam(ValidateUserOptions{allowPublicQuery: true}),
		m.Users.GetUserByID,
	)
	users.Get("/:id/friend-requests",
		ValidateUserIDParam(),
		m.Users.GetPendingFriendRequests)
	users.Post("/:id/friend-requests",
		ValidateUserIDParam(),
		m.Users.CreateAddFriendRequest)
	users.Put(
		"/:id/friend-requests/:requestId",
		ValidateUserIDParam(),
		m.Users.RespondFriendRequest)

	// TODO: need to check if this user is in the conversation
	conversations := authorized.Group("/conversations")
	conversations.Get("/:id", m.Conversations.GetConversationByID)
	conversations.Get("/:id/messages", m.Conversations.GetMessagesOfConversation)
	conversations.Get("/", m.Conversations.GetConversationsOfUser)
	conversations.Post("/", m.Conversations.CreateNewIndividualConversation)

	authorized.Get("/messages/:id", m.Messages.GetMessageByID)

	authorized.Post("/onboarding", m.Onboardings.PostOnboardingForm())

	authorized.Post("/feedback", m.Feedbacks.CreateFeedback)

	return nil
}
