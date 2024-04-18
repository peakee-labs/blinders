package practiceapi

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/collecting"

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

	var rsp collecting.SuggestPracticeUnitResponse
	loggedEvent, err := s.Logger.GetSuggestPracticeUnitLogByUserID(userOID)
	if err != nil || len(loggedEvent) == 0 {
		// try to return some pre-defined document here.
		log.Printf("practice: cannot get log event from Logger, err: %v, event num: %v\n", err, len(loggedEvent))
		usr, err := s.Db.Matches.GetMatchInfoByUserID(userOID)
		if err != nil {
			rsp, _ = s.GetRandomPracticeUnitWithLangCode(DefaultLanguageLocale)
			goto response
		}

		// we could index the most 'active' language with index 0 that mark that specific is currently actively learning by user.
		lang := strings.Split(usr.Learnings[0], "-")[0] // we only take the Two-character primary language subtags (ex: "en-US" => "en")
		rsp, err = s.GetRandomPracticeUnitWithLangCode(lang)
		if err != nil {
			// use pre-defined language tag as default language tag
			rsp, _ = s.GetRandomPracticeUnitWithLangCode(DefaultLanguageLocale)
		}
	} else {
		// currently, randomly return practice unit to user
		idx := rand.Intn(len(loggedEvent))
		rsp = loggedEvent[idx].Response
	}

response:
	return ctx.Status(fiber.StatusOK).JSON(rsp)
}

func (s Service) HandleGetRandomLanguageUnit(ctx *fiber.Ctx) error {
	localeCode := ctx.Query("lang")
	unit, err := s.GetRandomPracticeUnitWithLangCode(localeCode)
	if err != nil {
		// use pre-defined language tag as default language tag
		unit, _ = s.GetRandomPracticeUnitWithLangCode(DefaultLanguageLocale)
	}
	return ctx.Status(fiber.StatusOK).JSON(unit)
}

// GetRandomPracticeUnitWithLangCode return random practiceunit with given langCode
func (s Service) GetRandomPracticeUnitWithLangCode(langCode string) (collecting.SuggestPracticeUnitResponse, error) {
	// user's learning language code with RFC-5646 format
	units, ok := DefaultLanguageUnit[langCode]
	if !ok {
		return collecting.SuggestPracticeUnitResponse{}, fmt.Errorf("language unit with given language is not existed")
	}

	idx := rand.Intn(len(units))
	return units[idx], nil
}
