package weathervibes

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Spinner struct {
	spinner spinner.Model
}

func NewSpinner() Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Spinner{spinner: s}
}

func (s Spinner) View() string {
	return s.spinner.View()
}

func (s Spinner) Tick() tea.Cmd {
	return s.spinner.Tick
}
