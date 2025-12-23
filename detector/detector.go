// metadata.go
package detector

import "time"

type TargetType string

const (
	HTTP     TargetType    = "http"
	TCP      TargetType    = "tcp"
	Interval time.Duration = time.Second * time.Duration(300)
)

type Prober interface {
	Probe() Report
}

type ScheduledProbe struct {
	Prober Prober
}

type HTTPTarget struct {
	Name    string
	URL     string
	Timeout time.Duration
}

type TCPTarget struct {
	Name    string
	Address string
	Timeout time.Duration
}

type Report struct {
	Target    string        `json:"target"`
	Success   bool          `json:"success"`
	Error     error         `json:"error,omitempty"`
	Latency   time.Duration `json:"latency"`
	CheckedAT time.Time     `json:"checked_at"`
}
