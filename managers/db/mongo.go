package db

import (
	"context"
	"github.com/cemayan/url-shortener/common"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBContext = context.TODO()

type MongodbManager interface {
	New() *mongo.Client
}

type MongoDBService struct {
	configs common.Mongo
	_log    *log.Entry
}

// New  serves to connect to db
// When DB connection is successful then model migration is started
func (d MongoDBService) New() *mongo.Client {
	client, err := mongo.Connect(MongoDBContext, options.Client().ApplyURI(d.configs.Uri))
	if err != nil {
		panic("failed to connect database")
	}
	d._log.WithFields(log.Fields{"service": "database"}).Infoln("Connection Opened to MongoDB")

	return client
}

func NewMongodbManager(configs common.Mongo, _log *log.Entry) MongodbManager {
	return &MongoDBService{configs: configs, _log: _log}
}
