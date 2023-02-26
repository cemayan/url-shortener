package redis

import (
	"context"
	"github.com/cemayan/url-shortener/common/ports/output"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"time"
)

var Ctx = context.Background()

// A RedisSvc  contains the required dependencies for this service
type RedisSvc struct {
	client *redis.Client
	log    *log.Entry
}

func (r RedisSvc) Remove(key string) error {
	//TODO implement me
	panic("implement me")
}

// Get is command of Redis's Get
func (r RedisSvc) Get(key string) (string, error) {
	return r.client.Get(Ctx, key).Result()
}

// Set is command of Redis's Set
func (r RedisSvc) Set(key string, data string) error {
	err := r.client.Set(Ctx, key, data, 30*time.Minute).Err()
	if err != nil {
		r.log.WithFields(log.Fields{"method": "Set"}).Infof("An error occurred when set value to redis %v", err)
	}
	return err
}

func NewRedisHandler(client *redis.Client, log *log.Entry) output.RedisPort {
	return &RedisSvc{client: client, log: log}
}
