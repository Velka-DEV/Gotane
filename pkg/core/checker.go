package core

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
)

type CheckerConfig struct {
	Threads         int
	WaitBeforeRetry time.Duration
	OutputToFile    bool
	OutputDirectory string
	OutputFree      bool
	OutputInvalid   bool
	OutputLocked    bool
	OutputUnknown   bool
}

type CheckProcessArgs struct {
	Combo   *Combo
	Client  *ClientWithProxy
	Retries uint16
}

type CheckProcess func(args *CheckProcessArgs) *CheckResult
type OutputProcess func(args *CheckProcessArgs, result *CheckResult)

type Checker struct {
	Config     *CheckerConfig
	Infos      *CheckerInfo
	queue      *ComboQueue
	clientPool *HTTPClientPool
	waitGroup  *sync.WaitGroup
	workerPool *ants.PoolWithFunc

	checkProcess  CheckProcess
	outputProcess OutputProcess
}

// internalOutputProcess is the internal output process, which will be called after the check process
func (c *Checker) internalOutputProcess(args *CheckProcessArgs, result *CheckResult) {
	if c.outputProcess != nil {
		c.outputProcess(args, result)
	}

	if !c.Config.OutputToFile {
		return
	}

	shouldOutput := false

	switch result.Status {
	case CheckStatusFree:
		shouldOutput = c.Config.OutputFree
	case CheckStatusInvalid:
		shouldOutput = c.Config.OutputInvalid
	case CheckStatusLocked:
		shouldOutput = c.Config.OutputLocked
	case CheckStatusUnknown:
		shouldOutput = c.Config.OutputUnknown
	case CheckStatusError:
		shouldOutput = false
	default:
		shouldOutput = true
	}

	var outputBasePath string

	if c.Config.OutputDirectory == "" {
		outputBasePath, _ = os.Getwd()
	} else {
		outputBasePath = c.Config.OutputDirectory
	}

	if shouldOutput {
		err := WriteResultToFile(result, c.Infos, outputBasePath)
		if err != nil {
			// Display error
			fmt.Println(err)
		}
	}
}

// internalCheckProcess is the internal check process, which will be called recursively if the status is CheckStatusRetry
func (c *Checker) internalCheckProcess(args *CheckProcessArgs) *CheckResult {
	result := c.checkProcess(args)

	switch result.Status {
	case CheckStatusHit:
		atomic.AddUint64(&c.Infos.Hit, 1)
	case CheckStatusInvalid:
		atomic.AddUint64(&c.Infos.Invalid, 1)
	case CheckStatusRetry:
		newClient := c.clientPool.Get()
		args.Client = newClient
		time.Sleep(c.Config.WaitBeforeRetry)
		return c.internalCheckProcess(args)
	case CheckStatusFree:
		atomic.AddUint64(&c.Infos.Free, 1)
	case CheckStatusLocked:
		atomic.AddUint64(&c.Infos.Locked, 1)
	case CheckStatusError:
		atomic.AddUint64(&c.Infos.Error, 1)
	}

	atomic.AddUint64(&c.Infos.Checked, 1)

	c.internalOutputProcess(args, result)

	return result
}

// Start starts the checker
func (c *Checker) Start() {
	pool, err := ants.NewPoolWithFunc(c.Config.Threads, func(args interface{}) {
		for c.Infos.State == CheckerStatePaused {
			time.Sleep(100 * time.Millisecond)
		}

		checkArgs := args.(*CheckProcessArgs)
		c.internalCheckProcess(checkArgs)
		c.waitGroup.Done()
	})

	if err != nil {
		panic(err)
	}

	c.workerPool = pool
	go c.Infos.StartUpdateLoop()

	defer c.workerPool.Release()

	for i := 0; i < c.queue.size; i++ {
		combo, ok := c.queue.Dequeue()

		if !ok {
			panic("could not dequeue combo (this should not appen)")
		}

		client := c.clientPool.Get()
		c.waitGroup.Add(1)
		err := c.workerPool.Invoke(&CheckProcessArgs{Combo: combo, Client: client})
		if err != nil {
			fmt.Println(err)
		}
	}

	go func() {
		c.waitGroup.Wait()
		c.Stop()
	}()
}

// StartAndWait starts the checker and waits for it to finish
func (c *Checker) StartAndWait() {
	c.Start()
	c.waitGroup.Wait()
	c.Stop()
}

// Stop stops the checker
func (c *Checker) Stop() {
	c.Infos.State = CheckerStateStopped

	// Release the worker pool
	if c.workerPool != nil {
		c.workerPool.Release()
	}

	// Wait for all running tasks to complete
	c.waitGroup.Wait()
}

// IsRunning returns true if the checker is running
func (c *Checker) IsRunning() bool {
	return c.Infos.State == CheckerStateRunning
}

// Pause pauses the checker
// It will not stop the current checks, but will prevent new ones from starting
func (c *Checker) Pause() {
	c.Infos.State = CheckerStatePaused
}

// Resume resumes the checker
func (c *Checker) Resume() {
	c.Infos.State = CheckerStateRunning
}

// SetThreads sets the number of threads, and adjusts the worker pool accordingly
func (c *Checker) SetThreads(threads int) {
	c.Config.Threads = threads
	c.workerPool.Tune(threads)
}

// NewChecker creates a new checker
func NewChecker(
	config *CheckerConfig,
	queue *ComboQueue,
	pool *HTTPClientPool,
	checkProcess CheckProcess,
	outputProcess OutputProcess,
) *Checker {
	return &Checker{
		Config:        config,
		Infos:         NewCheckerInfo(uint64(queue.size)),
		queue:         queue,
		clientPool:    pool,
		checkProcess:  checkProcess,
		outputProcess: outputProcess,
		waitGroup:     &sync.WaitGroup{},
	}
}
