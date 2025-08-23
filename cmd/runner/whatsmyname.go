package runner

import (
	"fmt"
	"sync"

	"github.com/ppmpreetham/vesper/pkg/types"
	"github.com/ppmpreetham/vesper/sites"
	"github.com/ppmpreetham/vesper/tools"
)

// RunWhatsMyNameDatabase runs the WhatsMyName database enumeration
func RunWhatsMyNameDatabase(username string, config types.Config) int {
	var wg sync.WaitGroup
	jobs := make(chan sites.WhatsmynameSiteData, config.BufferSize)
	results := make(chan tools.ReturnData, config.BufferSize)

	// Start worker pool
	for i := 0; i < config.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for site := range jobs {
				result := tools.WhatsMyNameCheckURL(username, site)
				results <- result
			}
		}()
	}

	// Send jobs
	go func() {
		for _, site := range sites.WhatsmynameSites {
			jobs <- site
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
