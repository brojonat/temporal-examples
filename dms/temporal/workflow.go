package temporal

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// WorkflowAuction is a workflow that runs for some specified time and receives
// incoming bids. The current top bid is queryable. At the end of the auction,
// the workflow sends the results via HTTP (i.e., webhook) until it receives a
// 200.

const (
	// query types
	QueryTypeState = "state"

	// signal types
	SignalTypeDeactivate = "deactivate"
)

type RunDMSWFRequest struct {
	ID        string        `json:"id"`
	StartTime time.Time     `json:"start_time"`
	Duration  time.Duration `json:"duration"`
	Message   string        `json:"message"`
	Webhook   string        `json:"webhook"`
}

type DMSTimeoutPayload struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func RunDMSWF(ctx workflow.Context, r RunDMSWFRequest) error {

	// initialization for main selector loop
	doLoop := true
	timedOut := false
	deactivated := false
	selector := workflow.NewSelector(ctx)

	// register a handler to return the current top bid
	err := workflow.SetQueryHandler(ctx, QueryTypeState, func() (string, error) {
		if deactivated {
			return "switch was deactivated", nil
		}
		if timedOut {
			return "switch timed out", nil
		}
		return fmt.Sprintf("%s until timeout", time.Until(r.StartTime.Add(r.Duration))), nil
	})
	if err != nil {
		return err
	}

	// receive deactivation
	deactivateChan := workflow.GetSignalChannel(ctx, SignalTypeDeactivate)
	selector.AddReceive(deactivateChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, nil)
		doLoop = false
		deactivated = true
	})

	// receive timeout; uses a separate goroutine that will block until the
	// dms times out before sending on the timeoutChan.
	timeoutChan := workflow.NewChannel(ctx)
	workflow.Go(ctx, func(ictx workflow.Context) {
		wait := time.Until(time.Now().Add(r.Duration))
		workflow.AwaitWithTimeout(ictx, wait, func() bool { return false })
		timeoutChan.Send(ictx, nil)
	})
	selector.AddReceive(timeoutChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, nil)
		doLoop = false
		timedOut = true
	})

	// loop receive bids until the dms is deactivated or times out
	for doLoop {
		selector.Select(ctx)
	}

	if timedOut {
		// send the webhook with the results
		rp := temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 5.0,
			MaximumInterval:    time.Second * 100,
			MaximumAttempts:    0, // Unlimited
		}
		aopts := workflow.ActivityOptions{
			StartToCloseTimeout: 60 * time.Minute,
			RetryPolicy:         &rp,
			HeartbeatTimeout:    60 * time.Second,
		}
		ctx = workflow.WithActivityOptions(ctx, aopts)
		err = workflow.ExecuteActivity(ctx, RunDMSTimeoutWebhook, r.Webhook, DMSTimeoutPayload{ID: r.ID, Message: r.Message}).Get(ctx, nil)
	}

	return err
}
