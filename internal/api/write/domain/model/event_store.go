package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventStore struct {
	ID            primitive.ObjectID `bson:"_id"`
	EventData     primitive.M        `bson:"eventData" json:"eventData"`
	EventTime     int64              `bson:"eventTime" json:"eventTime"`
	EventName     string             `bson:"eventName" json:"eventName"`
	AggregateType string             `bson:"aggregateType" json:"aggregateType"`
	AggregateID   string             `bson:"aggregateId" json:"aggregateId"`
}
