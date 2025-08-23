package main

import (
	"time"

	"github.com/ppmpreetham/vesper/cmd/flags"
	"github.com/ppmpreetham/vesper/cmd/output"
	"github.com/ppmpreetham/vesper/cmd/runner"
	"github.com/ppmpreetham/vesper/tools"
)

func main() {
	tools.PrintLogo()

	// Parse command line arguments
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
