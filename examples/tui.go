package main

import (
	"os"

	"github.com/Velka-DEV/Gotane/v2/pkg/tui/views"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := views.NewHomeView()
	tea.NewProgram(&model, tea.WithAltScreen(), tea.WithOutput(os.Stderr)).Run()
}
