package queue

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/marllef/awesome-queue/modules/queue/internal/model"
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
	group   string
}

var queues = make(map[string]*queue)

func NewQueue(name string, container database.RedisDB, log *logger.Logger) *queue {
	newQueue := &queue{
		name:    name,
		client:  container.GetClient(),
		context: container.GetContext(),
		log:     log,
		group:   "",
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

func (q *queue) JoinGroup(group string) (err error) {
	q.client.XGroupCreateMkStream(q.context, q.name, group, "0")
	if err = q.client.XGroupCreateConsumer(q.context, q.name, group, configs.Env.RedisConsumerID).Err(); err != nil {
		return err
	}
	q.group = group

	q.log.Infof("Successfuly joined on group '%s'", group)
	return nil
}

func (q *queue) Proccess(handler model.HandlerFunc) error {

	args := &redis.XReadGroupArgs{
		Group:    q.group,
		Consumer: configs.Env.RedisConsumerID,
		Streams:  []string{q.name, ">"},
		Block:    0,
	}

	for {
		entries, err := q.client.XReadGroup(q.context, args).Result()

		if err != nil {
			return err
		}

		for _, entrie := range entries {
			for _, message := range entrie.Messages {
				messageID := message.ID
				values := message.Values
				q.log.Infof("[%s] New message received. Processing...", q.name)
				if err := handler(messageID, values); err != nil {
					q.log.Errorf("Message [%s] - Failed to proccess message: %v", messageID, err)
					continue
				}

				if err := q.client.XAck(q.context, q.name, q.group, messageID).Err(); err != nil {
					q.log.Warnf("Message [%s] - Failed to ack message: %s", messageID, err)
					continue
				}

				q.log.Infof("Message [%s] - Success to proccess message.", messageID)
			}
		}
	}
}
