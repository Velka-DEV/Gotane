package main

import (
	"errors"
	"github.com/panjf2000/ants/v2"
	"net/http"
	"net/url"
	"sync"
)

type CheckProcess func(args *CheckProcessArgs) CheckResult
type LogProcess func(args *CheckedComboArgs)
type OutputProcess func(args *CheckedComboArgs)

type Checker struct {
	Options     *CheckerOptions `json:"options"`
	Infos       *CheckerInfo    `json:"infos"`
	clients     []*http.Client
	combos      []*Combo
	state       CheckerState
	clientIndex int

	logProcess    LogProcess
	checkProcess  CheckProcess
	outputProcess OutputProcess
	wg            *sync.WaitGroup
	workersPool   *ants.PoolWithFunc
}

func (checker *Checker) internalCheckProcess(args *CheckProcessArgs) CheckResult {

	if checker.checkProcess == nil {
		panic("check process is not set")
	}

	result := checker.checkProcess(args)

	checkedComboArgs := &CheckedComboArgs{
		Combo:  args.Combo,
		Result: result,
	}

	if checker.logProcess != nil {
		checker.logProcess(checkedComboArgs)
	}

	if checker.outputProcess != nil {
		checker.outputProcess(checkedComboArgs)
	}

	switch result {
	case CheckResultInvalid:
		checker.Infos.IncrementInvalid()
	case CheckResultFree:
		checker.Infos.IncrementFree()
	case CheckResultHit:
		checker.Infos.IncrementHit()
	case CheckResultLocked:
		checker.Infos.IncrementLocked()
	case CheckResultError:
		checker.Infos.IncrementError()
	}

	if result != CheckResultError && result != CheckResultRetry {
		checker.Infos.IncrementChecked()
	}

	return result
}

/*
func OutputProcess(combo *Combo, checkResult CheckResult) {
	outputText := []byte(combo.Raw)

	currentPath, err := os.Getwd()

	outputPath := os.Getwd() + "/output.txt"

	switch checkResult {
	case CheckResultInvalid:
	case CheckResultFree:
	case CheckResultHit:
	case CheckResultLocked:
	}
}
*/

func (checker *Checker) getNextClient() *http.Client {

	if len(checker.clients) == 1 {
		return checker.clients[0]
	}

	client := checker.clients[checker.clientIndex]

	if (checker.clientIndex + 1) >= len(checker.clients) {
		checker.clientIndex = 0
	} else {
		checker.clientIndex++
	}

	return client
}

func (checker *Checker) WaitUntilEnd() {
	checker.wg.Wait()
	checker.state = CheckerStateEnded
}

func (checker *Checker) Start() (bool, error) {
	if checker.state != CheckerStateIdle {
		return false, errors.New("checker is not in idle state")
	}

	pool, err := ants.NewPoolWithFunc(checker.Options.Threads, func(i interface{}) {
		args := i.(*CheckProcessArgs)

		checker.internalCheckProcess(args)
		checker.wg.Done()
	})

	if err != nil {
		return false, err
	}

	checker.Infos.StartUpdater()
	checker.workersPool = pool
	checker.state = CheckerStateRunning

	if len(checker.clients) == 0 {
		panic("No available clients")
	}

	for _, combo := range checker.combos {

		args := &CheckProcessArgs{
			Combo:  combo,
			Client: checker.getNextClient(),
		}

		checker.wg.Add(1)
		_ = checker.workersPool.Invoke(args)
	}

	go checker.WaitUntilEnd()

	return true, nil
}

func (checker *Checker) Stop() (bool, error) {
	if checker.state != CheckerStateRunning {
		return false, errors.New("checker is not running")
	}

	checker.workersPool.Release()
	checker.state = CheckerStateEnded

	return true, nil
}

func (checker *Checker) Pause() (bool, error) {
	if checker.state != CheckerStateRunning {
		return false, errors.New("checker is not running")
	}

	//TODO: Check if there is not a better way to pause the pool
	checker.workersPool.Tune(0)
	checker.state = CheckerStatePaused

	return true, nil
}

func (checker *Checker) Resume() (bool, error) {
	if checker.state != CheckerStatePaused {
		return false, errors.New("checker is not paused")
	}

	//TODO: Check if there is not a better way to resume the pool
	checker.workersPool.Tune(checker.Options.Threads)
	checker.state = CheckerStatePaused

	return true, nil
}

func (checker *Checker) SetThreads(threads int) {
	checker.Options.Threads = threads
	checker.workersPool.Tune(threads)
}

func newChecker(options *CheckerOptions, proxies []*Proxy, combos []*Combo, checkProcess CheckProcess, outputProcess OutputProcess, logProcess LogProcess) *Checker {

	clients := make([]*http.Client, len(proxies))

	if len(proxies) > 0 {
		for i, proxy := range proxies {

			proxyUrl, err := url.Parse(proxy.ToString())

			if err != nil {
				continue
			}

			clients[i] = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyUrl),
				},
			}
		}
	} else {
		clients = append(clients, &http.Client{})
	}

	return &Checker{
		Options:       options,
		Infos:         NewCheckerInfo(uint64(len(combos))),
		combos:        combos,
		state:         CheckerStateIdle,
		checkProcess:  checkProcess,
		outputProcess: outputProcess,
		logProcess:    logProcess,
		clients:       clients,
		wg:            &sync.WaitGroup{},
	}
}
