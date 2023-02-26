package mongodb

import (
	"github.com/cemayan/url-shortener/common"
	"github.com/cemayan/url-shortener/common/adapters/database"
	"github.com/cemayan/url-shortener/common/domain"
	"github.com/cemayan/url-shortener/common/ports/output"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApiRepoSvc struct {
	mongoClient *mongo.Client
	configs     common.Mongo
	log         *log.Entry
}

func (a ApiRepoSvc) CreateEvent(event domain.Events) error {
	collection := a.mongoClient.Database(a.configs.DbName).Collection("events")
	_, err := collection.InsertOne(database.MongoDBContext, &event)
	if err != nil {
		a.log.WithFields(log.Fields{"method": "SaveRecord"}).Errorf("An errror occurred %v", err.Error())
		return err
	}

	return nil
}

func NewApiRepo(mongoClient *mongo.Client, configs common.Mongo, log *log.Entry) output.MongoPort {
	return &ApiRepoSvc{
		mongoClient: mongoClient,
		configs:     configs,
		log:         log,
	}
}
