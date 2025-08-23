package runner

import (
	"fmt"
	"sync"

	"github.com/ppmpreetham/vesper/pkg/types"
	"github.com/ppmpreetham/vesper/sites"
	"github.com/ppmpreetham/vesper/tools"
)

// RunSherlockDatabase runs the Sherlock database enumeration
func RunSherlockDatabase(username string, config types.Config) int {
	var wg sync.WaitGroup
	jobs := make(chan types.SherlockJob, config.BufferSize)
	results := make(chan tools.ReturnData, config.BufferSize)

	// Start worker pool for Sherlock
	for i := 0; i < config.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				result := tools.SherlockCheckURL(username, job.Data, job.Name)
				results <- result
			}
		}()
	}

	// Send Sherlock jobs
	go func() {
		for siteName, site := range sites.SherlockSites {
			jobs <- types.SherlockJob{
				Name: siteName,
				Data: site,
			}
		}
		close(jobs)
	}()

	// Wait and close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and print results
	foundCount := 0
	for result := range results {
		if result.Status == "FOUND" {
			foundCount++
			tools.Green("Found: ")
			fmt.Print(result.Name, " at ")
			tools.BoldGreen(result.URL)
			fmt.Print("\n")
		}
	}

	return foundCount
}
