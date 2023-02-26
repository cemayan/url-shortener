package service

import (
	"encoding/json"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cemayan/url-shortener/common"
	"github.com/cemayan/url-shortener/common/domain"
	"github.com/cemayan/url-shortener/common/ports/output"
	"github.com/cemayan/url-shortener/config/url_encoder"
	"github.com/cemayan/url-shortener/util"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
)

const (
	urlPrefix = "https://tiny.url/"
)

type UrlEncodeSvc struct {
	mongoPort  output.MongoPort
	pulsarPort output.PulsarPort
	configs    *url_encoder.AppConfig
	log        *log.Entry
}

func (u *UrlEncodeSvc) RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func (u *UrlEncodeSvc) NewUrlString() string {
	return u.RandomString(6)
}

func (u *UrlEncodeSvc) Consume() {
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

		if event.EventName == common.UrlEncoded {

			metadata := util.GetEventMetadata()

			var eventDetail common.EventDataDetail

			err := json.Unmarshal(event.EventData, &eventDetail)
			if err != nil {
				return
			}

			newUrlString := u.NewUrlString()

			eventDetail.ShortUrl = urlPrefix + newUrlString
			eventDetail.UrlString = newUrlString

			err = u.mongoPort.CreateEvent(domain.Events{
				ID: primitive.NewObjectID(),
				EventData: bson.M{
					"shortUrl": eventDetail.ShortUrl,
					"longUrl":  eventDetail.LongUrl,
					"urlStr":   newUrlString,
				},
				EventDate:     metadata.EventDate,
				EventName:     common.UrlCreated,
				AggregateType: "",
				AggregateId:   metadata.AggregateId,
				UserId:        event.UserId,
			})
			if err != nil {
				//TODO
				return
			} else {

				eventDetailBytes, err := json.Marshal(eventDetail)
				if err != nil {
					return
				}

				err = u.pulsarPort.PublishEvent(common.EventModel{
					AggregateId:   metadata.AggregateId,
					AggregateType: event.AggregateType,
					EventData:     eventDetailBytes,
					EventDate:     metadata.EventDate,
					EventName:     common.UrlCreated,
					UserId:        event.UserId,
				}, u.configs.Pulsar.Topic)
				if err != nil {
					//TODO
					return
				}
			}

		}

		err = consumer.Ack(msg)
		if err != nil {
			u.log.WithFields(logrus.Fields{"method": "Consume", "message": err.Error(), "data": fmt.Sprintf("%v||%v", msg.Topic(), msg.ID().EntryID())}).Log(logrus.ErrorLevel)
		}

	}
}

func NewUrlEncodeService(mongoPort output.MongoPort, pulsarPort output.PulsarPort, configs *url_encoder.AppConfig, log *log.Entry) *UrlEncodeSvc {
	return &UrlEncodeSvc{
		mongoPort:  mongoPort,
		pulsarPort: pulsarPort,
		configs:    configs,
		log:        log,
	}
}
