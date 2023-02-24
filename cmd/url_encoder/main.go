package main

import (
	"encoding/json"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cemayan/url-shortener/common"
	pulsar_handler "github.com/cemayan/url-shortener/common/adapters/pulsar"
	"github.com/cemayan/url-shortener/common/ports/output"
	"github.com/cemayan/url-shortener/config/url_encoder"
	"github.com/cemayan/url-shortener/internal/api/write/adapter/mongodb"
	"github.com/cemayan/url-shortener/internal/api/write/domain/model"
	"github.com/cemayan/url-shortener/managers/db"
	"github.com/cemayan/url-shortener/managers/hook"
	"github.com/cemayan/url-shortener/managers/mq"
	"github.com/cemayan/url-shortener/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"os"
	"time"
)

var _log *logrus.Logger
var configs *url_encoder.AppConfig
var v *viper.Viper

func init() {
	_log = logrus.New()
	if os.Getenv("env") == "dev" {
		_log.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
	} else {
		_log.SetFormatter(&logrus.JSONFormatter{})
		_log.SetOutput(os.Stdout)
	}
	_log.Out = os.Stdout

	v = viper.New()
	_configs := url_encoder.NewConfig(v)

	env := os.Getenv("env")
	appConfig, err := _configs.GetConfig(env)
	configs = appConfig
	if err != nil {
		return
	}

	fluentdHook := hook.NewFluentdHook()
	_log.Hooks.Add(fluentdHook)
}

func main() {

	var mongoManager db.MongodbManager
	mongoManager = db.NewMongodbManager(configs.Mongo, _log.WithFields(logrus.Fields{"service": "mongo-db-manager"}))

	var mongoPort output.MongoPort
	mongoPort = mongodb.NewApiRepo(mongoManager.New(), configs.Mongo, _log.WithFields(logrus.Fields{"service": "mongo-db-manager"}))

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               configs.Pulsar.Url,
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		_log.WithFields(logrus.Fields{"method": "NewPulsarClient", "message": err.Error()}).Log(logrus.FatalLevel)
		return
	}

	var pulsarManager mq.PulsarManager
	pulsarManager = mq.NewPulsarManager(client, configs.Pulsar, _log.WithFields(logrus.Fields{"service": "event-handler"}))

	var pulsarPort output.PulsarPort
	pulsarPort = pulsar_handler.NewPulsarHandler(pulsarManager, configs.Pulsar, _log.WithFields(logrus.Fields{"service": "event-handler"}))

	channel := make(chan pulsar.ConsumerMessage, 10)

	consumer := pulsarPort.ConsumeEvent(configs.Pulsar.Topic, fmt.Sprintf("%s-%v", "sub", rand.Int63()), channel)
	defer consumer.Close()

	for cm := range channel {
		consumer := cm.Consumer
		msg := cm.Message
		var event common.EventModel

		err := msg.GetSchemaValue(&event)
		if err != nil {
			_log.WithFields(logrus.Fields{"method": "GetSchemaValue", "message": err.Error(), "data": string(msg.Payload())}).Log(logrus.ErrorLevel)
		}

		if event.EventName == common.UrlEncoded {

			metadata := util.GetEventMetadata()

			var eventDetail common.EventDataDetail

			err := json.Unmarshal(event.EventData, &eventDetail)
			if err != nil {
				return
			}

			eventDetail.ShortUrl = "shortttt"

			err = mongoPort.CreateEvent(model.Events{
				ID: primitive.NewObjectID(),
				EventData: bson.M{
					"shortUrl": "shortttt",
					"longUrl":  eventDetail.LongUrl,
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

				err = pulsarPort.PublishEvent(common.EventModel{
					AggregateId:   metadata.AggregateId,
					AggregateType: event.AggregateType,
					EventData:     eventDetailBytes,
					EventDate:     metadata.EventDate,
					EventName:     common.UrlCreated,
					UserId:        event.UserId,
				}, configs.Pulsar.Topic)
				if err != nil {
					//TODO
					return
				}
			}

		}

		err = consumer.Ack(msg)
		if err != nil {
			_log.WithFields(logrus.Fields{"method": "ConsumerAck", "message": err.Error(), "data": fmt.Sprintf("%v||%v", msg.Topic(), msg.ID().EntryID())}).Log(logrus.ErrorLevel)
		}

	}
}
