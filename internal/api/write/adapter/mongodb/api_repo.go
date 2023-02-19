package mongodb

import (
	"github.com/cemayan/url-shortener/config/api"
	"github.com/cemayan/url-shortener/internal/api/write/adapter/database"
	"github.com/cemayan/url-shortener/internal/api/write/domain/model"
	"github.com/cemayan/url-shortener/internal/api/write/domain/port/output"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApiRepoSvc struct {
	mongoClient *mongo.Client
	configs     *api.AppConfig
	log         *log.Entry
}

func (a ApiRepoSvc) CreateUserUrl(userUrl model.UserUrl) error {
	collection := a.mongoClient.Database(a.configs.Mongo.DbName).Collection("userUrls")
	_, err := collection.InsertOne(database.MongoDBContext, &userUrl)
	if err != nil {
		a.log.WithFields(log.Fields{"method": "SaveRecord"}).Errorf("An errror occurred %v", err.Error())
		return err
	}

	return nil
}

func NewApiRepo(mongoClient *mongo.Client, configs *api.AppConfig, log *log.Entry) output.MongoPort {
	return &ApiRepoSvc{
		mongoClient: mongoClient,
		configs:     configs,
		log:         log,
	}
}
