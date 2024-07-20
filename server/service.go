package main

import "github.com/gofiber/fiber/v2"

type Service struct {
	PathPrefix string
	Fiber      FiberService
}

type FiberService interface {
	InitFiberRoutes(fiber.Router)
}
