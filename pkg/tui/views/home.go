package views

import (
	"fmt"
	"math"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomeModel struct {
	CheckerName        string
	CheckerVersion     string
	CheckerAuthor      string
	CheckerDescription string
	timePassed         float64
	width              int
	height             int
}

func NewHomeView() HomeModel {
	return HomeModel{
		CheckerName:        "Gotane",
		CheckerVersion:     "v2.0.0",
		CheckerAuthor:      "Velka-DEV",
		CheckerDescription: "A simple and fast checker written in Go",
	}
}

func (m HomeModel) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg {
		return timer.TickMsg{}
	})
}

func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, tea.Quit
		case "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case timer.TickMsg:
		m.timePassed += 0.001
	}
	return m, tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg {
		return timer.TickMsg{}
	})
}

// View renders the application UI.
func (m HomeModel) View() string {
	// RGB wave calculation for the title
	scaledSine := (math.Sin(m.timePassed) + 1) / 2
	colorIndex := int(scaledSine*150) + 50
	color := lipgloss.Color(fmt.Sprintf("%d", colorIndex))
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(color).Align(lipgloss.Center).PaddingTop(2)

	versionColor := lipgloss.Color("201") // aqua
	authorColor := lipgloss.Color("201")  // hot pink
	versionStyle := lipgloss.NewStyle().Foreground(versionColor).Underline(true).Align(lipgloss.Center)
	authorStyle := lipgloss.NewStyle().Foreground(authorColor).Underline(true).Align(lipgloss.Center)

	buttonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Padding(0, 1).
		Align(lipgloss.Center)

	infoStyle := lipgloss.NewStyle().Align(lipgloss.Center).PaddingTop(1)

	// Render content with the styles
	title := titleStyle.Width(m.width).Render(m.CheckerName)
	version := versionStyle.Render("Version") + ": " + m.CheckerVersion
	author := authorStyle.Render("Author") + ": " + m.CheckerAuthor
	info := infoStyle.Width(m.width).Render(version + " | " + author)
	button := buttonStyle.Width(m.width).Render("Press \"ENTER\" to start configuration")

	// Calculate the required padding to push the button to the bottom
	linesUsed := 3                           // 1 for title, 1 for info, 1 for button
	paddingLines := m.height - linesUsed - 4 // 4 is an estimated value considering potential line breaks, adjust as needed

	padding := lipgloss.NewStyle().Height(paddingLines).Render("") // Create a block of padding

	// Join the content vertically
	content := lipgloss.JoinVertical(lipgloss.Top,
		title,
		info,
		padding,
		button,
	)

	return content
}
