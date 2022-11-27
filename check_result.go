package gotane

type CheckResult uint

const (
	// CheckResultInvalid is the result when the combo is invalid
	CheckResultInvalid = iota

	// CheckResultRetry is the result when the check has failed (mostly network exception)
	CheckResultRetry

	// CheckResultError is the result when an unhandled error has occurred
	CheckResultError

	// CheckResultFree is the result when the combo is free
	CheckResultFree

	// CheckResultHit is the result when the combo is premium / subscribed
	CheckResultHit

	// CheckResultLocked is the result when the combo is locked
	CheckResultLocked
)

func (c CheckResult) String() string {
	switch c {
	case CheckResultInvalid:
		return "invalid"
	case CheckResultRetry:
		return "retry"
	case CheckResultError:
		return "error"
	case CheckResultFree:
		return "free"
	case CheckResultHit:
		return "hit"
	case CheckResultLocked:
		return "locked"
	}

	return "unknown"
}
