package mq

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cemayan/url-shortener/common"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"time"
)

var producerInstance *pulsar.Producer

type PulsarManager interface {
	GetPulsarProducer(producerOptions pulsar.ProducerOptions) pulsar.Producer
	GetPulsarConsumer(consumerOptions pulsar.ConsumerOptions) pulsar.Consumer
	New() pulsar.Client
}

type PulsarSvc struct {
	pulsarClient pulsar.Client
	pulsarConfig common.Pulsar
	log          *log.Entry
}

func (p PulsarSvc) New() pulsar.Client {
	pulsarClient, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               p.pulsarConfig.Url,
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})

	if err != nil {
		p.log.WithFields(logrus.Fields{"method": "NewPulsarClient", "message": err.Error()}).Log(logrus.FatalLevel)

	}

	return pulsarClient
}

func (p PulsarSvc) GetPulsarProducer(producerOptions pulsar.ProducerOptions) pulsar.Producer {

	if producerInstance == nil {
		producer, err := p.pulsarClient.CreateProducer(producerOptions)

		if err != nil {
			p.log.Fatal(err)
		}

		producerInstance = &producer
	}
	return *producerInstance
}

func (p PulsarSvc) GetPulsarConsumer(consumerOptions pulsar.ConsumerOptions) pulsar.Consumer {
	consumer, err := p.pulsarClient.Subscribe(consumerOptions)
	if err != nil {
		p.log.Fatal(err)
	}
	return consumer
}

func NewPulsarManager(pulsarClient pulsar.Client, pulsarConfig common.Pulsar, log *log.Entry) PulsarManager {
	return &PulsarSvc{pulsarClient: pulsarClient, pulsarConfig: pulsarConfig, log: log}
}
