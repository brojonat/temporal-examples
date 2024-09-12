package temporal

import (
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
	QueryTypeTopBid = "top-bid"

	// signal types
	SignalTypePlaceBid = "place-bid"
)

type QueryResultTopBid struct {
	Bidder string  `json:"bidder"`
	Amount float64 `json:"amount"`
}

type RunAuctionWFRequest struct {
	StartTime    time.Time     `json:"start_time"`
	Duration     time.Duration `json:"duration"`
	Item         string        `json:"item"`
	ReservePrice float64       `json:"reserve_price"`
	Webhook      string        `json:"webhook"`
}

type AuctionBid struct {
	Item   string  `json:"item"`
	Bidder string  `json:"bidder"`
	Amount float64 `json:"amount"`
}

func RunAuctionWF(ctx workflow.Context, r RunAuctionWFRequest) error {
	// register a handler to return the current top bid
	topBid := AuctionBid{Item: r.Item}
	err := workflow.SetQueryHandler(ctx, QueryTypeTopBid, func() (AuctionBid, error) {
		return topBid, nil
	})
	if err != nil {
		return err
	}

	doLoop := true
	var signal AuctionBid

	bidChan := workflow.GetSignalChannel(ctx, SignalTypePlaceBid)
	auctionOverChan := workflow.NewChannel(ctx)

	selector := workflow.NewSelector(ctx)

	// receive auction bids
	selector.AddReceive(bidChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &signal)
		if signal.Amount > topBid.Amount {
			topBid = signal
		}
	})
	// receive auction over
	selector.AddReceive(auctionOverChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, nil) // FIXME?
		doLoop = false
	})
	// loop receive bids until the auction is over
	workflow.Go(ctx, func(ctx workflow.Context) {
		for doLoop {
			selector.Select(ctx)
		}
	})

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
	err = workflow.ExecuteActivity(ctx, RunAuctionCompleteWebhook, r.Webhook, topBid).Get(ctx, nil)
	return err
}
