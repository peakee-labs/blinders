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

	rsp, err := s.Transport.Request(ctx.Context(), s.ConsumerMap[transport.Suggest], requestJSON)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get practice unit"})
	}

	suggestLog, err := utils.ParseJSON[logging.SuggestLanguageUnitEvent](rsp)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get practice unit"})
	}

	// Pass this log to the logging service,
	if err := s.Transport.Push(ctx.Context(), s.ConsumerMap[transport.Logging], rsp); err != nil {
		log.Printf("suggest: cannot push practice unit suggest log to log service, err: %v", err)
	}
	// communicate with pypackage service to get suggestion for this user then save the document into db.

	// pass the response to client application
	return ctx.Status(fiber.StatusOK).JSON(suggestLog.Response)
}

//type Payload struct {
//	Text   string `json:"text"`
//	UserID string `json:"userID"`
//}

// HandleTextSuggestion is now deprecated
//func (s *Service) HandleTextSuggestion(ctx *fiber.Ctx) error {
//	user := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
//
//	req := new(Payload)
//	if err := json.Unmarshal(ctx.Body(), req); err != nil {
//		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
//			"message": err.Error(),
//		})
//	}
//
//	userData, err := db.GetUserData(user.AuthID)
//	if err != nil {
//		return ctx.Status(400).JSON(fiber.Map{
//			"error": fmt.Sprintf("suggestion: cannot get data of user, err: (%s)", err.Error()),
//		})
//	}
//
//	suggestions, err := s.Suggester.TextCompletion(ctx.Context(), userData, req.Text)
//	if err != nil {
//		return ctx.Status(400).JSON(fiber.Map{
//			"error":       err.Error(),
//			"suggestions": []string{},
//		})
//	}
//
//	return ctx.Status(200).JSON(fiber.Map{
//		"suggestions": suggestions,
//	})
//}

//type ChatSuggestionPayload struct {
//	UserID   string          `json:"userID"`
//	Messages []ClientMessage `json:"messages"`
//}

//type ClientMessage struct {
//	Timestamp any    `json:"time"`
//	ID        string `json:"id"`
//	Content   string `json:"content"`
//	FromID    string `json:"senderId"`
//	ChatID    string `json:"roomId"`
//	Sender    string `json:"sender"`
//	Receiver  string `json:"receiver"`
//}

// HandleChatSuggestion is now deprecated
//func (s *Service) HandleChatSuggestion(ctx *fiber.Ctx) error {
//	req := new(ChatSuggestionPayload)
//	if err := json.Unmarshal(ctx.Body(), req); err != nil {
//		return ctx.Status(400).JSON(fiber.Map{
//			"suggestions": []string{},
//		})
//	}
//
//	// should communicate with user service
//	userData, err := db.GetUserData(req.UserID)
//	if err != nil {
//		return ctx.Status(400).JSON(fiber.Map{
//			"suggestions": []string{},
//		})
//	}
//
//	var msgs []suggest.Message
//	for _, msg := range req.Messages {
//		msgs = append(msgs, msg.ToCommonMessage())
//	}
//
//	suggestions, err := s.Suggester.ChatCompletion(ctx.Context(), userData, msgs)
//	if err != nil {
//		return ctx.Status(400).JSON(fiber.Map{
//			"suggestions": []string{},
//		})
//	}
//
//	return ctx.Status(200).JSON(fiber.Map{
//		"suggestions": suggestions,
//	})
//}
