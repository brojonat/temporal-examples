package temporal

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// WorkflowPoll is a workflow that runs for some specified time and receives
// incoming votes on specified poll options. The current state of the poll is
// queryable. At the end of the poll, the workflow sends the results via HTTP
// (i.e., webhook) until it receives a 200.

const (
	// query types
	QueryTypeState = "state"

	// signal types
	SignalTypeVote = "vote"
)

type PollResult struct {
	Prompt string             `json:"prompt"`
	Votes  map[string]float64 `json:"votes"`
}

type RunPollWFRequest struct {
	StartTime time.Time     `json:"start_time"`
	Duration  time.Duration `json:"duration"`
	Prompt    string        `json:"prompt"`
	Options   []string      `json:"options"`
	Webhook   string        `json:"webhook"`
}

type PollVote struct {
	Prompt string  `json:"prompt"`
	Option string  `json:"option"`
	Amount float64 `json:"amount"`
}

func RunPollWF(ctx workflow.Context, r RunPollWFRequest) error {
	// register a handler to return the current poll state
	results := PollResult{Prompt: r.Prompt, Votes: make(map[string]float64)}
	for _, o := range r.Options {
		results.Votes[o] = 0.
	}
	err := workflow.SetQueryHandler(ctx, QueryTypeState, func() (PollResult, error) {
		return results, nil
	})
	if err != nil {
		return err
	}

	// initialization for main selector loop
	doLoop := true
	var signal PollVote
	selector := workflow.NewSelector(ctx)

	// receive poll votes
	bidChan := workflow.GetSignalChannel(ctx, SignalTypeVote)
	selector.AddReceive(bidChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &signal)
		if _, ok := results.Votes[signal.Option]; !ok {
			return
		}
		results.Votes[signal.Option] += signal.Amount
	})

	// receive poll over; uses a separate goroutine that will block until the
	// poll is over before sending on the pollOverChan.
	pollOverChan := workflow.NewChannel(ctx)
	workflow.Go(ctx, func(ictx workflow.Context) {
		wait := time.Until(time.Now().Add(r.Duration))
		workflow.AwaitWithTimeout(ictx, wait, func() bool { return false })
		pollOverChan.Send(ictx, nil)
	})
	selector.AddReceive(pollOverChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, nil)
		doLoop = false
	})

	// loop receive votes until the poll is over
	for doLoop {
		selector.Select(ctx)
	}

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
	err = workflow.ExecuteActivity(ctx, RunPollCompleteWebhook, r.Webhook, results).Get(ctx, nil)
	return err
}
