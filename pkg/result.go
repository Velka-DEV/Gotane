package pkg

import "fmt"

type CheckStatus uint

const (
	CheckStatusInvalid CheckStatus = iota
	CheckStatusHit
	CheckStatusFree
	CheckStatusError
	CheckStatusRetry
	CheckStatusLocked
	CheckStatusUnknown
)

type CheckResult struct {
	Combo    *Combo
	Proxy    string
	Status   CheckStatus
	Captures map[string]string
}

func (r CheckResult) String() string {
	return fmt.Sprintf("%s:%s", r.Combo.String(), r.Proxy)
}

func (s CheckStatus) String() string {
	switch s {
	case CheckStatusInvalid:
		return "invalid"
	case CheckStatusHit:
		return "hit"
	case CheckStatusFree:
		return "free"
	case CheckStatusError:
		return "error"
	case CheckStatusRetry:
		return "retry"
	case CheckStatusLocked:
		return "locked"
	case CheckStatusUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}
