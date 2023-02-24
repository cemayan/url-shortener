package output

import "github.com/cemayan/url-shortener/internal/api/write/domain/model"

type MongoPort interface {
	CreateEvent(event model.Events) error
}
