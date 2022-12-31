package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/marllef/awesome-queue/pkg/configs"
	"github.com/marllef/awesome-queue/pkg/frameworks/database"
	"github.com/marllef/awesome-queue/pkg/queue/types"
	"github.com/marllef/awesome-queue/pkg/utils/logger"
)

type Queue interface {
	JoinGroup(group string) (err error)
	Publish(values types.Values) error
	Proccess(handler types.HandlerFunc) error
}

type queue struct {
	key     string
	client  *redis.Client
	context context.Context
	log     *logger.Logger
}

var queues = make(map[string]*queue)

func NewQueue(key string, container database.RedisDB, log *logger.Logger) *queue {
	newQueue := &queue{
		key:     key,
		client:  container.GetClient(),
		context: container.GetContext(),
		log:     log,
	}

	queues[key] = newQueue
	return newQueue
}

func GetAllQueues() map[string]*queue {
	return queues
}

func (q *queue) Publish(values map[string]interface{}) error {
	args := &redis.XAddArgs{
		Stream: q.key,
		ID:     "",
		Values: values,
		MinID: fmt.Sprintf("%d", time.Now().Add(-2 * time.Second).UnixMilli()),
	}

	if err := q.client.XAdd(q.context, args).Err(); err != nil {
		return err
	}

	return nil
}

func (qc *queue) Proccess(group string, handler types.HandlerFunc) error {

	args := &redis.XReadGroupArgs{
		Group:    group,
		Consumer: configs.Env.RedisConsumerID,
		Streams:  []string{qc.key, ">"},
		Block:    0,
	}

	for {
		entries, err := qc.client.XReadGroup(qc.context, args).Result()
		if err != nil {
			qc.log.Errorf("Failed to read...: %v", err)
			return err
		}

		for _, entrie := range entries {
			for _, message := range entrie.Messages {
				messageID := message.ID
				values := message.Values
				qc.log.Infof("[%s:proccess] New message(%s) received.", qc.key, messageID)
				if err := handler(messageID, values); err != nil {
					qc.log.Errorf("[%s] Failed to proccess message(%s): %v", qc.key, messageID, err)
					continue
				}

				if err := qc.client.XAck(qc.context, qc.key, group, messageID).Err(); err != nil {
					qc.log.Warnf("[%s:proccess] Failed to ack message(%s): %v", qc.key, messageID, err)
					continue
				}

				qc.log.Infof("[%s:proccess] Success on proccess message(%s).",qc.key, messageID)
			}
		}
	}
}
