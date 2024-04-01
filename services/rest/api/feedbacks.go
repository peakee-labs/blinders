package restapi

import (
	"blinders/packages/db/models"
	"blinders/packages/db/repo"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedbacksService struct {
	Repo *repo.FeedbacksRepo
}

func NewFeedbacksService(repo *repo.FeedbacksRepo) *FeedbacksService {
	return &FeedbacksService{Repo: repo}
}

type CreateFeedbackDTO struct {
	Comment string `json:"comment"`
}

func (s FeedbacksService) CreateFeedback(ctx *fiber.Ctx) error {
	// TODO: currently accept request from anonymous user
	userIDParams := ctx.Params("id")
	userID, _ := primitive.ObjectIDFromHex(userIDParams)

	feedback, err := utils.ParseJSON[models.Feedback](ctx.Body())
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot unmarshal feedback from request body"})
	}
	feedback.UserID = userID
	_, err = s.Repo.InsertNewFeedback(*feedback)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot save feedback"})
	}
	return ctx.SendStatus(fiber.StatusOK)
}
