package core

import "time"

type TransportConfig struct {
	DisableKeepAlive    bool
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
}

type CheckerBuilder struct {
	config          *CheckerConfig
	transportConfig *TransportConfig
	queue           *ComboQueue
	pool            *HTTPClientPool

	checkProcess  CheckProcess
	outputProcess OutputProcess
}

func NewCheckerBuilder() *CheckerBuilder {
	return &CheckerBuilder{
		transportConfig: &TransportConfig{
			DisableKeepAlive:    false,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 2,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}

func (b *CheckerBuilder) WithConfig(config *CheckerConfig) *CheckerBuilder {
	b.config = config
	return b
}

func (b *CheckerBuilder) WithTransportConfig(transportConfig *TransportConfig) *CheckerBuilder {
	b.transportConfig = transportConfig
	return b
}

func (b *CheckerBuilder) WithCombos(combos []*Combo) *CheckerBuilder {

	if len(combos) == 0 {
		panic("combos are empty")
	}

	b.queue = NewComboQueue(combos)
	return b
}

func (b *CheckerBuilder) WithProxies(proxies []string, scheme string) *CheckerBuilder {
	b.pool = NewHTTPClientPoolWithProxies(proxies, scheme, b.transportConfig)
	return b
}

func (b *CheckerBuilder) WithCheckProcess(checkProcess CheckProcess) *CheckerBuilder {
	b.checkProcess = checkProcess
	return b
}

func (b *CheckerBuilder) WithOutputProcess(outputProcess OutputProcess) *CheckerBuilder {
	b.outputProcess = outputProcess
	return b
}

func (b *CheckerBuilder) Build() *Checker {
	if b.config == nil {
		panic("checker config is nil")
	}

	if b.queue == nil {
		panic("combos are nil or empty")
	}

	if b.checkProcess == nil {
		panic("check process is nil")
	}

	if b.pool == nil {
		b.pool = NewHTTPClientPool(b.transportConfig)
	}

	return NewChecker(
		b.config,
		b.queue,
		b.pool,
		b.checkProcess,
		b.outputProcess,
	)
}
