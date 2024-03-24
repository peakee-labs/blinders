package exploreapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"blinders/packages/auth"
	"blinders/packages/db/models"
	"blinders/packages/explore"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Core        explore.Explorer
	RedisClient *redis.Client
	MatchCh     chan models.MatchInfo
}

func NewService(
	exploreCore explore.Explorer,
	redisClient *redis.Client,
) *Service {
	return &Service{
		Core:        exploreCore,
		RedisClient: redisClient,
		MatchCh:     make(chan models.MatchInfo),
	}
}

// HandleGetMatches returns 5 users that similarity with current user.
func (s *Service) HandleGetMatches(ctx *fiber.Ctx) error {
	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok || userAuth == nil {
		log.Println("cannot get auth user")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get user"})
	}

	candidates, err := s.Core.Suggest(userAuth.ID)
	if err != nil {
		log.Println("cannot get suggest for user", userAuth.ID, "err", err)
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(candidates)
}

// HandleAddUserMatch will add match-related information to match db
func (s *Service) HandleAddUserMatch(ctx *fiber.Ctx) error {
	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok || userAuth == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get user"})
	}

	userMatch := new(models.MatchInfo)
	if err := json.Unmarshal(ctx.Body(), userMatch); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "service: match information required in body",
		})
	}

	// currently user.AuthID is firebaseUID
	if userMatch.UserID.Hex() != userAuth.ID {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Insufficient permissions. The requester and the user in the request body must match.",
		})
	}

	info, err := s.Core.AddUserMatchInformation(*userMatch)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Errorf("service: cannot add user information, %v", err).Error(),
		})
	}
	s.MatchCh <- info
	return ctx.Status(fiber.StatusOK).JSON(info)
}

func (s *Service) Loop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	for match := range s.MatchCh {
		value := map[string]string{
			"id": match.UserID.Hex(),
		}

		// fire user-match-entry-created event to the stream
		if err := s.RedisClient.XAdd(ctx, &redis.XAddArgs{
			Stream:     "match:embed",
			NoMkStream: false,
			Values:     value,
		}).Err(); err != nil {
			fmt.Println("error", err)
			break
		}
		fmt.Println("fired event to stream", value)
	}
}
