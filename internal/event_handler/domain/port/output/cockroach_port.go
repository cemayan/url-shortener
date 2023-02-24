package output

import "github.com/cemayan/url-shortener/internal/event_handler/domain/model"

type CockroachPort interface {
	CreateUserUrl(user *model.UserUrl) (*model.UserUrl, error)
}
