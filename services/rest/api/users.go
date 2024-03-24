package restapi

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"blinders/packages/auth"
	"blinders/packages/db/models"
	"blinders/packages/db/repo"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersService struct {
	Repo *repo.UsersRepo
}

func NewUsersService(repo *repo.UsersRepo) *UsersService {
	return &UsersService{
		Repo: repo,
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
		Name:          userDTO.Name,
		Email:         userDTO.Email,
		ImageURL:      userDTO.ImageURL,
		FirebaseUID:   userAuth.AuthID,
		Conversations: make([]models.EmbeddedConversation, 0),
		FriendIDs:     make([]primitive.ObjectID, 0),
	})
	if err != nil {
		log.Println("can not create user:", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "can not create user",
		})
	}

	return ctx.Status(http.StatusCreated).JSON(user)
}
