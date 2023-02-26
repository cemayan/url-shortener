package middleware

import (
	"github.com/cemayan/url-shortener/common"
	"github.com/cemayan/url-shortener/config/user"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// Protected protect routes
func Protected(configs *user.AppConfig) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(configs.Secret),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(400).JSON(common.Response{
			Message: "Missing or malformed JWT",
		})

	}
	return c.Status(401).JSON(common.Response{
		Message: "Invalid or expired JWT",
	})
}
