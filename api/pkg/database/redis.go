package database

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
)

var Redis *redis.Client

var (
	ErrorNotFound = redis.Nil
)

type RedisOptions struct {
	Sentinel      bool
	MasterName    string
	SentinelAddrs []string

	Addr string

	Password string
	DB       int
}

func ConnectRedis(options RedisOptions) {
	if options.Sentinel {
		Redis = redis.NewFailoverClient(buildSentinelOptions(options))
	} else {
		Redis = redis.NewClient(buildRedisOptions(options))
	}
}

func buildSentinelOptions(options RedisOptions) *redis.FailoverOptions {
	return &redis.FailoverOptions{
		MasterName:    options.MasterName,
		SentinelAddrs: options.SentinelAddrs,
		Password:      options.Password,
		DB:            options.DB,
	}
}

func buildRedisOptions(options RedisOptions) *redis.Options {
	return &redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	}
}

func Set(key string, value interface{}, exp time.Duration) error {
	return Redis.Set(context.Background(), key, value, exp).Err()
}

func Get(key string) (string, error) {
	val, err := Redis.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", ErrorNotFound
	}
	return val, err
}
