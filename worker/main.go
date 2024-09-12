package worker

import (
	"context"
	"log"
	"log/slog"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const TaskQueue = "temporal-examples"

func RunWorker(ctx context.Context, l *slog.Logger, thp string) error {
	// connect to temporal
	c, err := client.Dial(client.Options{
		Logger:   l,
		HostPort: thp,
	})
	if err != nil {
		log.Fatalf("Couldn't initialize Temporal client. Exiting.\nError: %s", err)
	}
	defer c.Close()

	// register workflows
	w := worker.New(c, TaskQueue, worker.Options{})

	// register activities
	return w.Run(worker.InterruptCh())

}
