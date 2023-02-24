package util

import (
	"github.com/cemayan/url-shortener/common"
	"github.com/google/uuid"
	"time"
)

func GetUnixTime() int64 {
	now := time.Now()
	return now.Unix()
}

func GetEventId() uuid.UUID {
	return uuid.New()
}

func GetEventMetadata() *common.EventModel {

	eventId := GetEventId().String()
	eventName := common.UrlEncoded
	eventDate := GetUnixTime()

	return &common.EventModel{
		AggregateId: eventId,
		EventDate:   eventDate,
		EventName:   eventName,
	}
}
