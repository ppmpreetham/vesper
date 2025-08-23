package tui

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/ppmpreetham/vesper/pkg/types"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Bold(true).
			Padding(0, 1)

	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0087BD"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))
)

// item represents a selectable item in the list
type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// runFinishedMsg is sent when the runner is done
type runFinishedMsg struct {
	results []types.JobResult
	elapsed time.Duration
}

// foundResultMsg represents a new found result
type foundResultMsg struct {
	result string
}
