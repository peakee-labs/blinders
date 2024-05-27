package exploreapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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
	Transport        transport.Transport
}

func NewService(
	exploreCore explore.Explorer,
	redisClient *redis.Client,
	transport transport.Transport,
) *Service {
	return &Service{
		Core:        exploreCore,
		RedisClient: redisClient,
		Transport:   transport,
	}
}

func (s *Service) HandleGetMatchingProfile(ctx *fiber.Ctx) error {
	matchID := ctx.Params("id")
	if matchID == "" {
		log.Println("cannot match id is empty")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "match id is empty"})
	}
	matchOID, err := primitive.ObjectIDFromHex(matchID)
	if err != nil {
		log.Println("cannot convert match id to object id", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid match id"})
	}

	profile, err := s.Core.GetMatchingProfile(matchOID)
	if err != nil {
		log.Println("cannot get matching profile", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get matching profile"})
	}
	return ctx.Status(fiber.StatusOK).JSON(profile)
}

// HandleGetMatches returns 5 users that similarity with current user.
func (s *Service) HandleGetMatches(ctx *fiber.Ctx) error {
	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok || userAuth == nil {
		log.Println("cannot get auth user")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot get user"})
	}
	userOID, _ := primitive.ObjectIDFromHex(userAuth.ID)

	candidates, err := s.Core.SuggestWithContext(userOID)
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

type matchUserBody struct {
	Gender    string   `json:"gender"`
	Major     string   `json:"major"`
	Native    string   `json:"native"`    // language code with RFC-5646 format
	Country   string   `json:"country"`   // ISO-3166 format
	Learnings []string `json:"learnings"` // languages code with RFC-5646 format
	Interests []string `json:"interests"`
	Age       int      `json:"age"`
}

func (s *Service) HandleAddMatchingProfile(ctx *fiber.Ctx) error {
	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panic("cannot get auth user")
	}
	userOID, _ := primitive.ObjectIDFromHex(userAuth.ID)
	userMatch := new(matchUserBody)

	if err := json.Unmarshal(ctx.Body(), userMatch); err != nil {
		log.Println("cannot unmarshal user match", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot unmarshal user match"})
	}
	matchInformation, err := utils.JSONConvert[matchingdb.MatchInfo](userMatch)
	if err != nil {
		log.Println("cannot convert match info", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot convert match info"})
	}
	matchInformation.UserID = userOID

	err = s.AddUserMatch(matchInformation)
	if err != nil {
		log.Println("cannot add user match", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot add user match"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (s Service) HandleUpdateMatchingProfile(ctx *fiber.Ctx) error {
	userAuth, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panic("cannot get auth user")
	}
	userOID, _ := primitive.ObjectIDFromHex(userAuth.ID)
	currentInformation, err := s.Core.GetMatchingProfile(userOID)
	if err != nil {
		log.Println("cannot get current matching profile", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "matching profile must be created first"})
	}
	userMatch := new(matchUserBody)
	if err := json.Unmarshal(ctx.Body(), userMatch); err != nil {
		log.Println("cannot unmarshal user match", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot unmarshal user match"})
	}

	matchInformation, err := utils.JSONConvert[matchingdb.MatchInfo](userMatch)
	if err != nil {
		log.Println("cannot convert match info", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot convert match info"})
	}
	matchInformation.UserID = userOID
	matchInformation.SetID(currentInformation.ID)
	matchInformation.SetInitTime(currentInformation.CreatedAt.Time())

	embed, err := s.HandleGetEmbedding(matchInformation)
	if err != nil {
		log.Println("cannot get embedding", err)
		return fmt.Errorf("cannot get embedding")
	}

	_, err = s.Core.UpdaterUserMatchInformation(matchInformation)
	if err != nil {
		log.Println("cannot update user match", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot update user match"})
	}

	if err := s.Core.UpdateEmbedding(userOID, embed); err != nil {
		log.Println("cannot update user embed", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot update user embed"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// InternalHandleAddUserMatch will add match-related information to match db
func (s *Service) InternalHandleAddUserMatch(ctx *fiber.Ctx) error {
	userMatch := new(matchingdb.MatchInfo)
	if err := json.Unmarshal(ctx.Body(), userMatch); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "service: match information required in body",
		})
	}

	err := s.AddUserMatch(userMatch)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.SendStatus(http.StatusOK)
}

const UserEmbedFormat = "[BEGIN]gender: %s[SEP]age: %v[SEP]job: %s[SEP]native language: %s[SEP]learning language: %s[SEP]country: %s[SEP]interests: %s[END]"

func (s *Service) HandleGetEmbedding(info *matchingdb.MatchInfo) ([]float32, error) {
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	requestBytes, _ := json.Marshal(requestPayload)
	response, err := s.Transport.Request(ctx, s.Transport.ConsumerID(transport.Embed), requestBytes)
	if err != nil {
		log.Println("cannot request embed vector", err)
		return nil, fmt.Errorf("failed to embed user match")
	}

	rsp, err := utils.ParseJSON[transport.EmbeddingResponse](response)
	if err != nil {
		log.Println("cannot get embed from embedder response", err)
		return nil, fmt.Errorf("can not get embed from embedder response")
	}

	return rsp.Embedded, nil
}

func (s *Service) AddUserMatch(info *matchingdb.MatchInfo) error {
	embed, err := s.HandleGetEmbedding(info)
	if err != nil {
		log.Println("cannot get embedding", err)
		return fmt.Errorf("cannot get embedding")
	}

	info, err = s.Core.AddUserMatchInformation(info)
	if err != nil {
		log.Println("cannot add user match information", err)
		return fmt.Errorf("cannot add user match information")
	}

	if err := s.Core.AddEmbedding(info.UserID, embed); err != nil {
		log.Println("cannot add user embed", err)
		return fmt.Errorf("cannot add user embed")
	}

	return nil
}
