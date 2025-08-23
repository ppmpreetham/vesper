package tui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ppmpreetham/vesper/cmd/runner"
	"github.com/ppmpreetham/vesper/pkg/types"
	"github.com/ppmpreetham/vesper/tools"
)

// Model represents the TUI state
type Model struct {
	state         int
	usernameInput textinput.Model
	databaseList  list.Model
	timeoutInput  textinput.Model
	spinner       spinner.Model
	err           error
	results       []types.JobResult
	elapsedTime   time.Duration
	width, height int
	foundResults  []string
}

const (
	stateUsername = iota
	stateDatabase
	stateTimeout
	stateRunning
	stateDone
)

// Initialize creates a new TUI model
func Initialize() Model {
	// Username input
	usernameInput := textinput.New()
	usernameInput.Placeholder = "Enter username to search"
	usernameInput.Focus()
	usernameInput.CharLimit = 50
	usernameInput.Width = 30

	// Database selection
	databaseItems := []list.Item{
		item{title: "whatsmyname", desc: "Use WhatsMyName database (stable)"},
		item{title: "sherlock", desc: "Use Sherlock database"},
		item{title: "maigret", desc: "Use Maigret database"},
		item{title: "all", desc: "Use all available databases"},
	}

	databaseDelegate := list.NewDefaultDelegate()
	databaseDelegate.ShowDescription = true
	databaseDelegate.SetSpacing(1)

	databaseList := list.New(databaseItems, databaseDelegate, 0, 0)
	databaseList.Title = "Select Database"
	databaseList.SetShowStatusBar(false)
	databaseList.SetFilteringEnabled(false)
	databaseList.SetShowHelp(false)
	databaseList.SetShowPagination(false)

	// Timeout input
	timeoutInput := textinput.New()
	timeoutInput.Placeholder = "HTTP timeout in seconds (default: 7)"
	timeoutInput.CharLimit = 2
	timeoutInput.Width = 30

	// Spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		state:         stateUsername,
		usernameInput: usernameInput,
		databaseList:  databaseList,
		timeoutInput:  timeoutInput,
		spinner:       s,
		foundResults:  make([]string, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.databaseList.SetSize(msg.Width-4, msg.Height/2)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			switch m.state {
			case stateUsername:
				if m.usernameInput.Value() == "" {
					m.err = fmt.Errorf("username cannot be empty")
					return m, nil
				}
				m.err = nil
				m.state = stateDatabase
				return m, nil

			case stateDatabase:
				m.state = stateTimeout
				return m, nil

			case stateTimeout:
				timeout := 7 // default
				if m.timeoutInput.Value() != "" {
					var err error
					timeout, err = strconv.Atoi(m.timeoutInput.Value())
					if err != nil || timeout < 1 || timeout > 60 {
						m.err = fmt.Errorf("timeout must be a number between 1 and 60")
						return m, nil
					}
				}
				m.err = nil
				m.state = stateRunning

				// Set up the configuration
				selectedItem, ok := m.databaseList.SelectedItem().(item)
				if !ok {
					m.err = fmt.Errorf("failed to get selected database")
					m.state = stateDatabase
					return m, nil
				}

				config := types.Config{
					Database:   selectedItem.title,
					Timeout:    timeout,
					NumWorkers: 250,
					BufferSize: 1000,
				}

				// Intercept the found results by setting up a custom handler
				tools.SetFoundCallback(func(name, url string) {
					// Store found results for later display
					m.foundResults = append(m.foundResults, fmt.Sprintf("Found: %s at %s", name, url))
				})

				// Run the enumeration in a goroutine
				return m, tea.Batch(
					m.spinner.Tick,
					func() tea.Msg {
						tools.SetHTTPTimeout(time.Duration(config.Timeout) * time.Second)
						startTime := time.Now()
						results := runner.RunEnumeration(m.usernameInput.Value(), config)
						elapsedTime := time.Since(startTime)
						return runFinishedMsg{results: results, elapsed: elapsedTime}
					},
				)
			}
		}

	case runFinishedMsg:
		m.results = msg.results
		m.elapsedTime = msg.elapsed
		m.state = stateDone
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case foundResultMsg:
		// Add new found result to our list
		m.foundResults = append(m.foundResults, msg.result)
		return m, nil
	}

	// Handle input updates
	switch m.state {
	case stateUsername:
		var cmd tea.Cmd
		m.usernameInput, cmd = m.usernameInput.Update(msg)
		cmds = append(cmds, cmd)

	case stateDatabase:
		var cmd tea.Cmd
		m.databaseList, cmd = m.databaseList.Update(msg)
		cmds = append(cmds, cmd)

	case stateTimeout:
		var cmd tea.Cmd
		m.timeoutInput, cmd = m.timeoutInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	// Don't call PrintLogo here to avoid clearing the screen on every render

	switch m.state {
	case stateUsername:
		view := "\n" + titleStyle.Render(" VESPER OSINT ") + "\n\n"
		view += focusedStyle.Render("Username: ") + "\n"
		view += m.usernameInput.View() + "\n\n"
		view += "Press Enter to continue, Ctrl+C to quit\n"

		if m.err != nil {
			view += "\n" + errorStyle.Render(m.err.Error())
		}
		return view

	case stateDatabase:
		view := "\n" + titleStyle.Render(" VESPER OSINT ") + "\n\n"
		view += focusedStyle.Render("Username: ") + blurredStyle.Render(m.usernameInput.Value()) + "\n\n"
		view += m.databaseList.View() + "\n\n"
		view += "↑/↓: Select • Enter: Confirm • Ctrl+C: Quit\n"
		return view

	case stateTimeout:
		view := "\n" + titleStyle.Render(" VESPER OSINT ") + "\n\n"
		view += blurredStyle.Render("Username: ") + blurredStyle.Render(m.usernameInput.Value()) + "\n\n"

		selectedItem, ok := m.databaseList.SelectedItem().(item)
		dbName := "whatsmyname"
		if ok {
			dbName = selectedItem.title
		}
		view += blurredStyle.Render("Database: ") + blurredStyle.Render(dbName) + "\n\n"

		view += focusedStyle.Render("Timeout (seconds): ") + "\n"
		view += m.timeoutInput.View() + "\n\n"
		view += "Press Enter to run, Ctrl+C to quit\n"

		if m.err != nil {
			view += "\n" + errorStyle.Render(m.err.Error())
		}
		return view

	case stateRunning:
		view := "\n" + titleStyle.Render(" VESPER OSINT ") + "\n\n"
		view += infoStyle.Render("Running enumeration for ") +
			focusedStyle.Render(m.usernameInput.Value()) + "...\n\n"

		// Display spinner
		view += m.spinner.View() + " Scanning...\n\n"

		// Display last 10 found results (to avoid flooding the screen)
		if len(m.foundResults) > 0 {
			view += successStyle.Render("Latest findings:") + "\n"

			// Get the last 10 results or less if there aren't 10 yet
			startIdx := 0
			if len(m.foundResults) > 10 {
				startIdx = len(m.foundResults) - 10
			}

			for i := startIdx; i < len(m.foundResults); i++ {
				view += "  " + m.foundResults[i] + "\n"
			}

			// Show count if there are more than shown
			if len(m.foundResults) > 10 {
				view += fmt.Sprintf("\n%s\n", infoStyle.Render(
					fmt.Sprintf("...and %d more (full results will be shown at the end)",
						len(m.foundResults)-10)))
			}
		}

		return view

	case stateDone:
		// This won't actually show in the TUI since we quit after finishing
		return "Done!"

	default:
		return "Unknown state"
	}
}

// GetResults returns the final results from the TUI
func (m Model) GetResults() ([]types.JobResult, time.Duration, bool, []string) {
	return m.results, m.elapsedTime, m.state == stateDone, m.foundResults
}
