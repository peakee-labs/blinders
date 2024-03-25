package restapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"blinders/packages/db/models"
	"blinders/packages/transport"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s UsersService) GetPendingFriendRequests(ctx *fiber.Ctx) error {
	userID, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		log.Println("invalid user id:", err)
		return err
	}

	requests, err := s.FriendRequestsRepo.GetFriendRequestByTo(
		userID,
		models.FriendStatusPending,
	)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	if len(requests) == 0 {
		requests = make([]models.FriendRequest, 0)
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

	var user models.User
	err = s.Repo.Col.FindOne(context.Background(), bson.M{
		"_id":     userID,
		"friends": bson.M{"$all": []primitive.ObjectID{friendID}},
	}).Decode(&user)
	if err != mongo.ErrNoDocuments {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "user already added as friend",
		})
	}

	r, err := s.FriendRequestsRepo.InsertNewRawFriendRequest(
		models.FriendRequest{
			From:   userID,
			To:     friendID,
			Status: models.FriendStatusPending,
		})
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	event := transport.AddFriendEvent{
		Event:              transport.Event{Type: transport.AddFriend},
		UserID:             friendID.Hex(),
		AddFriendRequestID: r.ID.Hex(),
		Action:             transport.InitFriendRequest,
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

	var status models.FriendRequestStatus
	switch payload.Action {
	case AcceptAddFriend:
		status = models.FriendStatusAccepted
	case DenyAddFriend:
		status = models.FriendStatusDenied
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
		err = s.Repo.AddFriend(request.From, request.To)
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
		Event:              transport.Event{Type: transport.AddFriend},
		UserID:             request.From.Hex(),
		AddFriendRequestID: requestID.Hex(),
		Action:             action,
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
