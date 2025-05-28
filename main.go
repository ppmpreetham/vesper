package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/ppmpreetham/vesper/sites"
	"github.com/ppmpreetham/vesper/tools"
)

func main() {
	startTime := time.Now()

	username := "ppmpreetham"
	var wg sync.WaitGroup
	buffersize := 1000

	jobs := make(chan sites.WhatsmynameSiteData, buffersize)
	results := make(chan tools.ReturnData, buffersize)

	// Start worker pool
	numWorkers := 200
	for i := 0; i < numWorkers; i++ {
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
			fmt.Println("Found:", result.Name, "at", result.URL)
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("\nExecution completed in %s\n", elapsedTime)
	fmt.Printf("Found username on %d sites\n", foundCount)
}
