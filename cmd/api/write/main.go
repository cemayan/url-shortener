package main

import (
	"github.com/cemayan/url-shortener/config/api"
	"github.com/cemayan/url-shortener/handlers/db"
	"github.com/cemayan/url-shortener/internal/api/read/adapter/database"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

var _log *logrus.Logger
var configs *api.AppConfig
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
	_configs := api.NewConfig(v)

	env := os.Getenv("env")
	appConfig, err := _configs.GetConfig(env)
	configs = appConfig
	if err != nil {
		return
	}

	dbHandler := db.NewCockroachDbHandler(&configs.Cockroach, _log.WithFields(logrus.Fields{"service": "database"}))
	_db := dbHandler.New()
	database.DB = _db
	//util.MigrateDB(_db, _log.WithFields(logrus.Fields{"service": "database"}))
}

func main() {

}
