package router

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cemayan/url-shortener/common/adapters/mongodb"
	pulsar_handler "github.com/cemayan/url-shortener/common/adapters/pulsar"
	"github.com/cemayan/url-shortener/common/ports/output"
	"github.com/cemayan/url-shortener/config/api"
	"github.com/cemayan/url-shortener/internal/api/write/domain/port/input"
	"github.com/cemayan/url-shortener/internal/api/write/domain/service"
	"github.com/cemayan/url-shortener/managers/db"
	"github.com/cemayan/url-shortener/managers/mq"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
	"time"
)

// SetupRoutes creates the fiber's routes
// api/v1 is root group.
// Before the reach services interface is configured
func SetupRoutes(app *fiber.App, configs *api.AppConfig, _log *log.Entry) {

	apiGroup := app.Group("/api", logger.New())
	v1Group := apiGroup.Group("/v1")

	var mongoManager db.MongodbManager
	mongoManager = db.NewMongodbManager(configs.Mongo, _log)

	var mongoPort output.MongoPort
	mongoPort = mongodb.NewApiRepo(mongoManager.New(), configs.Mongo, _log)

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               configs.Pulsar.Url,
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}

	var pulsarManager mq.PulsarManager
	pulsarManager = mq.NewPulsarManager(client, configs.Pulsar, _log)

	var pulsarPort output.PulsarPort
	pulsarPort = pulsar_handler.NewPulsarHandler(pulsarManager, configs.Pulsar, _log)

	var apiUseCase input.ApiUseCase
	apiUseCase = service.NewApiService(mongoPort, pulsarPort, configs, _log)

	usGroup := v1Group.Group("/url")
	usGroup.Post("/", apiUseCase.CreateUrl)
	usGroup.Put("/:id", nil)
	usGroup.Delete("/", nil)
}
