package temporal

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func RunHeartWF(ctx workflow.Context) error {
	// Run the activity for some arbitrarily long period. If the activity ceases
	// to report progress (i.e., ceases to heartbeat) for >5 seconds, then we
	// consider the activity "dead" and we want to error out. For the purposes
	// of this example, under the hood, activity has some deterministic
	// process that will control whether or not it reports progress.
	aopts := workflow.ActivityOptions{
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 1},
		StartToCloseTimeout: 60 * time.Minute,
		HeartbeatTimeout:    5 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, aopts)
	workflow.ExecuteActivity(ctx, RunHeartActivity).Get(ctx, nil)

	// Regardless of what caused the activity error, we want to restart this
	// workflow as new, so return a new ContinueAsNewError. You can pass in
	// parameters here as you normally would calling a workflow.
	return workflow.NewContinueAsNewError(ctx, RunHeartWF)
}
