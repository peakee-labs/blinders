package collectingapi

import (
	"fmt"
	"log"

	"blinders/packages/auth"
	"blinders/packages/collecting"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Service) HandlePushEvent(ctx *fiber.Ctx) error {
	e, err := utils.ParseJSON[collecting.GenericEvent](ctx.Body())
	if err != nil {
		log.Printf("cannot get generic event from request's body, err: %v\n", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get log from body"})
	}

	switch e.Type {
	case collecting.EventTypeSuggestPracticeUnit:
		event, err := utils.JSONConvert[collecting.SuggestPracticeUnitEvent](e.Payload)
		if err != nil {
			log.Printf("cannot get suggest practice unit event from payload")
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "mismatch event type and event payload"})
		}

		eventLog, err := s.Collector.AddRawSuggestPracticeUnitLog(&collecting.SuggestPracticeUnitEventLog{
			SuggestPracticeUnitEvent: *event,
		})
		if err != nil {
			log.Printf("logger: cannot add raw translate log, err: %v", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot append translate log"})
		}

		return ctx.Status(fiber.StatusOK).JSON(eventLog)
	case collecting.EventTypeTranslate:
		event, err := utils.JSONConvert[collecting.TranslateEvent](e.Payload)
		if err != nil {
			log.Printf("cannot get translate event from payload")
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "mismatch event type and event payload"})
		}

		eventLog, err := s.Collector.AddRawTranslateLog(&collecting.TranslateEventLog{
			TranslateEvent: *event,
		})
		if err != nil {
			log.Printf("logger: cannot add raw translate log, err: %v", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot append translate log"})
		}

		return ctx.Status(fiber.StatusOK).JSON(eventLog)
	default:
		log.Printf("receive undefined event")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unsupported event type: " + e.Type})
	}
}

func (s Service) HandleGetEvent(ctx *fiber.Ctx) error {
	authUser, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Println("collector:  cannot get auth user from fiber context")
		return fmt.Errorf("cannot get auth user")
	}
	userOID, err := primitive.ObjectIDFromHex(authUser.ID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get userID from request"})
	}
	logs, err := s.Collector.GetTranslateLogByUserID(userOID)
	if err != nil {
		log.Println("collecting: cannot get logs of user", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user's log event"})
	}

	return ctx.Status(fiber.StatusOK).JSON(logs)
}
