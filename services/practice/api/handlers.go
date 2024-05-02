package practiceapi

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/transport"

	"github.com/gofiber/fiber/v2"
)

func (s Service) HandleGetPracticeUnitFromAnalyzeExplainLog(ctx *fiber.Ctx) error {
	authUser := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if authUser == nil {
		return fmt.Errorf("cannot get user auth information")
	}

	req := transport.GetCollectingLogRequest{
		Request: transport.Request{Type: transport.GetExplainLog},
		UserID:  authUser.ID,
	}

	reqBytes, _ := json.Marshal(req)
	response, err := s.Transport.Request(
		ctx.Context(),
		s.ConsumerMap[transport.CollectingGet],
		reqBytes,
	)
	if err != nil {
		log.Println("cannot get explain log:", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "cannot get explain log"})
	}

	var jsonResponse map[string]any
	_ = json.Unmarshal(response, &jsonResponse)

	return ctx.Status(http.StatusOK).JSON(jsonResponse)
}

func (s Service) HandleGetRandomLanguageUnit(ctx *fiber.Ctx) error {
	unitType := ctx.Query("type", "DEFAULT")

	switch strings.ToUpper(unitType) {
	case "DEFAULT":
		idx := rand.Intn(len(DefaultSimplePracticeUnits))
		unit := DefaultSimplePracticeUnits[idx]
		return ctx.Status(fiber.StatusOK).JSON(unit)

	case "EXPLAIN":
		idx := rand.Intn(len(ExplainLogSamples))
		unit := ExplainLogSamples[idx]
		return ctx.Status(fiber.StatusOK).JSON(unit)

	default:
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid unit type"})
	}
}
