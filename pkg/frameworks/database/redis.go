package database

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/marllef/awesome-queue/pkg/configs"
)

type RedisDB interface {
	GetClient() *redis.Client
	GetContext() context.Context
}

type redisDB struct {
	client  *redis.Client
	context context.Context
}

func NewRedisDB(ctx context.Context) (container *redisDB, err error) {
	if err = configs.LoadEnvs(); err != nil {
		return nil, err
	}

	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", configs.Env.RedisHost, configs.Env.RedisPort),
		Password: configs.Env.RedisPassword,
		DB:       configs.Env.RedisDatabase,
	}

	return &redisDB{
		client:  redis.NewClient(opts),
		context: ctx,
	}, nil
}

func (rdb *redisDB) GetClient() *redis.Client {
	return rdb.client
}

func (rdb *redisDB) GetContext() context.Context {
	return rdb.context
}
