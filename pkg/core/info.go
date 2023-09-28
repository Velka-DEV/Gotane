package core

import "time"

type CheckerState uint

const (
	CheckerStateIdle CheckerState = iota
	CheckerStateRunning
	CheckerStatePaused
	CheckerStateStopped
)

type CheckerInfo struct {
	Total   uint64
	Checked uint64
	Hit     uint64
	Invalid uint64
	Free    uint64
	Error   uint64
	Locked  uint64
	Retry   uint64

	Progress float64
	Cpm      float64

	State     CheckerState
	StartTime time.Time
}

func NewCheckerInfo(total uint64) *CheckerInfo {
	return &CheckerInfo{
		Total: total,
	}
}

func (c *CheckerInfo) StartUpdateLoop() {
	c.StartTime = time.Now()
	c.State = CheckerStateRunning

	for {
		if c.State == CheckerStateStopped {
			break
		}

		c.UpdateProgress()
		c.UpdateCpm()
	}
}

func (c *CheckerInfo) UpdateProgress() {
	c.Progress = float64(c.Checked) / float64(c.Total) * 100
}

func (c *CheckerInfo) UpdateCpm() {
	before := c.Checked
	time.Sleep(time.Second * 6)
	c.Cpm = float64(c.Checked-before) / 6 * 60
}

func (c *CheckerInfo) Elapsed() time.Duration {
	return time.Since(c.StartTime)
}

func (c *CheckerInfo) Remaining() time.Duration {
	if c.Cpm == 0 {
		return 0
	}
	return time.Duration(float64(c.Total-c.Checked)/c.Cpm) * time.Minute
}

func (c *CheckerInfo) ETA() time.Time {
	return time.Now().Add(c.Remaining())
}
