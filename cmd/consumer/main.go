package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/marllef/awesome-queue/modules/queue"
	"github.com/marllef/awesome-queue/pkg/frameworks/database"
	"github.com/marllef/awesome-queue/pkg/utils/logger"
)

func main() {
	log := logger.Default()
	ctx := context.Background()

	// New Redis Connection
	redis, err := database.NewRedisDB(ctx)
	if err != nil {
		log.Errorf("Error on connect to Redis: %v", err)
		return
	}

	// New queue
	message_queue := queue.NewQueue("Messages", redis, log)
	notification_queue := queue.NewQueue("Notification", redis, log)

	// Join on Consumer Groups
	if err = message_queue.JoinGroup("message-consumer-group"); err != nil {
		log.Errorf("Error on join group: %s", err)
	}

	if err = notification_queue.JoinGroup("notification-consumer-group"); err != nil {
		log.Errorf("Error on join group: %s", err)
	}

	go message_queue.Proccess(func(id string, values map[string]interface{}) error {
		log.Infof("[Processing] Message Id: %s | Message Type: %v", id, values["type"])
		return nil
	})

	go notification_queue.Proccess(func(id string, values map[string]interface{}) error {
		log.Infof("[Processing] Message Id: %s | Message Type: %v", id, values["type"])
		return nil
	})

	channel := make(chan os.Signal, 5)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-channel:
			log.Infof("Interrupt received, exiting")
			return
		default:
			continue
		}
	}
}
