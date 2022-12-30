package queue

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/marllef/awesome-queue/pkg/configs"
	"github.com/marllef/awesome-queue/pkg/frameworks/database"
	"github.com/marllef/awesome-queue/pkg/utils/logger"
)

type Queue interface{}

type queue struct {
	name    string
	client  *redis.Client
	context context.Context
	log     *logger.Logger
}

var queues = make(map[string]*queue)

func NewQueue(name string, container database.RedisDB, log *logger.Logger) *queue {
	newQueue := &queue{
		name:    name,
		client:  container.GetClient(),
		context: container.GetContext(),
		log:     log,
	}

	queues[name] = newQueue
	return newQueue
}

func GetAllQueues() map[string]*queue {
	return queues
}

func (q *queue) Publish(values map[string]interface{}) error {
	args := &redis.XAddArgs{
		Stream: q.name,
		ID:     "",
		Values: values,
	}

	if err := q.client.XAdd(q.context, args).Err(); err != nil {
		return err
	}

	return nil
}

func (q *queue) JoinConsumerGroup(group_name string) (err error) {
	q.client.XGroupCreateMkStream(q.context, q.name, group_name, "0")
	if err = q.client.XGroupCreateConsumer(q.context, q.name, group_name, configs.Env.RedisConsumerID).Err(); err != nil {
		return err
	}
	return nil
}

func (q *queue) Consume() error {

	return nil
}

func (q *queue) Proccess(handle func() error) error {
	return handle()
}
