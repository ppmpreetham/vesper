package runner

import (
	"fmt"
	"sync"

	"github.com/ppmpreetham/vesper/pkg/types"
	"github.com/ppmpreetham/vesper/sites"
	"github.com/ppmpreetham/vesper/tools"
)

// RunMaigretDatabase runs the Maigret database enumeration
func RunMaigretDatabase(username string, config types.Config) int {
	var wg sync.WaitGroup
	jobs := make(chan types.MaigretJob, config.BufferSize)
	results := make(chan tools.ReturnData, config.BufferSize)
	sitesChecked := 0
	var sitesCheckedMutex sync.Mutex

	// Start worker pool for Maigret
	for i := 0; i < config.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				result := tools.MaigretCheckURL(username, job.Data, job.Name)
				results <- result

				// Update sites checked counter
				sitesCheckedMutex.Lock()
				sitesChecked++
				currentCount := sitesChecked
				sitesCheckedMutex.Unlock()

				// Report progress
				tools.NotifyProgress(currentCount)
			}
		}()
	}

	// Send Maigret jobs
	go func() {
		for siteName, site := range sites.MaigretSites {
			jobs <- types.MaigretJob{
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
			handled := tools.NotifyFound(result.Name, result.URL)

			if !handled {
				tools.Green("Found: ")
				fmt.Print(result.Name, " at ")
				tools.BoldGreen(result.URL)
				fmt.Print("\n")
			}
		}
	}

	return foundCount
}
