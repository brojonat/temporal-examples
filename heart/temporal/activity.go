package temporal

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

func RunHeartActivity(ctx context.Context) error {
	// setup some timers
	ticker := time.NewTicker(time.Second)
	ticks := 0
	maxTicks := 20

	for {
		select {
		case t := <-ticker.C:
			// Heartbeat on every tick to indicate progress until we hit
			// maxTicks, at which point we stop reporting progress to simulate
			// the activity process dying.
			if ticks < maxTicks {
				activity.RecordHeartbeat(ctx, t)
			}
			ticks++
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
