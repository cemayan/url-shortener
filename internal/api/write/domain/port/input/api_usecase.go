package input

import "github.com/gofiber/fiber/v2"

type ApiUseCase interface {
	CreateUrl(c *fiber.Ctx) error
}
