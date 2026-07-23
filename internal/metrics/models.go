package metrics

import "time"

type Metrics struct {
	TraceID   string    `json:"trace_id"`
	Timestamp time.Time `json:"timestamp"`
	Hostname  string    `json:"hostname"`
	CPU       int       `json:"cpu"`
	Memory    int       `json:"memory"`
	Disk      int       `json:"disk"`
	Uptime    string    `json:"uptime"`
	Processes int       `json:"processes"`
}
