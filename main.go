package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/ppmpreetham/vesper/sites"
	"github.com/ppmpreetham/vesper/tools"
	"github.com/ppmpreetham/vesper/utils"
)

func main() {
	utils.PrintLogo() // logo goes here

	// Define flags before parsing
	helpFlag := flag.Bool("help", false, "Show help message")
	versionFlag := flag.Bool("version", false, "Show version information")
	databaseFlag := flag.String("database", "", "Enumerate on a specific database (default: all)")

	// Parse command-line flags
	flag.Parse()
	flag.Usage = func() {
		fmt.Println("Usage: vesper <username> [options]")
		fmt.Println("Options:")
		fmt.Println("  -h, --help\t\tShow this help message")
		fmt.Println("  -v, --version\t\tShow version information")
		fmt.Println("  -d, --database\tEnumerate on a specific database (default: all)\nList of databases:\n\t- sherlock\n\t- maigret\n\t- whatsmyname")
	}

	// Handle flags
	if *helpFlag {
		flag.Usage()
		return
	}

	if *versionFlag {
		fmt.Println("Vesper version 1.0.0") // Update version as needed
		return
	}

	// USERNAME arg
	username := flag.Arg(0)
	if username == "" {
		fmt.Println("Error: Username is required")
		flag.Usage()
		return
	}

	if username[0] == '@' {
		username = username[1:] // Remove '@' if present
	}

	fmt.Println("Starting enumeration for username:", username)

	startTime := time.Now()

	var wg sync.WaitGroup
	buffersize := 1000

	// Check which database to use
	switch *databaseFlag {
	case "sherlock", "":
		fmt.Println("Using Sherlock database for enumeration...")

		jobs := make(chan sites.SherlockSiteData, buffersize)
		results := make(chan tools.ReturnData, buffersize)
		siteNames := make(chan string, buffersize)

		// Start worker pool for Sherlock
		numWorkers := 250
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case site, ok := <-jobs:
						if !ok {
							return
						}
						siteName := <-siteNames
						result := tools.SherlockCheckURL(username, site, siteName)
						results <- result
					}
				}
			}()
		}

		// Send Sherlock jobs using the static map
		go func() {
			for siteName, site := range sites.SherlockSites {
				jobs <- site
				siteNames <- siteName
			}
			close(jobs)
			close(siteNames)
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

	case "whatsmyname":
		// Default behavior - use WhatsMyName
		jobs := make(chan sites.WhatsmynameSiteData, buffersize)
		results := make(chan tools.ReturnData, buffersize)

		// Start worker pool
		numWorkers := 250
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

	case "maigret":
		fmt.Println("Maigret database not implemented yet")
		return

	default:
		fmt.Printf("Unknown database: %s\n", *databaseFlag)
		fmt.Println("Available databases: sherlock, whatsmyname, maigret")
		return
	}
}
