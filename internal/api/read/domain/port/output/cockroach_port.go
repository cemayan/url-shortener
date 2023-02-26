package output

import "github.com/cemayan/url-shortener/internal/api/read/domain/model"

type CockroachPort interface {
	GetUserUrl(longUrl string) (model.UserUrl, error)
}
