package temporal

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	QueryTontineState      = "query-tontine-state"
	SignalParticipantAlive = "signal-participant-alive"
)

type SignalParticipant struct {
	Name  string
	Alive bool
}

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

	// Start a goroutine to handle the signal reception loop
	workflow.Go(ctx, func(ctx workflow.Context) {
		// setup listener for deceased participants
		signalChannel := workflow.GetSignalChannel(ctx, SignalParticipantAlive)
		// Create a workflow.Selector to handle multiple events
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(signalChannel, func(c workflow.ReceiveChannel, more bool) {
			var sp SignalParticipant
			c.Receive(ctx, &sp)
			for i, p := range state.Participants {
				if p.Name == sp.Name {
					state.Participants[i].Alive = sp.Alive
				}
			}
		})
		for {
			selector.Select(ctx)
		}
	})

	// Periodic payouts
	for len(state.Participants) > 1 {
		workflow.Sleep(ctx, payoutInterval)

		// Remove deceased participants
		var updatedParticipants []Participant
		for _, p := range state.Participants {
			if p.Alive {
				updatedParticipants = append(updatedParticipants, p)
			}
		}
		state.Participants = updatedParticipants

		// terminate if no participants
		if len(state.Participants) == 0 {
			break
		}

		// distribute regular payouts
		if err := workflow.ExecuteActivity(ctx, RegularPayoutActivity, state).Get(ctx, &state.TotalFund); err != nil {
			return err
		}
	}

	// distribute to the final remaining participant; if there's no survivors, distribute back to the bank
	rp := "bank"
	if len(state.Participants) == 1 {
		rp = state.Participants[0].Name
	}
	return workflow.ExecuteActivity(ctx, FinalPayoutActivity, rp, state.TotalFund).Get(ctx, nil)
}
