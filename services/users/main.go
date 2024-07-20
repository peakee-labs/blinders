package users

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"blinders/packages/auth"
	"blinders/packages/utils"

	"blinders/services/users/repo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Auth               *auth.Manager
	UsersRepo          *repo.UsersRepo
	FriendRequestsRepo *repo.FriendRequestsRepo
}

func NewService(auth *auth.Manager, db *mongo.Database) *Service {
	return &Service{
		Auth:               auth,
		UsersRepo:          repo.NewUsersRepo(db),
		FriendRequestsRepo: repo.NewFriendRequestsRepo(db),
	}
}

func (s Service) InitFiberRoutes(r fiber.Router) {
	authorized := r.Group("/", s.Auth.FiberAuthMiddleware())
	authorized.Get("/", s.GetUsers)
	authorized.Get("/self", s.GetSelfFromAuth)
	authorized.Post("/self", s.CreateNewUserBySelf)

	authorized = r.Group("/", s.Auth.FiberAuthMiddleware(auth.Config{WithUser: true}))
	authorized.Get("/:id", ValidateUserIDParam(ValidateOptions{PublicQuery: true}), s.GetUserByID)
	authorized.Get("/:id/friend-requests", ValidateUserIDParam(), s.GetPendingFriendRequests)
	authorized.Post("/:id/friend-requests", ValidateUserIDParam(), s.CreateAddFriendRequest)
	authorized.Put("/:id/friend-requests/:requestId", ValidateUserIDParam(), s.RespondFriendRequest)
}

func (s Service) GetSelfFromAuth(ctx *fiber.Ctx) error {
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

func (s Service) GetUserByID(ctx *fiber.Ctx) error {
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

func (s Service) GetUsers(ctx *fiber.Ctx) error {
	email := ctx.Query("email", "")
	if email != "" {
		user, err := s.UsersRepo.GetUserByEmail(email)
		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.Status(http.StatusOK).JSON([]repo.User{user})
	}

	return nil
}

type CreateUserDTO struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

func (s Service) CreateNewUserBySelf(ctx *fiber.Ctx) error {
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

	user, err := s.UsersRepo.InsertNewRawUser(repo.User{
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

func (s Service) GetPendingFriendRequests(ctx *fiber.Ctx) error {
	userID, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		log.Println("invalid user id:", err)
		return err
	}

	requests, err := s.FriendRequestsRepo.GetFriendRequestByTo(
		userID,
		repo.FriendStatusPending,
	)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	if len(requests) == 0 {
		requests = make([]repo.FriendRequest, 0)
	}
	return ctx.Status(http.StatusOK).JSON(requests)
}

type AddFriendRequest struct {
	FriendID string `json:"friendId"`
}

func (s Service) CreateAddFriendRequest(ctx *fiber.Ctx) error {
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

	var user repo.User
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
		repo.FriendRequest{
			From:   userID,
			To:     friendID,
			Status: repo.FriendStatusPending,
		})
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
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

func (s Service) RespondFriendRequest(ctx *fiber.Ctx) error {
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

	var status repo.FriendRequestStatus
	switch payload.Action {
	case AcceptAddFriend:
		status = repo.FriendStatusAccepted
	case DenyAddFriend:
		status = repo.FriendStatusDenied
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

	return ctx.Status(http.StatusAccepted).JSON(request)
}
