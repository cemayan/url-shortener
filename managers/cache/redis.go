package cache

import (
	"github.com/cemayan/url-shortener/common"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type RedisManager interface {
	New() *redis.Client
}

type RedisService struct {
	configs common.Redis
	_log    *log.Entry
}

// New create redis client
func (r RedisService) New() *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     r.configs.Address + ":" + r.configs.AddressPort,
		Password: "",
		DB:       0,
	})
}

func NewRedisManager(configs common.Redis, _log *log.Entry) RedisManager {
	return &RedisService{configs: configs, _log: _log}
}
