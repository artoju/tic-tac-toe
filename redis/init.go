package redis

import (
	"github.com/artoju/tic-tac-toe/config"
	"github.com/go-redis/redis/v8"
)

func Init(conf *config.Config) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     conf.RedisHandler.Address,
		Password: conf.RedisHandler.Password,
		DB:       conf.RedisHandler.Database,
	})

	return client, nil
}
