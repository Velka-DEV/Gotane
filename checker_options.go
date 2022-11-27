package gotane

type CheckerOptions struct {
	// MaxConcurrentChecks is the maximum number of concurrent checks
	Threads int `json:"threads"`

	// MaxRetries is the maximum number of retries for a proxy
	MaxRetries int `json:"max_retries"`
}
