package restapi

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"blinders/packages/auth"
	"blinders/packages/db/models"
	"blinders/packages/db/repo"
	"blinders/packages/transport"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersService struct {
	Repo               *repo.UsersRepo
	FriendRequestsRepo *repo.FriendRequestsRepo
	Transporter        transport.Transport
	ConsumerMap        transport.ConsumerMap
}

func NewUsersService(
	repo *repo.UsersRepo,
	frRepo *repo.FriendRequestsRepo,
	transporter transport.Transport,
	consumerMap transport.ConsumerMap,
) *UsersService {
	return &UsersService{
		Repo:               repo,
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

	user, err := s.Repo.GetUserByFirebaseUID(userAuth.AuthID)
	if errors.Is(err, mongo.ErrNoDocuments) {
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

	user, err := s.Repo.GetUserByID(oid)
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
		user, err := s.Repo.GetUserByEmail(email)
		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.Status(http.StatusOK).JSON([]models.User{user})
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

	user, err := s.Repo.InsertNewRawUser(models.User{
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
