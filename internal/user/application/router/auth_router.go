package router

import (
	"github.com/cemayan/url-shortener/config/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func prometheusHandler2() http.Handler {
	return promhttp.Handler()
}

// SetupAuthRoutes creates the fiber's routes
// api/v1 is root group.
// Before the reach services interface is configured
func SetupAuthRoutes(app *fiber.App, log *log.Entry, configs *user.AppConfig) {

	api := app.Group("/api", logger.New())
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	auth.Post("/getToken", nil)
	auth.Post("/validateToken", nil)

}
