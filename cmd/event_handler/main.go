package main

import (
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cemayan/url-shortener/common"
	pulsar_handler "github.com/cemayan/url-shortener/common/adapters/pulsar"
	"github.com/cemayan/url-shortener/common/ports/output"
	"github.com/cemayan/url-shortener/config/event_handler"
	"github.com/cemayan/url-shortener/internal/event_handler/adapter/cockroach"
	"github.com/cemayan/url-shortener/internal/event_handler/adapter/database"
	"github.com/cemayan/url-shortener/internal/event_handler/domain/model"
	cockroack_output "github.com/cemayan/url-shortener/internal/event_handler/domain/port/output"
	"github.com/cemayan/url-shortener/internal/event_handler/util"
	"github.com/cemayan/url-shortener/managers/db"
	"github.com/cemayan/url-shortener/managers/hook"
	"github.com/cemayan/url-shortener/managers/mq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"time"
)

var _log *logrus.Logger
var configs *event_handler.AppConfig
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
	_configs := event_handler.NewConfig(v)

	env := os.Getenv("env")
	appConfig, err := _configs.GetConfig(env)
	configs = appConfig
	if err != nil {
		return
	}

	dbHandler := db.NewCockroachDbManager(&configs.Cockroach, _log.WithFields(logrus.Fields{"service": "database"}))
	_db := dbHandler.New()
	database.DB = _db

	util.MigrateDB(_db, _log.WithFields(logrus.Fields{"service": "database"}))

	fluentdHook := hook.NewFluentdHook()
	_log.Hooks.Add(fluentdHook)
}

func main() {

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

		_log.Println(event)

		err = consumer.Ack(msg)
		if err != nil {
			_log.WithFields(logrus.Fields{"method": "ConsumerAck", "message": err.Error(), "data": fmt.Sprintf("%v||%v", msg.Topic(), msg.ID().EntryID())}).Log(logrus.ErrorLevel)
		}

		var cockroachPort cockroack_output.CockroachPort
		cockroachPort = cockroach.NewUserUrlRepo(database.DB, _log.WithFields(logrus.Fields{"service": "cockroach-repo"}))

		_, err = cockroachPort.CreateUserUrl(&model.UserUrl{
			UserId:   "",
			ShortUrl: "short",
			LongUrl:  "long",
		})
		if err != nil {
			return
		}

	}
}
