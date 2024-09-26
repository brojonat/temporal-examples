package worker

import (
	"context"
	"log"
	"log/slog"

	auction "github.com/brojonat/temporal-examples/auction/temporal"
	dms "github.com/brojonat/temporal-examples/dms/temporal"
	heart "github.com/brojonat/temporal-examples/heart/temporal"
	poll "github.com/brojonat/temporal-examples/poll/temporal"
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
	w.RegisterWorkflow(auction.RunAuctionWF)
	w.RegisterWorkflow(poll.RunPollWF)
	w.RegisterWorkflow(dms.RunDMSWF)
	w.RegisterWorkflow(heart.RunHeartWF)

	// register activities
	w.RegisterActivity(auction.RunAuctionCompleteWebhook)
	w.RegisterActivity(poll.RunPollCompleteWebhook)
	w.RegisterActivity(dms.RunDMSTimeoutWebhook)
	w.RegisterActivity(heart.RunHeartActivity)
	return w.Run(worker.InterruptCh())

}
