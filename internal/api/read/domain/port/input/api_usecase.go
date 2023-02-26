package input

import "github.com/gofiber/fiber/v2"

type ApiUseCase interface {
	Forward(c *fiber.Ctx) error
}
