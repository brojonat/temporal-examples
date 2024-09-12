package worker

import (
	"context"
	"log"
	"log/slog"

	"github.com/brojonat/temporal-examples/auction/temporal"
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
	w.RegisterWorkflow(temporal.RunAuctionWF)

	// register activities
	w.RegisterActivity(temporal.RunAuctionCompleteWebhook)
	return w.Run(worker.InterruptCh())

}
