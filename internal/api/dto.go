package api

import "time"

// Link represents a HATEOAS hypermedia link.
type Link struct {
	Rel    string `json:"rel"`
	Href   string `json:"href"`
	Method string `json:"method"`
}

// ResponseEnvelope is a base wrapper for HATEOAS support.
type ResponseEnvelope struct {
	Links []Link `json:"_links,omitempty"`
}

// HealthResponse represents the diagnostics status of the agent.
type HealthResponse struct {
	ResponseEnvelope
	Status    string    `json:"status"`
	Uptime    string    `json:"uptime"`
	Timestamp time.Time `json:"timestamp"`
}

// AgentActionResponse represents the result of start/stop/once actions.
type AgentActionResponse struct {
	ResponseEnvelope
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// ReportResponse represents metadata about a generated report file.
type ReportResponse struct {
	ResponseEnvelope
	Format    string    `json:"format"`
	Filename  string    `json:"filename"`
	Size      int64     `json:"size_bytes"`
	Timestamp time.Time `json:"timestamp"`
}