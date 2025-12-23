// tcp.go
package detector

import (
	"fmt"
	"net"
	"time"
)

func (t TCPTarget) Probe() Report {
	target := fmt.Sprintf("%v: %v", t.Name, t.Address)
	start := time.Now()

	conn, err := net.DialTimeout("tcp", t.Address, t.Timeout)
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

	conn.Close()

	return Report{
		Target:    target,
		Success:   true,
		Error:     nil,
		Latency:   latency,
		CheckedAT: time.Now(),
	}
}
