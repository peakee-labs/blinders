package suggestapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/suggest"

	"github.com/gofiber/fiber/v2"
)

type Payload struct {
	Text   string `json:"text"`
	UserID string `json:"userID"`
}

func (s *Service) HandleTextSuggestion() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)

		req := new(Payload)
		if err := json.Unmarshal(ctx.Body(), req); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		userData, err := db.GetUserData(user.AuthID)
		if err != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"error": fmt.Sprintf("suggestion: cannot get data of user, err: (%s)", err.Error()),
			})
		}

		suggestions, err := s.Suggester.TextCompletion(ctx.Context(), userData, req.Text)
		if err != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"error":       err.Error(),
				"suggestions": []string{},
			})
		}

		return ctx.Status(200).JSON(fiber.Map{
			"suggestions": suggestions,
		})
	}
}

type ChatSuggestionPayload struct {
	UserID   string          `json:"userID"`
	Messages []ClientMessage `json:"messages"`
}

type ClientMessage struct {
	Timestamp any    `json:"time"`
	ID        string `json:"id"`
	Content   string `json:"content"`
	FromID    string `json:"senderId"`
	ChatID    string `json:"roomId"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
}

func (s *Service) HandleChatSuggestion() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := new(ChatSuggestionPayload)
		if err := json.Unmarshal(ctx.Body(), req); err != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"suggestions": []string{},
			})
		}

		// should communicate with user service
		userData, err := db.GetUserData(req.UserID)
		if err != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"suggestions": []string{},
			})
		}

		msgs := []suggest.Message{}
		for _, msg := range req.Messages {
			msgs = append(msgs, msg.ToCommonMessage())
		}

		suggestions, err := s.Suggester.ChatCompletion(ctx.Context(), userData, msgs)
		if err != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"suggestions": []string{},
			})
		}

		return ctx.Status(200).JSON(fiber.Map{
			"suggestions": suggestions,
		})
	}
}
