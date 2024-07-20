package chat

import (
	"log"
	"net/http"
	"strconv"

	"blinders/packages/auth"
	"blinders/packages/utils"
	"blinders/services/chat/repo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Auth         *auth.Manager
	ConvsRepo    *repo.ConversationsRepo
	MessagesRepo *repo.MessagesRepo
}

func NewService(
	auth *auth.Manager,
	db *mongo.Database,
) *Service {
	return &Service{
		Auth:         auth,
		ConvsRepo:    repo.NewConversationsRepo(db),
		MessagesRepo: repo.NewMessagesRepo(db),
	}
}

func (s Service) InitFiberRoutes(r fiber.Router) {
	// TODO: need to check if this user is in the conversation
	conversations := r.Group(
		"/conversations",
		s.Auth.FiberAuthMiddleware(auth.Config{WithUser: true}),
	)
	conversations.Get("/:id", s.GetConversationByID)
	conversations.Get("/:id/messages", s.GetMessagesOfConversation)
	conversations.Get("/", s.GetConversationsOfUser)
	conversations.Post("/", s.CreateNewIndividualConversation)
}

func (s Service) GetConversationByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid id: " + err.Error(),
		})
	}

	conversation, err := s.ConvsRepo.GetConversationByID(oid)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "can not get conversation:" + err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(conversation)
}

func (s Service) GetConversationsOfUser(ctx *fiber.Ctx) error {
	userID := ctx.Locals(auth.UserIDKey).(primitive.ObjectID)

	queryType := ctx.Query("type", "all")
	switch queryType {
	case "all":
		conversations, err := s.ConvsRepo.GetConversationByMembers(
			[]primitive.ObjectID{userID})
		if err != nil {
			log.Println("can not get conversations:", err)
			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"error": "can not get conversations",
			})
		}
		return ctx.Status(http.StatusOK).JSON(conversations)
	case "individual":
		friendID, err := primitive.ObjectIDFromHex(
			ctx.Query("friendId", ""))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"error": "friend id is required",
			})
		}
		conversations, err := s.ConvsRepo.GetConversationByMembers(
			[]primitive.ObjectID{userID, friendID},
			repo.IndividualConversation)
		if err != nil {
			log.Println("can not get conversations:", err)
			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"error": "can not get conversations",
			})
		}
		return ctx.Status(http.StatusOK).JSON(conversations)
	case "group":
		conversations, err := s.ConvsRepo.GetConversationByMembers(
			[]primitive.ObjectID{userID}, repo.GroupConversation)
		if err != nil {
			log.Println("can not get conversations:", err)
			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"error": "can not get conversations",
			})
		}
		return ctx.Status(http.StatusOK).JSON(conversations)
	default:
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid query type, must be 'all', 'group' or 'individual'",
		})
	}
}

type CreateConversationDTO struct {
	Type repo.ConversationType `json:"type"`
}

type CreateGroupConvDTO struct {
	CreateConversationDTO `json:",inline"`
}

type CreateIndividualConvDTO struct {
	CreateConversationDTO `       json:",inline"`
	FriendID              string `json:"friendId"`
}

func (s Service) CreateNewIndividualConversation(ctx *fiber.Ctx) error {
	convDTO, err := utils.ParseJSON[CreateConversationDTO](ctx.Body())
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid payload to create conversation",
		})
	}

	switch convDTO.Type {
	case repo.IndividualConversation:
		{
			convDTO, err := utils.ParseJSON[CreateIndividualConvDTO](ctx.Body())
			if err != nil {
				return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
					"error": "invalid payload to create individual conversation",
				})
			}

			userID := ctx.Locals(auth.UserIDKey).(primitive.ObjectID)

			friendID, err := primitive.ObjectIDFromHex(convDTO.FriendID)
			if err != nil {
				return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
					"error": "invalid friend id",
				})
			}

			conv, err := s.ConvsRepo.InsertIndividualConversation(userID, friendID)
			if err != nil {
				return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
					"error": err.Error(),
				})
			}

			return ctx.Status(http.StatusCreated).JSON(conv)

		}
	}

	return nil
}

func (s Service) GetMessagesOfConversation(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("invalid id:", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid id",
		})
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "30"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid limit",
		})
	}
	messages, err := s.MessagesRepo.GetMessagesOfConversation(oid, int64(limit))
	if err != nil {
		log.Println("can not get messages:", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "can not get messages",
		})
	}

	return ctx.Status(http.StatusOK).JSON(messages)
}
