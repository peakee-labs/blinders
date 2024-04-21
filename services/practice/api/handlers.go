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
		return fmt.Errorf("cannot get user auth information")
	}

	var (
		req = transport.GetEventRequest{
			Request:   transport.Request{Type: transport.GetEvent},
			UserID:    authUser.ID,
			NumReturn: 1,
			Type:      collecting.EventTypeSuggestPracticeUnit,
		}
		rsp = new(transport.GetEventResponse)
	)

	transportBytes, _ := json.Marshal(req)

	response, err := s.Transport.Request(
		ctx.Context(),
		s.ConsumerMap[transport.Collecting],
		transportBytes,
	)
	if err != nil {
		log.Printf("practice: cannot get log event from collecting service, err: %v\n", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get practice unit"})
	}

	if err := json.Unmarshal(response, &rsp); err != nil {
		log.Printf("practice: cannot parse result from collecting service, err: %v\n", err)
	}
	if len(rsp.Data) == 0 {
		log.Printf("practice: response includes no event: \n")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get practice unit"})
	}

	event := rsp.Data[0]
	switch event.Type {
	case collecting.EventTypeSuggestPracticeUnit:
		return ctx.Status(fiber.StatusOK).JSON(event.Payload)

	default:
		log.Printf("practice: unsupported event type (%v)\n", event.Type)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get practice unit"})
	}
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
