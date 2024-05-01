package restapi

import (
	"log"
	"net/http"

	"blinders/packages/db/chatdb"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessagesService struct {
	Repo *chatdb.MessagesRepo
}

func NewMessagesService(repo *chatdb.MessagesRepo) *MessagesService {
	return &MessagesService{
		Repo: repo,
	}
}

func (s MessagesService) GetMessageByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("invalid id:", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid id",
		})
	}

	message, err := s.Repo.GetMessageByID(oid)
	if err != nil {
		log.Println("can not get message:", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "can not get message",
		})
	}

	return ctx.Status(http.StatusOK).JSON(message)
}
