package main

import (
	"context"
	"net/http"

	"github.com/marllef/awesome-queue/pkg/frameworks/database"
	"github.com/marllef/awesome-queue/pkg/frameworks/server"
	"github.com/marllef/awesome-queue/pkg/queue"
	"github.com/marllef/awesome-queue/pkg/queue/types"
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

	mailQueue := queue.NewQueue("mail", redis, log)

	mailConsumer := queue.NewConsumer(mailQueue, "mail-consumer-group")

	mailConsumer.Consume(func(id string, values map[string]interface{}) error {
		log.Infof("New Message Processing... | Values: %v", values)
		return nil
	})

	app := server.NewServer()

	app.AddRoute("pub", server.Route{
		Path:        "/pub",
		Middlewares: server.Middlewares{},
		Controller: func(res http.ResponseWriter, req *http.Request) {
			err := mailQueue.Publish(types.Values{
				"type":    "mail",
				"message": "Oi, Moanoite",
			})
			if err != nil {
				log.Errorf("Error on send mail: %v", err)
				res.WriteHeader(500)
				return
			}
		},
		Methods: []string{"GET"},
	})

	app.Serve()
}
