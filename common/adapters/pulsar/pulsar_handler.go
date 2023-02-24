package pulsar

import (
	"context"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cemayan/url-shortener/common"
	"github.com/cemayan/url-shortener/common/ports/output"
	"github.com/cemayan/url-shortener/managers/mq"
	log "github.com/sirupsen/logrus"
)

type PulsarSvc struct {
	pulsarManager mq.PulsarManager
	configs       common.Pulsar
	log           *log.Entry
}

var (
	avroSchemaDef = "{\"type\":\"record\",\"name\":\"Example\",\"namespace\":\"public\",\"fields\":[{\"name\":\"aggregate_id\",\"type\":\"string\"},{\"name\":\"aggregate_type\",\"type\":\"string\"},{\"name\":\"event_data\",\"type\":\"string\"},{\"name\":\"event_date\",\"type\":\"long\"},{\"name\":\"event_name\",\"type\":\"string\"},{\"name\":\"user_id\",\"type\":\"string\"}]}"
)

func (p PulsarSvc) ConsumeEvent(topicName string, subscriptionName string, channel chan pulsar.ConsumerMessage) pulsar.Consumer {
	consumer := p.pulsarManager.GetPulsarConsumer(pulsar.ConsumerOptions{
		Topic:            topicName,
		SubscriptionName: subscriptionName,
		Schema:           pulsar.NewAvroSchema(avroSchemaDef, nil),
		MessageChannel:   channel,
	})

	return consumer
}

func (p PulsarSvc) PublishEvent(event common.EventModel, topicName string) error {
	producer := p.pulsarManager.GetPulsarProducer(pulsar.ProducerOptions{
		Topic:  topicName,
		Schema: pulsar.NewAvroSchema(avroSchemaDef, nil),
	})

	producer.SendAsync(context.Background(), &pulsar.ProducerMessage{
		Value: event,
	}, func(id pulsar.MessageID, message *pulsar.ProducerMessage, err error) {

	})

	return nil
}

func NewPulsarHandler(pulsarManager mq.PulsarManager, configs common.Pulsar, log *log.Entry) output.PulsarPort {
	return &PulsarSvc{pulsarManager: pulsarManager, configs: configs, log: log}
}
