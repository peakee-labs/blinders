package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"blinders/packages/auth"
	"blinders/packages/db/usersdb"
	"blinders/packages/transport"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersService struct {
	UsersRepo          *usersdb.UsersRepo
	FriendRequestsRepo *usersdb.FriendRequestsRepo
	Transporter        transport.Transport
	ConsumerMap        transport.ConsumerMap
}

func NewUsersService(
	repo *usersdb.UsersRepo,
	frRepo *usersdb.FriendRequestsRepo,
	transporter transport.Transport,
	consumerMap transport.ConsumerMap,
) *UsersService {
	return &UsersService{
		UsersRepo:          repo,
		FriendRequestsRepo: frRepo,
		Transporter:        transporter,
		ConsumerMap:        consumerMap,
	}
}

func (s UsersService) GetSelfFromAuth(ctx *fiber.Ctx) error {
	userAuth := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if userAuth == nil {
		return fmt.Errorf("required user auth")
	}

	user, err := s.UsersRepo.GetUserByFirebaseUID(userAuth.AuthID)
	if err == mongo.ErrNoDocuments {
		return ctx.Status(http.StatusNotFound).JSON(nil)
	} else if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(user)
}

func (s UsersService) GetUserByID(ctx *fiber.Ctx) error {
	// TODO: need to check if this is a public query and eliminate private data
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("invalid id:", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid id",
		})
	}

	user, err := s.UsersRepo.GetUserByID(oid)
	if err != nil {
		log.Println("can not get user:", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "can not get user",
		})
	}

	return ctx.Status(http.StatusOK).JSON(user)
}

func (s UsersService) GetUsers(ctx *fiber.Ctx) error {
	email := ctx.Query("email", "")
	if email != "" {
		user, err := s.UsersRepo.GetUserByEmail(email)
		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.Status(http.StatusOK).JSON([]usersdb.User{user})
	}

	return nil
}

type CreateUserDTO struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

func (s UsersService) CreateNewUserBySelf(ctx *fiber.Ctx) error {
	userDTO, err := utils.ParseJSON[CreateUserDTO](ctx.Body())
	if err != nil {
		log.Println("invalid payload:", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid payload",
		})
	}
	if userDTO.Email == "" || userDTO.Name == "" {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid payload, require email and name",
		})
	}

	userAuth := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if userAuth == nil {
		return fmt.Errorf("required user auth")
	}

	user, err := s.UsersRepo.InsertNewRawUser(usersdb.User{
		Name:        userDTO.Name,
		Email:       userDTO.Email,
		ImageURL:    userDTO.ImageURL,
		FirebaseUID: userAuth.AuthID,
		FriendIDs:   make([]primitive.ObjectID, 0),
	})
	if err != nil {
		log.Println("can not create user:", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "can not create user",
		})
	}

	return ctx.Status(http.StatusCreated).JSON(user)
}

func (s UsersService) GetPendingFriendRequests(ctx *fiber.Ctx) error {
	userID, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		log.Println("invalid user id:", err)
		return err
	}

	requests, err := s.FriendRequestsRepo.GetFriendRequestByTo(
		userID,
		usersdb.FriendStatusPending,
	)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	if len(requests) == 0 {
		requests = make([]usersdb.FriendRequest, 0)
	}
	return ctx.Status(http.StatusOK).JSON(requests)
}

type AddFriendRequest struct {
	FriendID string `json:"friendId"`
}

func (s UsersService) CreateAddFriendRequest(ctx *fiber.Ctx) error {
	userID, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		log.Println("invalid user id:", err)
		return err
	}

	payload, err := utils.ParseJSON[AddFriendRequest](ctx.Body())
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid payload",
		})
	}
	friendID, err := primitive.ObjectIDFromHex(payload.FriendID)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid friend id",
		})
	}

	var user usersdb.User
	err = s.UsersRepo.FindOne(context.Background(), bson.M{
		"_id":     userID,
		"friends": bson.M{"$all": []primitive.ObjectID{friendID}},
	}).Decode(&user)
	if err != mongo.ErrNoDocuments {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "user already added as friend",
		})
	}

	r, err := s.FriendRequestsRepo.InsertNewRawFriendRequest(
		usersdb.FriendRequest{
			From:   userID,
			To:     friendID,
			Status: usersdb.FriendStatusPending,
		})
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	event := transport.AddFriendEvent{
		Event: transport.Event{Type: transport.AddFriend},
		Payload: transport.AddFriendPayload{
			UserID:             friendID.Hex(),
			AddFriendRequestID: r.ID.Hex(),
			Action:             transport.InitFriendRequest,
		},
	}
	notiPayload, _ := json.Marshal(event)
	err = s.Transporter.Push(
		context.Background(),
		s.ConsumerMap[transport.Notification],
		notiPayload,
	)
	if err != nil {
		log.Println("failed to push notification", err)
	}

	return ctx.Status(http.StatusCreated).JSON(r)
}

const (
	AcceptAddFriend string = "accept"
	DenyAddFriend   string = "deny"
)

type RespondFriendRequest struct {
	Action string `json:"action"`
}

func (s UsersService) RespondFriendRequest(ctx *fiber.Ctx) error {
	userID, _ := primitive.ObjectIDFromHex(ctx.Params("id"))
	requestID, err := primitive.ObjectIDFromHex(ctx.Params("requestId"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid request id",
		})
	}

	payload, err := utils.ParseJSON[RespondFriendRequest](ctx.Body())
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid payload",
		})
	}

	var status usersdb.FriendRequestStatus
	switch payload.Action {
	case AcceptAddFriend:
		status = usersdb.FriendStatusAccepted
	case DenyAddFriend:
		status = usersdb.FriendStatusDenied
	default:
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid action",
		})
	}

	request, err := s.FriendRequestsRepo.UpdateFriendRequestStatusByID(
		requestID,
		userID,
		status,
	)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	var action transport.AddFriendAction
	switch payload.Action {
	case AcceptAddFriend:
		// TODO: need to apply transaction
		err = s.UsersRepo.AddFriend(request.From, request.To)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"error": err.Error(),
			})
		}
		action = transport.AcceptFriendRequest
	case DenyAddFriend:
		action = transport.DenyFriendRequest
	}

	event := transport.AddFriendEvent{
		Event: transport.Event{Type: transport.AddFriend},
		Payload: transport.AddFriendPayload{
			UserID:             request.From.Hex(),
			AddFriendRequestID: requestID.Hex(),
			Action:             action,
		},
	}
	notiPayload, _ := json.Marshal(event)
	err = s.Transporter.Push(
		context.Background(),
		s.ConsumerMap[transport.Notification],
		notiPayload,
	)
	if err != nil {
		log.Println("failed to push notification", err)
	}

	return ctx.Status(http.StatusAccepted).JSON(request)
}
