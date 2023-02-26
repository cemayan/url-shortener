package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Events struct {
	ID            primitive.ObjectID `bson:"_id"`
	EventData     primitive.M        `bson:"eventData" json:"eventData"`
	EventDate     int64              `bson:"eventDate" json:"eventDate"`
	EventName     string             `bson:"eventName" json:"eventName"`
	AggregateType string             `bson:"aggregateType" json:"aggregateType"`
	AggregateId   string             `bson:"aggregateId" json:"aggregateId"`
	UserId        string             `bson:"userId" json:"userId"`
}
