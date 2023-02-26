package service

import (
	"fmt"
	"github.com/cemayan/url-shortener/common"
	common_port "github.com/cemayan/url-shortener/common/ports/output"
	"github.com/cemayan/url-shortener/config/api"
	"github.com/cemayan/url-shortener/internal/api/read/domain/port/input"
	"github.com/cemayan/url-shortener/internal/api/read/domain/port/output"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type ApiSvc struct {
	configs       *api.AppConfig
	log           *log.Entry
	redisPort     common_port.RedisPort
	cockroachPort output.CockroachPort
}

func (a ApiSvc) GetUserUrl(c *fiber.Ctx) error {

	url, err := a.cockroachPort.GetUserUrl("")
	if err != nil {
		return c.Status(400).JSON(common.Response{
			StatusCode: 400,
			Message:    fmt.Sprintf("Url not found or request is bad %v", err),
		})
	}

	return c.Status(200).JSON(common.Response{
		Data:       url,
		StatusCode: 200,
	})
}

func NewApiService(redisPort common_port.RedisPort, cockroachPort output.CockroachPort, configs *api.AppConfig, log *log.Entry) input.ApiUseCase {
	return &ApiSvc{
		configs:       configs,
		log:           log,
		redisPort:     redisPort,
		cockroachPort: cockroachPort,
	}
}
