package router

import (
	"github.com/cemayan/url-shortener/common/adapters/redis"
	common_port "github.com/cemayan/url-shortener/common/ports/output"
	"github.com/cemayan/url-shortener/config/api"
	"github.com/cemayan/url-shortener/internal/api/read/adapter/cockroach"
	"github.com/cemayan/url-shortener/internal/api/read/adapter/database"
	"github.com/cemayan/url-shortener/internal/api/read/domain/port/input"
	"github.com/cemayan/url-shortener/internal/api/read/domain/port/output"
	"github.com/cemayan/url-shortener/internal/api/read/domain/service"
	"github.com/cemayan/url-shortener/managers/cache"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// SetupRoutes creates the fiber's routes
// api/v1 is root group.
// Before the reach services interface is configured
func SetupRoutes(app *fiber.App, configs *api.AppConfig, _log *log.Entry) {

	apiGroup := app.Group("/api", logger.New())
	v1Group := apiGroup.Group("/v1")

	var redisManager cache.RedisManager
	redisManager = cache.NewRedisManager(configs.Redis, _log.WithFields(logrus.Fields{"service": "cache-manager"}))

	var redisPort common_port.RedisPort
	redisPort = redis.NewRedisHandler(redisManager.New(), _log.WithFields(logrus.Fields{"service": "cache-service"}))

	var cockroachPort output.CockroachPort
	cockroachPort = cockroach.NewUserUrlRepo(database.DB, _log.WithFields(logrus.Fields{"service": "cockroach-repo"}))

	var apiUseCase input.ApiUseCase
	apiUseCase = service.NewApiService(redisPort, cockroachPort, configs, _log)

	usGroup := v1Group.Group("/url")
	usGroup.Get("/:id", apiUseCase.GetUserUrl)
}
