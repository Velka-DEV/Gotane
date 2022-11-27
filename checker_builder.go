package gotane

import (
	"errors"
)

type CheckerBuilder struct {
	Options       *CheckerOptions `json:"options"`
	Proxies       []*Proxy        `json:"proxies"`
	Combos        []*Combo        `json:"combos"`
	CheckProcess  CheckProcess
	OutputProcess OutputProcess
	LogProcess    LogProcess
}

func (cb *CheckerBuilder) SetThreads(threads int) *CheckerBuilder {
	cb.Options.Threads = threads
	return cb
}

func (cb *CheckerBuilder) SetCheckProcess(process CheckProcess) *CheckerBuilder {
	cb.CheckProcess = process
	return cb
}

func (cb *CheckerBuilder) SetOutputProcess(outputProcess OutputProcess) *CheckerBuilder {
	cb.OutputProcess = outputProcess
	return cb
}

func (cb *CheckerBuilder) SetLogProcess(logProcess LogProcess) *CheckerBuilder {
	cb.LogProcess = logProcess
	return cb
}

func (cb *CheckerBuilder) WithProxy(proxy *Proxy) *CheckerBuilder {
	cb.Proxies = append(cb.Proxies, proxy)
	return cb
}

func (cb *CheckerBuilder) WithCombo(combo *Combo) *CheckerBuilder {
	cb.Combos = append(cb.Combos, combo)
	return cb
}

func (cb *CheckerBuilder) WithProxies(proxies []*Proxy) *CheckerBuilder {
	cb.Proxies = append(cb.Proxies, proxies...)
	return cb
}

func (cb *CheckerBuilder) WithCombos(combos []*Combo) *CheckerBuilder {
	cb.Combos = append(cb.Combos, combos...)
	return cb
}

func (cb *CheckerBuilder) SetMaxRetries(maxRetries int) *CheckerBuilder {
	cb.Options.MaxRetries = maxRetries
	return cb
}

func (cb *CheckerBuilder) WithProxiesFromUrlLines(lines []string) *CheckerBuilder {
	for _, line := range lines {
		proxy, err := NewProxyFromUrlString(line)
		if err != nil {
			continue
		}
		cb.WithProxy(proxy)
	}
	return cb
}

func (cb *CheckerBuilder) WithProxiesFromLines(lines []string, protocol string) *CheckerBuilder {
	for _, line := range lines {
		proxy, err := NewProxyFromString(line, protocol)
		if err != nil {
			continue
		}
		cb.WithProxy(proxy)
	}
	return cb
}

func (cb *CheckerBuilder) WithCombosFromLines(lines []string) *CheckerBuilder {
	for _, line := range lines {
		combo, err := NewComboFromString(line)
		if err != nil {
			continue
		}
		cb.WithCombo(combo)
	}
	return cb
}

func (cb *CheckerBuilder) Build() (*Checker, error) {

	if cb.Options.Threads == 0 {
		return nil, errors.New("threads must be greeter than zero")
	}

	if cb.Options.MaxRetries == 0 {
		return nil, errors.New("max retries value must be set greater than zero")
	}

	if len(cb.Combos) == 0 {
		return nil, errors.New("combos list cannot be empty")
	}

	if cb.CheckProcess == nil {
		return nil, errors.New("check process cannot be nil")
	}

	return newChecker(cb.Options, cb.Proxies, cb.Combos, cb.CheckProcess, cb.OutputProcess, cb.LogProcess), nil
}

func NewCheckerBuilder() *CheckerBuilder {
	return &CheckerBuilder{
		Options: &CheckerOptions{
			Threads:    10,
			MaxRetries: 3,
		},
	}
}
