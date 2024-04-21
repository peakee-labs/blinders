package collectingapi

import (
	"log"

	"blinders/packages/collecting"
	"blinders/packages/transport"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Service) HandlePushEvent(ctx *fiber.Ctx) error {
	req, err := utils.ParseJSON[transport.CollectEventRequest](ctx.Body())
	if err != nil {
		log.Printf("cannot get collect event from request's body, err: %v\n", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get event from body"})
	}

	if req.Request.Type != transport.CollectEvent {
		log.Printf("event type mismatch, type: %v\n", req.Request.Type)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get event from body"})
	}

	if req.Type != transport.CollectEvent {
		log.Printf("invalid request type: %v", req.Type)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request type"})
	}

	event := req.Data
	switch event.Type {
	case collecting.EventTypeSuggestPracticeUnit:
		event, err := utils.JSONConvert[collecting.SuggestPracticeUnitEvent](event.Payload)
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
		event, err := utils.JSONConvert[collecting.TranslateEvent](event.Payload)
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

	case collecting.EventTypeExplain:
		event, err := utils.JSONConvert[collecting.ExplainEvent](event.Payload)
		if err != nil {
			log.Printf("cannot get explain event from payload")
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "mismatch event type and event payload"})
		}

		eventLog, err := s.Collector.AddRawExplainLog(&collecting.ExplainEventLog{
			ExplainEvent: *event,
		})
		if err != nil {
			log.Printf("logger: cannot add raw explain log, err: %v", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot append explain log"})
		}

		return ctx.Status(fiber.StatusOK).JSON(eventLog)

	default:
		log.Printf("receive unsupport event, type: %v", event.Type)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unsupported event type: " + event.Type,
		})
	}
}

func (s Service) HandleGetEvent(ctx *fiber.Ctx) error {
	req, err := utils.ParseJSON[transport.GetEventRequest](ctx.Body())
	if err != nil {
		log.Println("collector: cannot get event request from body")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get request body"})
	}

	userOID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get userID from request"})
	}

	switch req.Type {
	case collecting.EventTypeExplain:
		logs, err := s.Collector.GetExplainLogByUserID(userOID)
		if err != nil {
			log.Println("collecting: cannot get logs of user", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user's log event"})
		}

		return ctx.Status(fiber.StatusOK).JSON(logs[:req.NumReturn])

	case collecting.EventTypeTranslate:
		logs, err := s.Collector.GetTranslateLogByUserID(userOID)
		if err != nil {
			log.Println("collecting: cannot get logs of user", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user's log event"})
		}

		return ctx.Status(fiber.StatusOK).JSON(logs[:req.NumReturn])

	case collecting.EventTypeSuggestPracticeUnit:
		logs, err := s.Collector.GetSuggestPracticeUnitLogByUserID(userOID)
		if err != nil {
			log.Println("collecting: cannot get logs of user", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user's log event"})
		}

		return ctx.Status(fiber.StatusOK).JSON(logs[:req.NumReturn])

	default:
		log.Printf("collecting: received undefined event type (%v)\n", req.Type)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unsupported event"})
	}
}
