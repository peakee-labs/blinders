package api

import (
	"log"

	"blinders/packages/transport"
	"blinders/services/embedder/core"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	// This service is used for local development only
	App      *fiber.App
	Embedder *core.Embedder
}

func NewService(app *fiber.App, embedder *core.Embedder) *Service {
	return &Service{
		App:      app,
		Embedder: embedder,
	}
}

func (s *Service) InitRoute() {
	embedderRoute := s.App.Group("/embedd")
	embedderRoute.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("hello from embedder service")
	})

	embedderRoute.Get("/", s.HandleEmbedding)
}

func (s Service) HandleEmbedding(ctx *fiber.Ctx) error {
	req := new(transport.EmbeddingRequest)
	if err := ctx.BodyParser(req); err != nil {
		log.Println("cannot parse request", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse request",
		})
	}

	embeddVector, err := s.Embedder.Embedding(req.Payload)
	if err != nil {
		log.Println("cannot get embedding", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot get embedding",
		})
	}

	response := transport.EmbeddingResponse{
		Embedded: embeddVector,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
