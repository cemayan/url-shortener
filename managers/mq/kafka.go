package mq

import (
	"github.com/cemayan/url-shortener/common"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type KafkaManager interface {
	GetKafkaWriter(addr string, topic string) *kafka.Writer
	GetKafkaReader(addr string, topic string) *kafka.Reader
}
type KafkaSvc struct {
	kafkaConfigs common.Kafka
	log          *log.Entry
}

func (k KafkaSvc) GetKafkaWriter(addr string, topic string) *kafka.Writer {

	return &kafka.Writer{
		Addr:                   kafka.TCP(addr),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		MaxAttempts:            3,
		ErrorLogger:            kafka.LoggerFunc(k.log.Errorf),
	}

}

func (k KafkaSvc) GetKafkaReader(addr string, topic string) *kafka.Reader {

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{addr},
		Topic:       topic,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		ErrorLogger: kafka.LoggerFunc(k.log.Errorf),
	})
}

func NewKafkaManager(kafkaConfigs common.Kafka, log *log.Entry) KafkaManager {
	return &KafkaSvc{kafkaConfigs: kafkaConfigs, log: log}
}
