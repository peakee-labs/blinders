package practiceapi

import "github.com/gofiber/fiber/v2"

// TODO: if you want to check if a collection is own by a user or not, use this middleware after main handler instead
func (Service) CheckFlashcardCollectionOwnership(_ *fiber.Ctx) error {
	return nil
}
