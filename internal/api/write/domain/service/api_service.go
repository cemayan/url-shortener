package service

import (
	"encoding/json"
	"fmt"
	"github.com/cemayan/url-shortener/common"
	"github.com/cemayan/url-shortener/common/ports/output"
	"github.com/cemayan/url-shortener/config/api"
	"github.com/cemayan/url-shortener/internal/api/write/domain/model"
	"github.com/cemayan/url-shortener/internal/api/write/domain/port/input"
	mongo_output "github.com/cemayan/url-shortener/internal/api/write/domain/port/output"
	"github.com/cemayan/url-shortener/internal/api/write/util"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApiSvc struct {
	configs    *api.AppConfig
	log        *log.Entry
	pulsarPort output.PulsarPort
	mongoPort  mongo_output.MongoPort
}

func (a ApiSvc) CreateUrl(c *fiber.Ctx) error {

	userReq := new(model.Payload)
	if err := c.BodyParser(userReq); err != nil {
		a.log.WithFields(log.Fields{"method": "CreateUrl"}).Errorf("An error occurred when parsing the payload %v \n", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.Response{
			Message:    fmt.Sprintf("Review your input %s", err),
			StatusCode: 400,
		})
	}

	metadata := util.GetEventMetadata()
	eventData := bson.M{
		"LongUrl": userReq.Url,
	}

	err := a.mongoPort.CreateEvent(model.Events{
		ID:            primitive.NewObjectID(),
		EventData:     eventData,
		EventDate:     metadata.EventDate,
		EventName:     metadata.EventName,
		AggregateType: "",
		AggregateId:   metadata.AggregateId,
		UserId:        "user1",
	})
	if err != nil {
		return c.Status(400).JSON(`test`)
	} else {

		userReqBytes, err := json.Marshal(userReq)
		if err != nil {
			return err
		}

		encodedUrlEvent := common.EventModel{
			EventData:     userReqBytes,
			EventDate:     metadata.EventDate,
			EventName:     metadata.EventName,
			AggregateType: "test",
			AggregateId:   metadata.AggregateId,
			UserId:        "user1",
		}

		err = a.pulsarPort.PublishEvent(encodedUrlEvent, "events")
		if err != nil {
			return err
		}
	}
	return c.JSON(`test`)
}

func NewApiService(mongoPort mongo_output.MongoPort, pulsarPort output.PulsarPort, configs *api.AppConfig, log *log.Entry) input.ApiUseCase {
	return &ApiSvc{
		configs:    configs,
		log:        log,
		mongoPort:  mongoPort,
		pulsarPort: pulsarPort,
	}
}
