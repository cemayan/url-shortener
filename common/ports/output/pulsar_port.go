package output

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cemayan/url-shortener/common"
)

type PulsarPort interface {
	PublishEvent(event common.EventModel, topicName string) error
	ConsumeEvent(topicName string, subscriptionName string, channel chan pulsar.ConsumerMessage) pulsar.Consumer
}
