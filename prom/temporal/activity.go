package temporal

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

func RunPromActivity(ctx context.Context) error {

	// setup some timers
	ticker := time.NewTicker(time.Second)
	ender := time.NewTimer(time.Minute)

	// setup dummy metrics
	mh := activity.GetMetricsHandler(ctx)
	labels := map[string]string{"prom_test_label": "foo"}
	mh = mh.WithTags(labels)
	c := mh.Counter("foo-counter")
	g := mh.Gauge("foo-gauge")

	// loop for a while and emit metrics
	loop := true
	for loop {
		select {
		case t := <-ticker.C:
			c.Inc(1)
			g.Update(float64(t.Unix()))
		case <-ender.C:
			loop = false
		case <-ctx.Done():
			loop = false
		}
	}
	return nil
}
