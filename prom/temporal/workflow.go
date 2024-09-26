package temporal

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func RunPromWF(ctx workflow.Context) error {
	// run the activity for some arbitrarily long period
	aopts := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, aopts)
	err := workflow.ExecuteActivity(ctx, RunPromActivity).Get(ctx, nil)
	return err
}
