package temporal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func doRequestPayout(ctx context.Context, name string, amount float64, endpoint string) error {
	payload := struct {
		Name   string  `json:"name"`
		Amount float64 `json:"amount"`
	}{
		Name:   name,
		Amount: amount,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	r, err := http.NewRequest(http.MethodPost, endpoint+"/payout", bytes.NewReader(b))
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return nil
	}
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("bad response (%d) and error reading body: %w", res.StatusCode, err)
	}
	return fmt.Errorf("bad response (%d): %s", res.StatusCode, b)
}

// RegularPayoutActivity distributes payouts to all living participants.
func RegularPayoutActivity(ctx context.Context, state TontineState, amount float64, endpoint string) (float64, error) {
	var g errgroup.Group
	for _, p := range state.Participants {
		g.Go(func() error {
			err := doRequestPayout(ctx, p.Name, amount, endpoint)
			if err == nil {
				state.TotalFund -= amount
			}
			return err
		})
	}

	// Wait for all tasks to complete
	if err := g.Wait(); err != nil {
		return state.TotalFund, err
	}
	return state.TotalFund, nil
}

// FinalPayoutActivity gives the remaining fund to the final survivor.
func FinalPayoutActivity(ctx context.Context, name string, amount float64, endpoint string) error {
	return doRequestPayout(ctx, name, amount, endpoint)
}
