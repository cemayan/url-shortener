package input

import "github.com/gofiber/fiber/v2"

type ApiUseCase interface {
	GetUserUrl(c *fiber.Ctx) error
}
