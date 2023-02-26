package output

import "github.com/cemayan/url-shortener/internal/api/read/domain/model"

type CockroachPort interface {
	GetUserUrl(urlStr string) (model.UserUrl, error)
}
