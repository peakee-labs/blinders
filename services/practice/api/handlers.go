package practiceapi

import (
	"log"
	"math/rand"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/logging"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var DefaultLanguageLocale = "en"

func (s Service) HandleSuggestLanguageUnit(ctx *fiber.Ctx) error {
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
		rsp = s.GetRandomPracticeUnitForUser(userOID)
	} else {
		// currently, randomly return practice unit to user
		idx := rand.Intn(len(loggedEvent))
		rsp = loggedEvent[idx].Response
	}

	return ctx.Status(fiber.StatusOK).JSON(rsp)
}

func (s Service) HandleGetRandomLanguageUnit(ctx *fiber.Ctx) error {
	localeCode := ctx.Query("lang")
	unit, err := s.GetRandomPracticeUnitWithLangCode(localeCode)
	if err != nil {
		// use pre-defined language tag as default language tag
		unit, err = s.GetRandomPracticeUnitWithLangCode(DefaultLanguageLocale)
	}
	return ctx.Status(fiber.StatusOK).JSON(unit)
}

// GetRandomPracticeUnitWithLangCode return random practiceunit with given langCode
func (s Service) GetRandomPracticeUnitWithLangCode(langCode string) (logging.SuggestPracticeUnitResponse, error) {
	// user's learning language code with RFC-5646 format
	units, ok := DefaultLanguageUnit[langCode]
	if !ok {
		return logging.SuggestPracticeUnitResponse{}, fmt.Errorf("language unit with given language is not existed")
	}

	idx := rand.Intn(len(units))
	return units[idx], nil
}
