package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ppmpreetham/vesper/cmd/output"
	"github.com/ppmpreetham/vesper/pkg/types"
)

// RunTUI starts the TUI interface for Vesper
func RunTUI() (types.Config, string, bool) {
	// Create and run the program
	p := tea.NewProgram(Initialize(), tea.WithAltScreen())

	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	// Get the final state
	model := m.(Model)
	results, elapsed, successful, foundItems := model.GetResults()

	// Print results if we have them
	if successful {
		// First print all found items clearly (no overwriting)
		fmt.Println("\n=== Found Profiles ===")
		for _, item := range foundItems {
			fmt.Println(item)
		}
		fmt.Println()

		// Then print the summary
		output.PrintSummary(results, elapsed)
	}

	// Return empty values and true to indicate we've handled everything
	return types.Config{}, "", true
}
