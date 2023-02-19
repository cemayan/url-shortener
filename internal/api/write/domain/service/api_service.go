package service

import (
	"github.com/cemayan/url-shortener/config/api"
	"github.com/cemayan/url-shortener/internal/api/write/domain/model/event"
	"github.com/cemayan/url-shortener/internal/api/write/domain/port/input"
	"github.com/cemayan/url-shortener/internal/api/write/domain/port/output"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type ApiSvc struct {
	configs    *api.AppConfig
	log        *log.Entry
	pulsarPort output.PulsarPort
	mongoPort  output.MongoPort
}

func (a ApiSvc) CreateUrl(c *fiber.Ctx) error {

	createEvent := event.NewUrlCreate(nil)

	err := a.pulsarPort.PublishEvent(createEvent, "events")
	if err != nil {
		return err
	}

	//err := a.mongoPort.CreateUserUrl(model.UserUrl{
	//	UserId: "test",
	//})
	//if err != nil {
	//	return c.Status(400).JSON(`test`)
	//}
	return c.JSON(`test`)
}

func NewApiService(mongoPort output.MongoPort, pulsarPort output.PulsarPort, configs *api.AppConfig, log *log.Entry) input.ApiUseCase {
	return &ApiSvc{
		configs:    configs,
		log:        log,
		mongoPort:  mongoPort,
		pulsarPort: pulsarPort,
	}
}
