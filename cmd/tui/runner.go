package tui

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ppmpreetham/vesper/pkg/types"
	"github.com/ppmpreetham/vesper/tools"
)

// FoundResultsChan is a channel to communicate found results to the progress model
var FoundResultsChan = make(chan tools.ReturnData, 100) // Buffer the channel to prevent blocking

// Global program reference to send messages
var globalProgram *tea.Program

// RunTUI starts the TUI interface for Vesper
func RunTUI() (types.Config, string, bool) {
	// First program: collect username and database
	p := tea.NewProgram(Initialize())

	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	// Get selections
	model := m.(Model)
	username, dbName, ready := model.GetSelections()

	// If not ready, just exit
	if !ready {
		return types.Config{}, "", true
	}

	// Create the config with selected values
	config := types.Config{
		Database:   dbName,
		Timeout:    7, // Default timeout
		NumWorkers: 250,
		BufferSize: 1000,
	}

	// Set the timeout in tools package
	tools.SetHTTPTimeout(time.Duration(config.Timeout) * time.Second)

	// Set up the found callback to use our channel
	var allFoundResults []tools.ReturnData
	tools.SetFoundCallback(func(name, url string) {
		result := tools.ReturnData{
			Name:   name,
			URL:    url,
			Status: "FOUND",
		}
		allFoundResults = append(allFoundResults, result)

		// IMPORTANT: NO DEBUG PRINTS IN TUI MODE
		// Send to channel, but don't block if channel is full
		select {
		case FoundResultsChan <- result:
			// Sent successfully
		default:
			// Channel full or closed, just continue
		}
	})

	// Second program: Show progress and run enumeration
	progModel := NewProgressModel(username, dbName, config)
	progressProgram := tea.NewProgram(progModel, tea.WithAltScreen()) // Use alt screen for clean display
	globalProgram = progressProgram

	// Run the channel listener in background
	go func() {
		for result := range FoundResultsChan {
			progressProgram.Send(foundMsg{result: result})
		}
	}()

	// Run the progress program
	finalModel, err := progressProgram.Run()
	if err != nil {
		fmt.Printf("Error running progress program: %v\n", err)
		os.Exit(1)
	}

	// Close the channel when done
	close(FoundResultsChan)

	// Get the final results
	finalProgressModel := finalModel.(ProgressModel)
	results := finalProgressModel.results
	foundResults := finalProgressModel.allResults

	// If we have results, show them in a table
	if len(foundResults) > 0 {
		tableModel := NewTableModel(foundResults)
		tableProgram := tea.NewProgram(tableModel, tea.WithAltScreen())

		if _, err := tableProgram.Run(); err != nil {
			fmt.Printf("Error running table program: %v\n", err)
		}
	} else {
		fmt.Println("No profiles found to display in table.")
	}

	// Print summary
	totalFound := 0
	for _, result := range results {
		totalFound += result.FoundCount
	}

	fmt.Printf("\nExecution completed\n")
	fmt.Printf("Found username on %d sites\n", totalFound)

	for _, result := range results {
		fmt.Printf("  - %s: %d matches\n", result.DatabaseName, result.FoundCount)
	}

	return config, username, true
}
