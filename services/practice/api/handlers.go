package practiceapi

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	"blinders/packages/auth"
	"blinders/packages/collecting"
	"blinders/packages/transport"

	"github.com/gofiber/fiber/v2"
)

var DefaultLanguageLocale = "en"

func (s Service) HandleSuggestLanguageUnit(ctx *fiber.Ctx) error {
	authUser := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if authUser == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user auth information"})
	}

	var unit collecting.SuggestPracticeUnitResponse
	transportBody := map[string]string{
		"userID": authUser.ID,
	}
	transportBytes, _ := json.Marshal(transportBody)

	response, err := s.Transport.Request(ctx.Context(), s.ConsumerMap[transport.Collecting], transportBytes)
	if err != nil {
		log.Printf("practice: cannot get log event from collecting service, err: %v\n", err)

		preferedLang := ctx.Query("lang")
		if preferedLang != "" {
			unit, err = s.GetRandomPracticeUnitWithLangCode(preferedLang)
			if err == nil {
				goto response
			}
			log.Printf("practice: cannot get predefined event, lang: %v, err: %v\n", preferedLang, err)
		}
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if err := json.Unmarshal(response, &unit); err != nil {
		log.Printf("practice: cannot parse result from collecting service, err: %v\n", err)
	}

response:
	return ctx.Status(fiber.StatusOK).JSON(unit)
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

// GetRandomPracticeUnitWithLangCode returns random practice-unit with given langCode
func (s Service) GetRandomPracticeUnitWithLangCode(langCode string) (collecting.SuggestPracticeUnitResponse, error) {
	// user's learning language code with RFC-5646 format
	units, ok := DefaultLanguageUnit[langCode]
	if !ok {
		return collecting.SuggestPracticeUnitResponse{}, fmt.Errorf("language unit with given language is not existed")
	}

	idx := rand.Intn(len(units))
	return units[idx], nil
}
