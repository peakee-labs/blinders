package exploreapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/db/matchingdb"
	"blinders/packages/explore"
	"blinders/packages/transport"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	Core             explore.Explorer
	RedisClient      *redis.Client
	EmbedderEndpoint string
}

func NewService(
	exploreCore explore.Explorer,
	redisClient *redis.Client,
	embedderEndpoint string,
) *Service {
	return &Service{
		Core:             exploreCore,
		RedisClient:      redisClient,
		EmbedderEndpoint: embedderEndpoint,
	}
}

// HandleGetMatches returns 5 users that similarity with current user.
func (s *Service) HandleGetMatches(ctx *fiber.Ctx) error {
	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok || userAuth == nil {
		log.Println("cannot get auth user")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user"})
	}
	userOID, _ := primitive.ObjectIDFromHex(userAuth.ID)

	candidates, err := s.Core.SuggestWithContext(userAuth.ID)
	if err != nil {
		goto returnRandomPool
	}

	return ctx.Status(fiber.StatusOK).JSON(candidates)

returnRandomPool:
	pool, err := s.Core.SuggestRandom(userOID)
	if err != nil {
		log.Println("cannot get suggest for user", userAuth.ID, "err", err)
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "cannot suggest users"})
	}
	return ctx.Status(fiber.StatusOK).JSON(pool)
}

// HandleAddUserMatch will add match-related information to match db
func (s *Service) HandleAddUserMatch(ctx *fiber.Ctx) error {
	userMatch := new(matchingdb.MatchInfo)
	if err := json.Unmarshal(ctx.Body(), userMatch); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "service: match information required in body",
		})
	}

	err := s.AddUserMatch(*userMatch)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.SendStatus(http.StatusOK)
}

const UserEmbedFormat = "[BEGIN]gender: %s[SEP]age: %v[SEP]job: %s[SEP]native language: %s[SEP]learning language: %s[SEP]country: %s[SEP]interests: %s[END]"

func (s *Service) AddUserMatch(info matchingdb.MatchInfo) error {
	requestPayload := transport.EmbeddingRequest{
		Request: transport.Request{Type: transport.Embedding},
		Payload: fmt.Sprintf(
			UserEmbedFormat,
			info.Gender,
			info.Age,
			info.Major,
			info.Native,
			strings.Join(info.Learnings, ", "),
			info.Country,
			strings.Join(info.Interests, ", "),
		),
	}
	requestBytes, _ := json.Marshal(requestPayload)
	code, body, errs := fiber.Post(s.EmbedderEndpoint).
		ContentType("application/json").
		Body(requestBytes).
		Bytes()
	if errs != nil || len(errs) > 0 {
		log.Println("cannot request embed vector", errs)
		return fmt.Errorf("failed to embed user match")
	} else if code != fiber.StatusOK {
		log.Println("cannot request embed vector, server response", code)
		return fmt.Errorf("failed to embed user match")
	}

	rsp, err := utils.ParseJSON[transport.EmbeddingResponse](body)
	if err != nil {
		log.Println("cannot get embed from embedder response", err)
		return fmt.Errorf("can not get embed from embedder response")
	}

	info, err = s.Core.AddUserMatchInformation(info)
	if err != nil {
		log.Println("cannot add user match information", err)
		return fmt.Errorf("cannot add user match information")
	}

	if err := s.Core.AddEmbedding(info.UserID, rsp.Embedded); err != nil {
		log.Println("cannot add user embed", err)
		return fmt.Errorf("cannot add user embed")
	}

	return nil
}
