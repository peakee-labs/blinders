package practiceapi

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/collecting"
	"blinders/packages/transport"

	"github.com/gofiber/fiber/v2"
)

var DefaultLanguageLocale = "en"

func (s Service) HandleGetLanguageUnit(ctx *fiber.Ctx) error {
	authUser := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if authUser == nil {
		return fmt.Errorf("cannot get user auth information")
	}

	var (
		numReturn = 1
		req       = transport.GetEventRequest{
			Request:   transport.Request{Type: transport.GetEvent},
			UserID:    authUser.ID,
			NumReturn: numReturn,
			Type:      collecting.EventTypeExplain,
		}
		rsp = new(transport.GetEventResponse)
	)

	transportBytes, _ := json.Marshal(req)

	response, err := s.Transport.Request(
		ctx.Context(),
		s.ConsumerMap[transport.CollectingGet],
		transportBytes,
	)

	if err != nil {
		log.Printf("practice: cannot get log event from collecting service, err: %v\n", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get practice unit"})
	}

	if err := json.Unmarshal(response, &rsp); err != nil {
		log.Printf("practice: cannot parse result from collecting service, err: %v\n", err)
	}

	if len(rsp.Data) != numReturn {
		log.Printf("practice: expected return %v event, got %v\n", numReturn, rsp)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get practice unit"})
	}

	event := rsp.Data[0]
	switch event.Type {
	case collecting.EventTypeExplain:
		return ctx.Status(fiber.StatusOK).JSON(event.Payload) //ExplainEvent

	default:
		log.Printf("practice: unsupported event type (%v)\n", event.Type)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get practice unit"})
	}
}

func (s Service) HandleGetRandomLanguageUnit(ctx *fiber.Ctx) error {
	localeCode := ctx.Query("lang")
	unitType := ctx.Query("type", string(collecting.EventTypeExplain))

	// event type currently is capitialized
	switch collecting.EventType(strings.ToUpper(unitType)) {
	case collecting.EventTypeExplain:
		unit, err := s.GetRandomExplainWithLangCode(localeCode)
		if err != nil {
			// use pre-defined language tag as default language tag
			unit, _ = s.GetRandomExplainWithLangCode(DefaultLanguageLocale)
		}
		return ctx.Status(fiber.StatusOK).JSON(unit)

	case collecting.EventTypeSuggestPracticeUnit:
		unit, err := s.GetRandomPracticeUnitWithLangCode(localeCode)
		if err != nil {
			// use pre-defined language tag as default language tag
			unit, _ = s.GetRandomPracticeUnitWithLangCode(DefaultLanguageLocale)
		}
		return ctx.Status(fiber.StatusOK).JSON(unit)

	default:
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid unit type"})
	}
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

// GetRandomExplainWithLangCode returns random practice-unit with given langCode
func (s Service) GetRandomExplainWithLangCode(langCode string) (collecting.ExplainEvent, error) {
	// user's learning language code with RFC-5646 format
	explainUnit, ok := DefaultExplain[langCode]
	if !ok {
		return collecting.ExplainEvent{}, fmt.Errorf("invalid lang code")
	}

	idx := rand.Intn(len(explainUnit))
	return explainUnit[idx], nil
}
