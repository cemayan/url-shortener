package db

import (
	"fmt"
	"github.com/cemayan/url-shortener/common"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

type CockroackDBManager interface {
	New() *gorm.DB
}

type DBService struct {
	configs *common.Cockroach
	_log    *log.Entry
}

// New  serves to connect to db
// When DB connection is successful then model migration is started
func (d DBService) New() *gorm.DB {

	newLogger := logger.New(
		d._log.Logger,
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	getenv := os.Getenv("env")
	var DSN string

	if getenv == "dev" {
		DSN = fmt.Sprintf("host=%s port=%s  user=%s  dbname=%s sslmode=disable",
			d.configs.Host,
			d.configs.Port,
			d.configs.User,
			d.configs.Name)
	} else {
		DSN = fmt.Sprintf("host=%s port=%s  user=%s password=%s  dbname=%s",
			d.configs.Host,
			d.configs.Port,
			d.configs.User,
			d.configs.Password,
			d.configs.Name)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: DSN,
	}), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic("failed to connect database")
	}

	d._log.WithFields(log.Fields{"service": "database"}).Infoln("Connection Opened to Database")

	return db
}

func NewCockroachDbManager(configs *common.Cockroach, _log *log.Entry) CockroackDBManager {
	return &DBService{configs: configs, _log: _log}
}
