package main

type CheckerState uint

const (
	// CheckerStateIdle is the state of the checker when it is not running
	CheckerStateIdle = iota

	// CheckerStateRunning is the state of the checker when is processing combos
	CheckerStateRunning

	// CheckerStatePaused is the state when the checker as been paused and waiting for user start
	CheckerStatePaused

	// CheckerStateEnded is the state when the checker has processed all the lines of the combo list
	CheckerStateEnded
)
