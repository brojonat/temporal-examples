package temporal

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	QueryTontineState      = "query-tontine-state"
	SignalParticipantAlive = "signal-participant-alive"
)

// Participant represents a member of the tontine.
type Participant struct {
	Name         string
	Contribution float64
	Alive        bool
}

// TontineState represents the current state of the tontine.
type TontineState struct {
	Participants []Participant
	TotalFund    float64
	LastPayout   time.Time
}

// TontineWorkflow is the Temporal workflow that coordinates the tontine.
func TontineWorkflow(ctx workflow.Context, participants []Participant, payoutInterval time.Duration) error {
	// Initialize tontine state
	state := TontineState{
		Participants: participants,
		TotalFund:    0,
		LastPayout:   workflow.Now(ctx),
	}

	// Set activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Collect initial contributions
	for i, p := range state.Participants {
		var contribution float64
		err := workflow.ExecuteActivity(ctx, CollectContributionActivity, p.Name, p.Contribution).Get(ctx, &contribution)
		if err != nil {
			return err
		}
		state.TotalFund += contribution
		state.Participants[i].Contribution = contribution
	}

	// Periodic payouts
	for len(state.Participants) > 1 {
		workflow.Sleep(ctx, payoutInterval)

		// Distribute payouts
		err := workflow.ExecuteActivity(ctx, PayoutActivity, state).Get(ctx, &state.TotalFund)
		if err != nil {
			return err
		}

		// Remove deceased participants
		var updatedParticipants []Participant
		for _, p := range state.Participants {
			var alive bool
			err := workflow.ExecuteActivity(ctx, CheckIfAliveActivity, p.Name).Get(ctx, &alive)
			if err != nil {
				return err
			}
			if alive {
				updatedParticipants = append(updatedParticipants, p)
			}
		}
		state.Participants = updatedParticipants
	}

	// Handle the final survivor
	if len(state.Participants) == 1 {
		return workflow.ExecuteActivity(ctx, FinalPayoutActivity, state.Participants[0].Name, state.TotalFund).Get(ctx, nil)
	}

	return nil
}
