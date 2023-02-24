package helper

import (
	"github.com/cemayan/url-shortener/internal/event_handler/domain/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
)

func MigrateDB(db *gorm.DB, _log *log.Entry) {
	if os.Getenv("env") == "test" {
		// ConnectDBForTesting  serves to connect to db for Testing
		// When DB connection is successful then model migration is started
		db.Migrator().DropTable(&model.UserUrl{})
		err := db.AutoMigrate(&model.UserUrl{})
		if err != nil {
			_log.WithFields(log.Fields{"method": "MigrateDB"}).Errorf("An error occured when migrating the db %v", err)
			return
		}

		_log.WithFields(log.Fields{"method": "MigrateDB"}).Println("Database migrated")
	} else {
		err := db.AutoMigrate(&model.UserUrl{})
		if err != nil {
			_log.WithFields(log.Fields{"method": "MigrateDB"}).Errorf("An error occured when migrating the db %v", err)
			return
		}

		_log.WithFields(log.Fields{"method": "MigrateDB"}).Println("Database migrated")
	}
}
