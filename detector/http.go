// http.go
package detector

import (
	"fmt"
	"net/http"
	"time"
)

func (t HTTPTarget) Probe() Report {
	target := fmt.Sprintf("%v: %v", t.Name, t.URL)
	start := time.Now()

	client := http.Client{
		Timeout: t.Timeout,
	}

	r, err := client.Get(t.URL)
	latency := time.Since(start)

	if err != nil {
		return Report{
			Target:    target,
			Success:   false,
			Error:     err,
			Latency:   latency,
			CheckedAT: time.Now(),
		}
	}
	defer r.Body.Close()

	success := r.StatusCode < 500

	return Report{
		Target:    target,
		Success:   success,
		Error:     err,
		Latency:   latency,
		CheckedAT: time.Now(),
	}
}
