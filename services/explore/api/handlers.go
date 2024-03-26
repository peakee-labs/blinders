package exploreapi

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/db/models"
	"blinders/packages/explore"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Core        explore.Explorer
	RedisClient *redis.Client
}

func NewService(
	exploreCore explore.Explorer,
	redisClient *redis.Client,
) *Service {
	return &Service{
		Core:        exploreCore,
		RedisClient: redisClient,
	}
}

// HandleGetMatches returns 5 users that similarity with current user.
func (s *Service) HandleGetMatches(ctx *fiber.Ctx) error {
	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok || userAuth == nil {
		log.Println("cannot get auth user")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user"})
	}

	candidates, err := s.Core.Suggest(userAuth.ID)
	if err != nil {
		log.Println("cannot get suggest for user", userAuth.ID, "err", err)
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(candidates)
}

// HandleAddUserMatch will add match-related information to match db, this api could only be called by inter-system service
func (s *Service) HandleAddUserMatch(ctx *fiber.Ctx) error {
	userMatch := new(models.MatchInfo)
	if err := json.Unmarshal(ctx.Body(), userMatch); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "service: match information required in body",
		})
	}

	// TODO: at here, we've to call embedder service to get the embed vector then add this entry into db.
	embedderURL := fmt.Sprintf("http://%s:%s/explore/embed", os.Getenv("EXPLORE_EMBEDDER_HOST"), os.Getenv("EXPLORE_EMBEDDER_PORT"))
	jsonBody, err := json.Marshal(userMatch)
	if err != nil {
		log.Println("cannot unmarshal body", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Errorf("service: cannot add user information, %v", err).Error(),
		})
	}
	fmt.Println(string(jsonBody))

	code, body, errs := fiber.Post(embedderURL).ContentType("application/json").Body(jsonBody).Bytes()
	if errs != nil || len(errs) > 0 {
		log.Println("cannot request embed vector", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Errorf("service: cannot get embed vector, %v", errs).Error(),
		})
	}
	if code != fiber.StatusOK {
		log.Println("cannot request embed vector, server response", code)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "service: cannot get embed vector",
		})
	}

	type response struct {
		Embed explore.EmbeddingVector `json:"embed"`
	}
	var rsp response
	if err := json.Unmarshal(body, &rsp); err != nil {
		log.Println("cannot get embed from server response", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Errorf("service: cannot unmarshall embed vector, err: %v", err),
		})
	}

	info, err := s.Core.AddUserMatchInformation(*userMatch)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Errorf("service: cannot add user information, %v", err).Error(),
		})
	}

	if err := s.Core.AddEmbedding(info.UserID, rsp.Embed); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Errorf("service: cannot add user embed, %v", err).Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(info)
}
