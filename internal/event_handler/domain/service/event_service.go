package service

import (
	"encoding/json"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cemayan/url-shortener/common"
	"github.com/cemayan/url-shortener/common/ports/output"
	"github.com/cemayan/url-shortener/config/event_handler"
	"github.com/cemayan/url-shortener/internal/event_handler/adapter/cockroach"
	"github.com/cemayan/url-shortener/internal/event_handler/adapter/database"
	"github.com/cemayan/url-shortener/internal/event_handler/domain/model"
	cockroach_output "github.com/cemayan/url-shortener/internal/event_handler/domain/port/output"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"math/rand"
)

type EventSvc struct {
	redisPort  output.RedisPort
	pulsarPort output.PulsarPort
	configs    *event_handler.AppConfig
	log        *log.Entry
}

func (u *EventSvc) Consume() {
	channel := make(chan pulsar.ConsumerMessage, 10)

	consumer := u.pulsarPort.ConsumeEvent(u.configs.Pulsar.Topic, fmt.Sprintf("%s-%v", "sub", rand.Int63()), channel)
	defer consumer.Close()

	for cm := range channel {
		consumer := cm.Consumer
		msg := cm.Message
		var event common.EventModel

		err := msg.GetSchemaValue(&event)
		if err != nil {
			u.log.WithFields(logrus.Fields{"method": "GetSchemaValue", "message": err.Error(), "data": string(msg.Payload())}).Log(logrus.ErrorLevel)
		}

		if event.EventName == common.UrlCreated {

			var cockroachPort cockroach_output.CockroachPort
			cockroachPort = cockroach.NewUserUrlRepo(database.DB, u.log.WithFields(logrus.Fields{"service": "cockroach-repo"}))

			var eventDetail common.EventDataDetail

			err := json.Unmarshal(event.EventData, &eventDetail)
			if err != nil {
				return
			}

			_, err = cockroachPort.CreateUserUrl(&model.UserUrl{
				UserId:    event.UserId,
				ShortUrl:  eventDetail.ShortUrl,
				LongUrl:   eventDetail.LongUrl,
				UrlString: eventDetail.UrlString,
			})
			if err != nil {
				u.log.WithFields(logrus.Fields{"method": "CockroachCreateUserUrl", "message": err.Error()}).Log(logrus.ErrorLevel)
				return
			}

			err = u.redisPort.Set(eventDetail.ShortUrl, eventDetail.LongUrl)
			if err != nil {
				u.log.WithFields(logrus.Fields{"method": "RedisSet", "message": err.Error()}).Log(logrus.ErrorLevel)
				return
			}

		}

		err = consumer.Ack(msg)
		if err != nil {
			u.log.WithFields(logrus.Fields{"method": "ConsumerAck", "message": err.Error(), "data": fmt.Sprintf("%v||%v", msg.Topic(), msg.ID().EntryID())}).Log(logrus.ErrorLevel)
		}
	}
}

func NewEventService(redisPort output.RedisPort, pulsarPort output.PulsarPort, configs *event_handler.AppConfig, log *log.Entry) *EventSvc {
	return &EventSvc{
		redisPort:  redisPort,
		pulsarPort: pulsarPort,
		configs:    configs,
		log:        log,
	}
}
