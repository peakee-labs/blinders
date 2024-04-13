package suggestapi

import (
	"log"
	"math/rand"

	"blinders/packages/auth"
	"blinders/packages/logging"

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

	var rsp logging.SuggestPracticeUnitResponse
	loggedEvent, err := s.Logger.GetSuggestPracticeUnitEventLogByUserID(userOID)
	if err != nil || len(loggedEvent) == 0 {
		log.Printf("practice: cannot get log event from Logger, err: %v, event num: %v\n", err, len(loggedEvent))
		// we could return some pre-defined document here.
		idx := rand.Intn(len(DefaultLanguageUnit))
		rsp = DefaultLanguageUnit[idx]
	} else {
		// currently, randomly return practice unit to user
		idx := rand.Intn(len(loggedEvent))
		rsp = loggedEvent[idx].Response
	}

	return ctx.Status(fiber.StatusOK).JSON(rsp)
}
