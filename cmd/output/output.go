package output

import (
	"fmt"
	"time"

	"github.com/ppmpreetham/vesper/pkg/types"
)

// PrintSummary prints a summary of the results
func PrintSummary(results []types.JobResult, elapsedTime time.Duration) {
	if len(results) == 0 {
		return
	}

	totalFoundCount := 0
	for _, result := range results {
		totalFoundCount += result.FoundCount
	}

	if len(results) > 1 {
		fmt.Printf("\n=== Summary ===\n")
		fmt.Printf("Execution completed in %s\n", elapsedTime)
		fmt.Printf("Total found username on %d sites across all databases\n", totalFoundCount)

		for _, result := range results {
			fmt.Printf("  - %s: %d matches\n", result.DatabaseName, result.FoundCount)
		}
	} else {
		fmt.Printf("\nExecution completed in %s\n", elapsedTime)
		fmt.Printf("Found username on %d sites\n", totalFoundCount)
	}
}
