package queue

import (
	"context"

	"github.com/marllef/awesome-queue/pkg/configs"
	"github.com/marllef/awesome-queue/pkg/queue/types"
	"github.com/marllef/awesome-queue/pkg/utils/logger"
)

type QueueConsumer interface{}

type queueConsumer struct {
	id      string
	context context.Context
	queue   *queue
	log     *logger.Logger
	groups  []string
}

func NewConsumer(queue *queue, groups ...string) *queueConsumer {
	consumer := &queueConsumer{
		id:      configs.Env.RedisConsumerID,
		queue:   queue,
		log:     logger.Default(),
		groups:  []string{},
		context: queue.context,
	}

	// Join on consumer groups
	joinedGroups := make([]string, 0)
	for _, group := range groups {
		err := consumer.JoinGroup(group)
		if err != nil {
			continue
		}

		joinedGroups = append(joinedGroups, group)
	}
	
	consumer.groups = joinedGroups

	return consumer
}

func (qc *queueConsumer) JoinGroup(group string) (err error) {
	qc.queue.client.XGroupCreateMkStream(qc.context, qc.queue.key, group, "0")
	if err = qc.queue.client.XGroupCreateConsumer(qc.context, qc.queue.key, group, qc.id).Err(); err != nil {
		return err
	}
	qc.groups = append(qc.groups, group)

	qc.log.Infof("Successfuly joined on group '%s'", group)
	return nil
}

func (qc *queueConsumer) LeaveGroup(group string) (err error) {
	if err = qc.queue.client.XGroupDelConsumer(qc.context, qc.queue.key, group, qc.id).Err(); err != nil {
		return err
	}

	qc.groups = make([]string, 0)
	for _, value := range qc.groups {
		if value == group {
			continue
		}
		qc.groups = append(qc.groups, value)
	}

	qc.log.Infof("Consumer successfuly leave of the group '%s'", group)
	return nil
}

func (qc *queueConsumer) Consume(handler types.HandlerFunc) {
	for _, group := range qc.groups {
		go qc.queue.Proccess(group, handler)
	}
}

// Get current consumer context
func (qc *queueConsumer) GetContext() context.Context {
	return qc.context
}

// Set current consumer context
func (qc *queueConsumer) SetContext(ctx context.Context) {
	qc.context = ctx
}

// Get current consumer logger
func (qc *queueConsumer) GetLogger() *logger.Logger {
	return qc.log
}

// Set current consumer logger
func (qc *queueConsumer) SetLogger(log *logger.Logger) {
	qc.log = log
}
