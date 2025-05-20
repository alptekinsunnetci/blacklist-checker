package models

// BlacklistResult represents the result of a DNSBL check
type BlacklistResult struct {
	IP        string `json:"ip"`
	Blacklist string `json:"blacklist"`
	Error     string `json:"error,omitempty"`
	Timestamp string `json:"timestamp"`
}

// CheckResult represents the aggregated results for an IP
type CheckResult struct {
	IP         string   `json:"ip"`
	Blacklists []string `json:"blacklists"`
	Error      string   `json:"error,omitempty"`
}

// Config represents the application configuration
type Config struct {
	Concurrency      int      `json:"concurrency"`
	Timeout          int      `json:"timeout"`
	Blacklists       []string `json:"blacklists"`
	OutputFormat     string   `json:"output_format"`
	CustomBlacklists []string `json:"custom_blacklists,omitempty"`
}
