package cockroach

import (
	"github.com/cemayan/url-shortener/internal/api/read/domain/model"
	"github.com/cemayan/url-shortener/internal/api/read/domain/port/output"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserUrlrepo struct {
	db  *gorm.DB
	log *log.Entry
}

func (u UserUrlrepo) GetUserUrl(urlStr string) (model.UserUrl, error) {
	var url model.UserUrl
	tx := u.db.Where("url_string = ? ", urlStr).First(&url)
	return url, tx.Error
}

func NewUserUrlRepo(db *gorm.DB, log *log.Entry) output.CockroachPort {
	return &UserUrlrepo{
		db:  db,
		log: log,
	}
}
