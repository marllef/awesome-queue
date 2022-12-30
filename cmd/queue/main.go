package main

import (
	"context"

	"github.com/marllef/awesome-queue/modules/queue"
	"github.com/marllef/awesome-queue/pkg/frameworks/database"
	"github.com/marllef/awesome-queue/pkg/utils/logger"
)

func main() {
	log := logger.Default()
	ctx := context.Background()

	// New Redis Connection
	container, err := database.NewRedisDB(ctx)
	if err != nil {
		log.Errorf("Error on connect to Redis: %v", err)
		return
	}

	// New queue
	messageQueue := queue.NewQueue("Messages", container, log)

	if err = messageQueue.JoinConsumerGroup("message-consumer-group"); err != nil {
		log.Errorf("Error on create group on Queue: %s", err)
	}

}
