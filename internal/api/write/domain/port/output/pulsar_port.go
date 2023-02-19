package output

import "github.com/cemayan/url-shortener/internal/api/write/domain/model/event"

type PulsarPort interface {
	PublishEvent(event event.Event, topicName string) error
}
