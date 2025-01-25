package temporal

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/testsuite"
)

func TestTontineWorkflow(t *testing.T) {
	ts := testsuite.WorkflowTestSuite{}
	env := ts.NewTestWorkflowEnvironment()

	env.RegisterActivity(CollectContributionActivity)
	env.RegisterActivity(PayoutActivity)
	env.RegisterActivity(CheckIfAliveActivity)
	env.RegisterActivity(FinalPayoutActivity)

	participants := []Participant{
		{Name: "Alice", Contribution: 100, Alive: true},
		{Name: "Bob", Contribution: 150, Alive: true},
		{Name: "Charlie", Contribution: 200, Alive: true},
	}

	env.ExecuteWorkflow(TontineWorkflow, participants, 24*time.Hour) // Simulate daily payouts

	assert.True(t, env.IsWorkflowCompleted())
	assert.NoError(t, env.GetWorkflowError())
}

func TestTontinePayoutDistribution(t *testing.T) {
	// Test the PayoutActivity independently
	state := TontineState{
		Participants: []Participant{
			{Name: "Alice", Contribution: 100, Alive: true},
			{Name: "Bob", Contribution: 150, Alive: true},
		},
		TotalFund:  250,
		LastPayout: time.Now(),
	}

	updatedTotal, err := PayoutActivity(context.Background(), state)
	assert.NoError(t, err)

	// Verify that the total fund decreases after payout
	assert.Less(t, updatedTotal, state.TotalFund)
}

func TestTontineSurvivorship(t *testing.T) {
	// Test the CheckIfAliveActivity independently
	participants := []Participant{
		{Name: "Alice", Contribution: 100, Alive: true},
		{Name: "Bob", Contribution: 150, Alive: true},
		{Name: "Charlie", Contribution: 200, Alive: false},
	}

	var updatedParticipants []Participant
	for _, participant := range participants {
		alive, err := CheckIfAliveActivity(context.Background(), participant.Name)
		assert.NoError(t, err)
		if alive {
			updatedParticipants = append(updatedParticipants, participant)
		}
	}

	// Verify that only alive participants remain
	assert.Equal(t, 2, len(updatedParticipants))
	assert.NotContains(t, updatedParticipants, Participant{Name: "Charlie", Contribution: 200, Alive: false})
}

func TestFinalPayout(t *testing.T) {
	// Test the FinalPayoutActivity independently
	participant := Participant{Name: "Alice", Contribution: 100, Alive: true}
	totalFund := 1000.0

	err := FinalPayoutActivity(context.Background(), participant.Name, totalFund)
	assert.NoError(t, err)
}
