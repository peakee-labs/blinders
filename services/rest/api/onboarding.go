package restapi

import (
	"encoding/json"
	"fmt"
	"os"

	"blinders/packages/auth"
	"blinders/packages/db/models"
	"blinders/packages/db/repo"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	OnboardingService struct {
		UserRepo    *repo.UsersRepo
		ExploreRepo *repo.MatchesRepo
	}
	OnboardingForm struct {
		Name      string   `json:"name"      form:"name"`
		Major     string   `json:"major"     form:"major"`
		Gender    string   `json:"gender"    form:"gender"`
		Native    string   `json:"native"    form:"native"`
		Country   string   `json:"country"   form:"country"`
		Learnings []string `json:"learnings" form:"learnings"`
		Interests []string `json:"interests" form:"interests"`
		Age       int      `json:"age"       form:"age"`
	}
)

func (s *OnboardingService) PostOnboardingForm() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
		if !ok || userAuth == nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get user"})
		}

		var formValue OnboardingForm
		if err := ctx.BodyParser(&formValue); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
		}
		uid, err := primitive.ObjectIDFromHex(userAuth.ID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get objectID from userID " + err.Error()})
		}
		matchInfo, err := utils.JSONConvert[models.MatchInfo](formValue)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot unmarshal match info from form value" + err.Error()})
		}
		matchInfo.UserID = uid

		// TODO: at here we must notify explore service to add new profile to explore db as well as vector db.
		// TODO: currently, we make a inter-system http request.
		embedderURL := fmt.Sprintf("http://%s:%s/explore", os.Getenv("EXPLORE_API_HOST"), os.Getenv("EXPLORE_API_PORT"))
		jsonBody, err := json.Marshal(matchInfo)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Errorf("service: cannot add user information, %v", err).Error(),
			})
		}
		code, _, errs := fiber.Post(embedderURL).Body(jsonBody).Bytes()
		if errs != nil || len(errs) > 0 {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Errorf("service: cannot get embed vector, %v", errs).Error(),
			})
		}
		if code != fiber.StatusOK {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "service: cannot get embed vector",
			})
		}
		return nil
	}
}

func NewOnboardingService(userRepo *repo.UsersRepo, matchRepo *repo.MatchesRepo) *OnboardingService {
	return &OnboardingService{
		UserRepo:    userRepo,
		ExploreRepo: matchRepo,
	}
}
