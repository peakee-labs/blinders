package main

import (
	"log"

	"blinders/packages/db/collectingdb"
	"blinders/packages/utils"
	collecting "blinders/services/collecting/core"

	"github.com/gofiber/fiber/v2"
)

type Manager struct {
	App               *fiber.App
	CollectingService *collecting.Service
}

func NewManager(app *fiber.App, db *collectingdb.CollectingDB) *Manager {
	return &Manager{
		App:               app,
		CollectingService: collecting.NewService(db.ExplainLogsRepo, db.TranslateLogsRepo),
	}
}

func (s *Manager) InitRoute() error {
	s.App.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("service healthy")
	})

	s.App.Post("/", func(c *fiber.Ctx) error {
		event, err := utils.ParseJSON[map[string]interface{}](c.Body())
		if err != nil {
			log.Println("can not parse request:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}
		err = s.CollectingService.HandlePushEvent(*event)
		if err != nil {
			log.Println("can not handle request:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{})
	})

	s.App.Get("/", func(c *fiber.Ctx) error {
		request, err := utils.ParseJSON[map[string]interface{}](c.Body())
		if err != nil {
			log.Println("can not parse request:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}
		res, err := s.CollectingService.HandleGetRequest(*request)
		if err != nil {
			log.Println("can not handle request:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}
		return c.Status(fiber.StatusOK).JSON(res)
	})
	return nil
}
