package main

import (
	"github.com/cemayan/url-shortener/config/event_handler"
	"github.com/cemayan/url-shortener/internal/event_handler/adapter/database"
	"github.com/cemayan/url-shortener/managers/db"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
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
	//util.MigrateDB(_db, _log.WithFields(logrus.Fields{"service": "database"}))
}

func main() {

}
