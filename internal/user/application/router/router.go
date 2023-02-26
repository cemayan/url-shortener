package router

import (
	"github.com/cemayan/url-shortener/config/user"
	"github.com/cemayan/url-shortener/internal/user/application/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
)

// SetupRoutes creates the fiber's routes
// api/v1 is root group.
// Before the reach services interface is configured
func SetupRoutes(app *fiber.App, log *log.Entry, configs *user.AppConfig) {

	api := app.Group("/api", logger.New())
	v1 := api.Group("/v1")

	v1.Get("/health", nil)

	userGroup := v1.Group("/user")
	userGroup.Get("/:id?", middleware.Protected(configs), nil)
	userGroup.Post("/", nil)
	userGroup.Put("/:id", middleware.Protected(configs), nil)
	userGroup.Delete("/:id", middleware.Protected(configs), nil)

}
