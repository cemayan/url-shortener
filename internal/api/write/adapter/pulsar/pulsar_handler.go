package pulsar

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cemayan/url-shortener/config/api"
	"github.com/cemayan/url-shortener/internal/api/write/domain/model/event"
	"github.com/cemayan/url-shortener/internal/api/write/domain/port/output"
	"github.com/cemayan/url-shortener/managers/mq"
	log "github.com/sirupsen/logrus"
)

type PulsarSvc struct {
	pulsarManager mq.PulsarManager
	configs       *api.AppConfig
	log           *log.Entry
}

func (p PulsarSvc) PublishEvent(event event.Event, topicName string) error {

	producer := p.pulsarManager.GetPulsarProducer(pulsar.ProducerOptions{
		Topic:  topicName,
		Schema: pulsar.NewBytesSchema(nil),
	})

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return errors.New("")
	}

	producer.SendAsync(context.Background(), &pulsar.ProducerMessage{
		Value: eventBytes,
	}, func(id pulsar.MessageID, message *pulsar.ProducerMessage, err error) {

	})

	return nil
}

func NewPulsarHandler(pulsarManager mq.PulsarManager, configs *api.AppConfig, log *log.Entry) output.PulsarPort {
	return &PulsarSvc{pulsarManager: pulsarManager, configs: configs, log: log}
}
