package main

import (
	"github.com/apache/pulsar-client-go/pulsar"
	pulsar_handler "github.com/cemayan/url-shortener/common/adapters/pulsar"
	"github.com/cemayan/url-shortener/common/adapters/redis"
	"github.com/cemayan/url-shortener/common/ports/output"
	"github.com/cemayan/url-shortener/config/event_handler"
	"github.com/cemayan/url-shortener/internal/event_handler/adapter/database"
	"github.com/cemayan/url-shortener/internal/event_handler/domain/service"
	"github.com/cemayan/url-shortener/internal/event_handler/helper"
	"github.com/cemayan/url-shortener/managers/cache"
	"github.com/cemayan/url-shortener/managers/db"
	"github.com/cemayan/url-shortener/managers/hook"
	"github.com/cemayan/url-shortener/managers/mq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

	helper.MigrateDB(_db, _log.WithFields(logrus.Fields{"service": "database"}))

	fluentdHook := hook.NewFluentdHook()
	_log.Hooks.Add(fluentdHook)
}

func main() {

	pulsarClient, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               configs.Pulsar.Url,
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})

	if err != nil {
		_log.WithFields(logrus.Fields{"method": "NewPulsarClient", "message": err.Error()}).Log(logrus.FatalLevel)

	}

	var pulsarManager mq.PulsarManager
	pulsarManager = mq.NewPulsarManager(pulsarClient, configs.Pulsar, _log.WithFields(logrus.Fields{"service": "event-handler"}))

	var pulsarPort output.PulsarPort
	pulsarPort = pulsar_handler.NewPulsarHandler(pulsarManager, configs.Pulsar, _log.WithFields(logrus.Fields{"service": "event-handler"}))

	var redisManager cache.RedisManager
	redisManager = cache.NewRedisManager(configs.Redis, _log.WithFields(logrus.Fields{"service": "cache-manager"}))

	var redisPort output.RedisPort
	redisPort = redis.NewRedisHandler(redisManager.New(), _log.WithFields(logrus.Fields{"service": "cache-service"}))

	eventService := service.NewEventService(redisPort, pulsarPort, configs, _log.WithFields(logrus.Fields{"service": "event-service"}))
	eventService.Consume()
}
