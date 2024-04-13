package suggestapi

import (
	"encoding/json"
	"log"

	"blinders/packages/auth"
	"blinders/packages/logging"
	"blinders/packages/transport"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"
)

// TODO: clarify that this handler
// - will get the language unit that suggested for user and return
// - make a new suggest request to pysuggest service then the content will by handle by the suggest service
// in 1st approach, this handle make no sense and the logic should put in pysuggest instead of this handler
// in 2nd approach, log should be pass to log service by the suggester
//
// current logic:
// this handler will be triggered when user want to retrieve 1 language unit (currently 1 word) that related to user's context
// then this handler will ask the pySuggester service to return 1 word that related to user context
// this handler pass the returned result from pySuggester service which include (suggestRequest and suggestResponse) to the log service in order to save user's suggested language unit (?)
// the suggestResponse then will be return to the client

func (s *Service) HandleSuggestLanguageUnit(ctx *fiber.Ctx) error {
	authUser := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if authUser == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user auth information"})
	}

	requestBody := map[string]string{
		"userID": authUser.ID,
	}

	requestJSON, _ := json.Marshal(requestBody)

	// communicate with pySuggest service to get suggestion for this user then save the document into db.
	rsp, err := s.Transport.Request(ctx.Context(), s.ConsumerMap[transport.Suggest], requestJSON)
	if err != nil {
		log.Println("cannot get response from py suggest service", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get practice unit"})
	}

	suggestLog, err := utils.ParseJSON[logging.SuggestPracticeUnitEvent](rsp)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get practice unit"})
	}

	// pass the log to the logging service
	event := logging.NewGenericEvent(logging.EventTypeSuggestPracticeUnit, *suggestLog)
	log.Println(event)

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Println("cannot marshal event", err)
	}

	if err := s.Transport.Push(ctx.Context(), s.ConsumerMap[transport.Logging], eventBytes); err != nil {
		log.Println("cannot push event to logging service", err)
	}

	// pass the response to client application
	return ctx.Status(fiber.StatusOK).JSON(suggestLog.Response)
}

func (s *Service) HandleGetPracticeUnit(ctx *fiber.Ctx) error {
	authUser := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if authUser == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user auth information"})
	}

	requestBody := map[string]string{
		"userID": authUser.ID,
	}

	requestJSON, _ := json.Marshal(requestBody)

	// communicate with pySuggest service to get suggestion for this user then save the document into db.
	rsp, err := s.Transport.Request(ctx.Context(), s.ConsumerMap[transport.Suggest], requestJSON)
	if err != nil {
		log.Println("cannot get response from py suggest service", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get practice unit"})
	}

	suggestLog, err := utils.ParseJSON[logging.SuggestPracticeUnitEvent](rsp)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get practice unit"})
	}

	// pass the log to the logging service
	if err := s.Transport.Push(ctx.Context(), s.ConsumerMap[transport.Logging], rsp); err != nil {
		log.Println("cannot push event to logging service", err)
	}

	// pass the response to client application
	return ctx.Status(fiber.StatusOK).JSON(suggestLog.Response)
}
