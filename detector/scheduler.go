// scheduler.go
package detector

import (
	"context"
	"time"
)

func RunProbe(ctx context.Context, p Prober, out chan<- Report) {
	ticker := time.NewTicker(Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			out <- p.Probe()
		}
	}
}
