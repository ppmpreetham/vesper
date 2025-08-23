package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ppmpreetham/vesper/cmd/runner"
	"github.com/ppmpreetham/vesper/pkg/types"
	"github.com/ppmpreetham/vesper/sites"
	"github.com/ppmpreetham/vesper/tools"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

// ProgressModel represents the progress tracking state
type ProgressModel struct {
	progress     progress.Model
	width        int
	currentStep  int
	totalSteps   int
	username     string
	config       types.Config
	dbName       string
	results      []types.JobResult
	allResults   []tools.ReturnData
	done         bool
	finalMessage string
	sitesChecked int
	totalSites   int
}

// Custom message types
type tickMsg time.Time
type doneMsg struct {
	results []types.JobResult
}
type foundMsg struct {
	result tools.ReturnData
}
type progressSiteMsg struct {
	sitesChecked int
	totalSites   int
}

// NewProgressModel creates a new progress tracking model
func NewProgressModel(username, dbName string, config types.Config) ProgressModel {
	prog := progress.New(progress.WithDefaultGradient())

	// Set total steps and calculate total sites based on database
	totalSteps := 1
	totalSites := 0

	if dbName == "all" {
		totalSteps = 3 // Sherlock, WhatsMyName, and Maigret
		totalSites = len(sites.SherlockSites) + len(sites.WhatsmynameSites) + len(sites.MaigretSites)
	} else {
		switch dbName {
		case "sherlock":
			totalSites = len(sites.SherlockSites)
		case "maigret":
			totalSites = len(sites.MaigretSites)
		default: // whatsmyname or empty
			totalSites = len(sites.WhatsmynameSites)
		}
	}

	return ProgressModel{
		progress:     prog,
		username:     username,
		dbName:       dbName,
		config:       config,
		currentStep:  0,
		totalSteps:   totalSteps,
		allResults:   make([]tools.ReturnData, 0),
		sitesChecked: 0,
		totalSites:   totalSites,
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		runEnumeration(m.username, m.config, m.dbName, m.currentStep),
	)
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.done {
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.done {
			return m, nil
		}
		return m, tickCmd()

	case progressSiteMsg:
		// Update the sites checked count
		m.sitesChecked = msg.sitesChecked

		// Calculate progress percentage
		var percent float64
		if m.totalSites > 0 {
			percent = float64(m.sitesChecked) / float64(m.totalSites)
		}

		// Ensure we don't go over 100%
		if percent > 1.0 {
			percent = 1.0
		}

		cmd := m.progress.SetPercent(percent)
		return m, cmd

	case foundMsg:
		// Add the found result to our collection
		m.allResults = append(m.allResults, msg.result)
		return m, nil

	case doneMsg:
		// One database search is done
		m.currentStep++
		m.results = append(m.results, msg.results...)

		if m.currentStep >= m.totalSteps {
			m.done = true
			// Force progress to 100%
			cmd := m.progress.SetPercent(1.0)

			// Once all done, show the final message
			if len(m.allResults) > 0 {
				m.finalMessage = fmt.Sprintf("Found %d profiles. Press any key to view results.", len(m.allResults))
			} else {
				m.finalMessage = "No profiles found. Press any key to exit."
			}
			return m, cmd
		}

		// If not all done, start the next database search
		return m, runEnumeration(m.username, m.config, m.dbName, m.currentStep)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m ProgressModel) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	pad := strings.Repeat(" ", padding)

	// Show current database being searched
	var currentDB string
	if m.dbName == "all" {
		switch m.currentStep {
		case 0:
			currentDB = "Searching Sherlock database..."
		case 1:
			currentDB = "Searching WhatsMyName database..."
		case 2:
			currentDB = "Searching Maigret database..."
		}
	} else {
		currentDB = fmt.Sprintf("Searching %s database...", m.dbName)
	}

	// Calculate progress metrics
	sitesCheckedStr := fmt.Sprintf("%d / %d sites checked", m.sitesChecked, m.totalSites)

	view := "\n" +
		titleStyle.Render(" VESPER OSINT ") + "\n\n" +
		pad + "Searching for: " + focusedStyle.Render(m.username) + "\n" +
		pad + currentDB + "\n\n" +
		pad + m.progress.View() + "\n" +
		pad + sitesCheckedStr + "\n\n"

	// Show counts of found items so far
	if len(m.allResults) > 0 {
		view += pad + fmt.Sprintf("Found %d profiles so far\n\n", len(m.allResults))
	}

	if m.done {
		view += pad + successStyle.Render(m.finalMessage) + "\n"
	} else {
		view += pad + helpStyle("Press q to cancel") + "\n"
	}

	return view
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// runEnumeration runs the enumeration in a goroutine and returns when done
func runEnumeration(username string, config types.Config, dbName string, step int) tea.Cmd {
	return func() tea.Msg {
		var results []types.JobResult

		// Set up the site progress reporter
		var sitesInCurrentDB int
		var sitesOffset int

		if dbName == "all" {
			switch step {
			case 0: // Sherlock
				sitesInCurrentDB = len(sites.SherlockSites)
				sitesOffset = 0
			case 1: // WhatsMyName
				sitesInCurrentDB = len(sites.WhatsmynameSites)
				sitesOffset = len(sites.SherlockSites)
			case 2: // Maigret
				sitesInCurrentDB = len(sites.MaigretSites)
				sitesOffset = len(sites.SherlockSites) + len(sites.WhatsmynameSites)
			}
		} else {
			switch dbName {
			case "sherlock":
				sitesInCurrentDB = len(sites.SherlockSites)
			case "maigret":
				sitesInCurrentDB = len(sites.MaigretSites)
			default: // whatsmyname or empty
				sitesInCurrentDB = len(sites.WhatsmynameSites)
			}
			sitesOffset = 0
		}

		// Set up a progress reporter
		tools.SetProgressCallback(func(sitesChecked int) {
			if globalProgram != nil {
				globalProgram.Send(progressSiteMsg{
					sitesChecked: sitesOffset + sitesChecked,
					totalSites:   sitesInCurrentDB,
				})
			}
		})

		// Run the actual database search
		if dbName == "all" {
			// For "all" mode, run each database sequentially based on the step
			switch step {
			case 0: // Sherlock
				count := runner.RunSherlockDatabase(username, config)
				results = []types.JobResult{{DatabaseName: "Sherlock", FoundCount: count}}
			case 1: // WhatsMyName
				tools.ResetHTTPClient()
				count := runner.RunWhatsMyNameDatabase(username, config)
				results = []types.JobResult{{DatabaseName: "WhatsMyName", FoundCount: count}}
			case 2: // Maigret
				tools.ResetHTTPClient()
				count := runner.RunMaigretDatabase(username, config)
				results = []types.JobResult{{DatabaseName: "Maigret", FoundCount: count}}
			}
		} else {
			// For single database mode
			switch dbName {
			case "sherlock":
				count := runner.RunSherlockDatabase(username, config)
				results = []types.JobResult{{DatabaseName: "Sherlock", FoundCount: count}}
			case "maigret":
				count := runner.RunMaigretDatabase(username, config)
				results = []types.JobResult{{DatabaseName: "Maigret", FoundCount: count}}
			default: // whatsmyname or empty
				count := runner.RunWhatsMyNameDatabase(username, config)
				results = []types.JobResult{{DatabaseName: "WhatsMyName", FoundCount: count}}
			}
		}

		// Final progress update (ensure 100% for this step)
		if globalProgram != nil {
			globalProgram.Send(progressSiteMsg{
				sitesChecked: sitesOffset + sitesInCurrentDB,
				totalSites:   sitesInCurrentDB,
			})
		}

		return doneMsg{results: results}
	}
}
