package main

import (
	"flag"
	"time"

	"github.com/ppmpreetham/vesper/cmd/flags"
	"github.com/ppmpreetham/vesper/cmd/output"
	"github.com/ppmpreetham/vesper/cmd/runner"
	"github.com/ppmpreetham/vesper/cmd/tui"
	"github.com/ppmpreetham/vesper/tools"
)

func main() {
	// Print the logo once at startup
	tools.PrintLogo()

	// Check if we should use TUI or CLI
	useTUIPtr := flag.Bool("tui", false, "Use interactive TUI interface")
	flag.Parse()
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError) // Reset flags for subsequent parsing

	// If --tui flag is provided, use the TUI interface
	if *useTUIPtr {
		// TUI mode will handle everything internally
		_, _, _ = tui.RunTUI()
		return
	}

	// Otherwise, use the traditional CLI interface
	config, username, shouldExit := flags.Parse()
	if shouldExit {
		return
	}

	// Set the timeout in tools package
	tools.SetHTTPTimeout(time.Duration(config.Timeout) * time.Second)

	// Start the execution
	startTime := time.Now()
	results := runner.RunEnumeration(username, config)
	elapsedTime := time.Since(startTime)

	// Print summary
	output.PrintSummary(results, elapsedTime)
}
