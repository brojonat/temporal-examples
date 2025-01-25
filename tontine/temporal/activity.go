package temporal

import (
	"context"
	"fmt"
)

// CollectContributionActivity collects the initial contribution from a participant.
func CollectContributionActivity(ctx context.Context, name string, contribution float64) (float64, error) {
	fmt.Printf("Collecting contribution from %s: %.2f\n", name, contribution)
	return contribution, nil
}

// PayoutActivity distributes payouts to all living participants.
func PayoutActivity(ctx context.Context, state TontineState) (float64, error) {
	if len(state.Participants) == 0 {
		return state.TotalFund, nil
	}

	payout := state.TotalFund / float64(len(state.Participants))
	fmt.Printf("Distributing payout of %.2f to each participant\n", payout)

	for _, p := range state.Participants {
		fmt.Printf("Payout to %s: %.2f\n", p.Name, payout)
	}

	return state.TotalFund - (payout * float64(len(state.Participants))), nil
}

// CheckIfAliveActivity checks if a participant is still alive.
func CheckIfAliveActivity(ctx context.Context, name string) (bool, error) {
	fmt.Printf("Checking if %s is alive\n", name)
	// For simplicity, simulate with random logic or a fixed value.
	return true, nil
}

// FinalPayoutActivity gives the remaining fund to the final survivor.
func FinalPayoutActivity(ctx context.Context, name string, totalFund float64) error {
	fmt.Printf("Final payout to %s: %.2f\n", name, totalFund)
	return nil
}
