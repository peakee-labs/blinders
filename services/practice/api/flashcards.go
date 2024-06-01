package practiceapi

import "github.com/gofiber/fiber/v2"

func (Service) HandleGetFlashcardCollectionByID(_ *fiber.Ctx) error {
	return nil
}

func (Service) HandleCreateFlashcardCollectionByID(_ *fiber.Ctx) error {
	return nil
}

func (Service) HandleUpdateFlashcardCollectionByID(_ *fiber.Ctx) error {
	return nil
}

func (Service) HandleDeleteFlashcardCollectionByID(_ *fiber.Ctx) error {
	return nil
}

// define one-time used type in the usage scope
type (
	AddFlashcardBody struct {
		FrontText string
		BackText  string
	}
	AddFlashcardResponse struct{}
)

func (Service) HandleAddFlashcardToCollection(_ *fiber.Ctx) error {
	return nil
}

func (Service) HandleUpdateFlashcardInCollection(_ *fiber.Ctx) error {
	return nil
}

func (Service) HandleRemoveFlashcardFromCollection(_ *fiber.Ctx) error {
	return nil
}
