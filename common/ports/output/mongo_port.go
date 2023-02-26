package output

import (
	"github.com/cemayan/url-shortener/common/domain"
)

type MongoPort interface {
	CreateEvent(event domain.Events) error
}
