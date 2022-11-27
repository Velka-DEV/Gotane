package gotane

import (
	"sync/atomic"
	"time"
)

type CheckerInfo struct {
	checked   uint64    `json:"checked"`
	invalid   uint64    `json:"invalid"`
	free      uint64    `json:"free"`
	hit       uint64    `json:"hit"`
	retry     uint64    `json:"retry"`
	error     uint64    `json:"error"`
	locked    uint64    `json:"locked"`
	banned    uint64    `json:"banned"`
	total     uint64    `json:"total"`
	cpm       int       `json:"cpm"`
	startDate time.Time `json:"start_date"`
}

func (c *CheckerInfo) IncrementChecked() {
	atomic.AddUint64(&c.checked, 1)
}

func (c *CheckerInfo) IncrementInvalid() {
	atomic.AddUint64(&c.invalid, 1)
}

func (c *CheckerInfo) IncrementFree() {
	atomic.AddUint64(&c.free, 1)
}

func (c *CheckerInfo) IncrementHit() {
	atomic.AddUint64(&c.hit, 1)
}

func (c *CheckerInfo) IncrementRetry() {
	atomic.AddUint64(&c.retry, 1)
}

func (c *CheckerInfo) IncrementLocked() {
	atomic.AddUint64(&c.locked, 1)
}

func (c *CheckerInfo) IncrementBanned() {
	atomic.AddUint64(&c.banned, 1)
}

func (c *CheckerInfo) IncrementError() {
	atomic.AddUint64(&c.error, 1)
}

func (c *CheckerInfo) GetTotal() uint64 {
	return atomic.LoadUint64(&c.total)
}

func (c *CheckerInfo) GetChecked() uint64 {
	return atomic.LoadUint64(&c.checked)
}

func (c *CheckerInfo) GetInvalid() uint64 {
	return atomic.LoadUint64(&c.invalid)
}

func (c *CheckerInfo) GetFree() uint64 {
	return atomic.LoadUint64(&c.free)
}

func (c *CheckerInfo) GetHit() uint64 {
	return atomic.LoadUint64(&c.hit)
}

func (c *CheckerInfo) GetRetry() uint64 {
	return atomic.LoadUint64(&c.retry)
}

func (c *CheckerInfo) GetLocked() uint64 {
	return atomic.LoadUint64(&c.locked)
}

func (c *CheckerInfo) GetBanned() uint64 {
	return atomic.LoadUint64(&c.banned)
}

func (c *CheckerInfo) GetError() uint64 {
	return atomic.LoadUint64(&c.error)
}

func (c *CheckerInfo) SetTotal(total uint64) {
	atomic.StoreUint64(&c.total, total)
}

func (c *CheckerInfo) GetCpm() int {
	return c.cpm
}

func (c *CheckerInfo) SetCpm(cpm int) {
	c.cpm = cpm
}

func (c *CheckerInfo) GetStartDate() time.Time {
	return c.startDate
}

func (c *CheckerInfo) GetElapsedTime() time.Duration {
	return time.Since(c.startDate)
}

func (c *CheckerInfo) StartUpdater() {

	c.startDate = time.Now()

	go func() {
		for {
			before := c.GetChecked()

			time.Sleep(time.Second * 6)

			c.SetCpm(int(c.GetChecked()-before) * 10)
		}
	}()
}

func NewCheckerInfo(total uint64) *CheckerInfo {
	return &CheckerInfo{
		checked: 0,
		invalid: 0,
		free:    0,
		hit:     0,
		retry:   0,
		locked:  0,
		banned:  0,
		error:   0,
		total:   total,
		cpm:     0,
	}
}
