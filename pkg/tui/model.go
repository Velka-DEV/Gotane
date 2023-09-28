package tui

import (
	"github.com/Velka-DEV/Gotane/v2/pkg/core"
	"github.com/Velka-DEV/Gotane/v2/pkg/tui/views"
)

type RunningModel struct {
	Checker *core.Checker
	Info    *core.CheckerInfo
	Config  *core.CheckerConfig
}

type ConsoleUiStage int

const (
	HomeStage ConsoleUiStage = iota
	ConfigurationStage
	RunningStage
)

type ConsoleUiModel struct {
	Stage   ConsoleUiStage
	Home    *views.HomeModel
	Config  *views.ConfigurationModel
	Running *RunningModel
}
