package router

import (
	"github.com/cemayan/url-shortener/config/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
)

// SetupRoutes creates the fiber's routes
// api/v1 is root group.
// Before the reach services interface is configured
func SetupRoutes(app *fiber.App, configs *api.AppConfig, _log *log.Entry) {

	apiGroup := app.Group("/api", logger.New())
	v1Group := apiGroup.Group("/v1")

	usGroup := v1Group.Group("/url")
	usGroup.Post("/", nil)
	usGroup.Put("/:id", nil)
	usGroup.Delete("/", nil)
}
