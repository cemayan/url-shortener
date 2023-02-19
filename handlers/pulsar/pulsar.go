package kafka

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cemayan/url-shortener/common"
	log "github.com/sirupsen/logrus"
)

type PulsarHandler interface {
	GetPulsarProducer(producerOptions pulsar.ProducerOptions) pulsar.Producer
	GetPulsarConsumer(consumerOptions pulsar.ConsumerOptions) pulsar.Consumer
}

type PulsarSvc struct {
	pulsarClient pulsar.Client
	pulsarConfig common.Pulsar
	log          *log.Entry
}

func (p PulsarSvc) GetPulsarProducer(producerOptions pulsar.ProducerOptions) pulsar.Producer {
	producer, err := p.pulsarClient.CreateProducer(producerOptions)

	if err != nil {
		p.log.Fatal(err)
	}

	return producer
}

func (p PulsarSvc) GetPulsarConsumer(consumerOptions pulsar.ConsumerOptions) pulsar.Consumer {
	consumer, err := p.pulsarClient.Subscribe(consumerOptions)
	if err != nil {
		p.log.Fatal(err)
	}
	return consumer
}

func NewPulsarHandler(pulsarClient pulsar.Client, pulsarConfig common.Pulsar, log *log.Entry) PulsarHandler {
	return &PulsarSvc{pulsarClient: pulsarClient, pulsarConfig: pulsarConfig, log: log}
}
