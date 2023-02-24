package cockroach

import (
	"github.com/cemayan/url-shortener/internal/event_handler/domain/model"
	"github.com/cemayan/url-shortener/internal/event_handler/domain/port/output"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserUrlrepo struct {
	db  *gorm.DB
	log *log.Entry
}

func (r UserUrlrepo) CreateUserUrl(user *model.UserUrl) (*model.UserUrl, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func NewUserUrlRepo(db *gorm.DB, log *log.Entry) output.CockroachPort {
	return &UserUrlrepo{
		db:  db,
		log: log,
	}
}
