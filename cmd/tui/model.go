package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the TUI state for input collection
type Model struct {
	state         int
	usernameInput textinput.Model
	databaseList  list.Model
	err           error
	width, height int
	// Final selections
	selectedUsername string
	selectedDatabase string
	ready            bool
}

const (
	stateUsername = iota
	stateDatabase
	stateReady
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

	return Model{
		state:         stateUsername,
		usernameInput: usernameInput,
		databaseList:  databaseList,
		ready:         false,
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
				m.selectedUsername = m.usernameInput.Value()
				m.state = stateDatabase
				return m, nil

			case stateDatabase:
				// Set up the configuration
				selectedItem, ok := m.databaseList.SelectedItem().(item)
				if !ok {
					m.err = fmt.Errorf("failed to get selected database")
					return m, nil
				}

				m.selectedDatabase = selectedItem.title
				m.ready = true
				m.state = stateReady

				// Exit the TUI loop - we have what we need
				return m, tea.Quit
			}
		}
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
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
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
		view += focusedStyle.Render("Username: ") + blurredStyle.Render(m.selectedUsername) + "\n\n"
		view += m.databaseList.View() + "\n\n"
		view += "↑/↓: Select • Enter: Confirm • Ctrl+C: Quit\n"
		return view

	case stateReady:
		return "Ready to start enumeration..."

	default:
		return "Unknown state"
	}
}

// GetSelections returns the selected username and database
func (m Model) GetSelections() (string, string, bool) {
	return m.selectedUsername, m.selectedDatabase, m.ready
}
