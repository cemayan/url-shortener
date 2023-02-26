package main

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cemayan/url-shortener/common/adapters/mongodb"
	pulsar_handler "github.com/cemayan/url-shortener/common/adapters/pulsar"
	"github.com/cemayan/url-shortener/common/ports/output"
	"github.com/cemayan/url-shortener/config/url_encoder"
	"github.com/cemayan/url-shortener/internal/url_encoder/domain/service"
	"github.com/cemayan/url-shortener/managers/db"
	"github.com/cemayan/url-shortener/managers/hook"
	"github.com/cemayan/url-shortener/managers/mq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

	encodeService := service.NewUrlEncodeService(mongoPort, pulsarPort, configs, _log.WithFields(logrus.Fields{"service": "url-service"}))

	encodeService.Consume()
}
