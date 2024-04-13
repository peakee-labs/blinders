package suggestapi

import (
	"math/rand"

	"blinders/packages/auth"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Service) HandleSuggestLanguageUnit(ctx *fiber.Ctx) error {
	authUser := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if authUser == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user auth information"})
	}

	userOID, err := primitive.ObjectIDFromHex(authUser.ID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	loggedEvent, err := s.Logger.GetSuggestPracticeUnitEventLogByUserID(userOID)
	if err != nil {
		// we could return some pre-defined document here.
		return s.HandleGetDefaultLanguageUnit(ctx)
	}
	// currently, randomly return practice unit to user
	idx := rand.Intn(len(loggedEvent))
	return ctx.Status(fiber.StatusOK).JSON(loggedEvent[idx].Response)
}

// HandleGetDefaultLanguageUnit returns 1 random pre-defined LanguageUnit.
func (s *Service) HandleGetDefaultLanguageUnit(ctx *fiber.Ctx) error {
	return nil
}
