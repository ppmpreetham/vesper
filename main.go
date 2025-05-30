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

	// both short and long flags
	helpFlag := flag.Bool("help", false, "Show help message")
	flag.BoolVar(helpFlag, "h", false, "Show help message")

	versionFlag := flag.Bool("version", false, "Show version information")
	flag.BoolVar(versionFlag, "v", false, "Show version information")

	databaseFlag := flag.String("database", "", "Enumerate on a specific database (default: all)")
	flag.StringVar(databaseFlag, "d", "", "Enumerate on a specific database (default: all)")

	// command-line flags
	flag.Parse()
	flag.Usage = func() {
		fmt.Println("Usage: vesper <username> [options]")
		fmt.Println("Options:")
		fmt.Println("  -h, --help\t\tShow this help message")
		fmt.Println("  -v, --version\t\tShow version information")
		fmt.Println("  -d, --database\tEnumerate on a specific database (default: all)\nList of databases:\n\t- sherlock\n\t- maigret\n\t- whatsmyname")
	}

	// flags
	if *helpFlag {
		flag.Usage()
		return
	}

	if *versionFlag {
		fmt.Println("Vesper version 1.0.0")
		return
	}

	// USERNAME arg
	username := flag.Arg(0)
	if username == "" {
		fmt.Println("Error: Username is required")
		flag.Usage()
		return
	}

	fmt.Println("Starting enumeration for username:", username)

	startTime := time.Now()

	var wg sync.WaitGroup
	buffersize := 1000

	// Check which database to use
	switch *databaseFlag {
	case "sherlock":
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
				for site := range jobs {
					siteName := <-siteNames
					result := tools.SherlockCheckURL(username, site, siteName)
					results <- result
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
		fmt.Println("Using WhatsMyName database for enumeration...")

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
		fmt.Println("Using Maigret database for enumeration...")

		jobs := make(chan sites.MaigretSiteData, buffersize)
		results := make(chan tools.ReturnData, buffersize)
		siteNames := make(chan string, buffersize)

		// Start worker pool for Maigret
		numWorkers := 250
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for site := range jobs {
					siteName := <-siteNames
					result := tools.MaigretCheckURL(username, site, siteName)
					results <- result
				}
			}()
		}

		// Send Maigret jobs using the static map
		go func() {
			for siteName, site := range sites.MaigretSites {
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

	case "": // Default case - run all databases
		fmt.Println("Using all databases for enumeration...")

		totalFoundCount := 0

		// Run Sherlock database first
		fmt.Println("\n=== Starting Sherlock database enumeration ===")

		jobs := make(chan sites.SherlockSiteData, buffersize)
		results := make(chan tools.ReturnData, buffersize)
		siteNames := make(chan string, buffersize)

		// Start worker pool for Sherlock
		numWorkers := 250
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for site := range jobs {
					siteName := <-siteNames
					result := tools.SherlockCheckURL(username, site, siteName)
					results <- result
				}
			}()
		}

		// Send Sherlock jobs
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

		// Collect and print Sherlock results
		sherlockCount := 0
		for result := range results {
			if result.Status == "FOUND" {
				sherlockCount++
				fmt.Println("Found:", result.Name, "at", result.URL)
			}
		}
		fmt.Printf("Sherlock database completed - Found %d matches\n", sherlockCount)
		totalFoundCount += sherlockCount

		// Reset WaitGroup for next database
		wg = sync.WaitGroup{}

		// Run WhatsMyName database second
		fmt.Println("\n=== Starting WhatsMyName database enumeration ===")

		wmnJobs := make(chan sites.WhatsmynameSiteData, buffersize)
		wmnResults := make(chan tools.ReturnData, buffersize)

		// Start worker pool for WhatsMyName
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for site := range wmnJobs {
					result := tools.WhatsMyNameCheckURL(username, site)
					wmnResults <- result
				}
			}()
		}

		// Send WhatsMyName jobs
		go func() {
			for _, site := range sites.WhatsmynameSites {
				wmnJobs <- site
			}
			close(wmnJobs)
		}()

		// Wait and close results
		go func() {
			wg.Wait()
			close(wmnResults)
		}()

		// Collect and print WhatsMyName results
		wmnCount := 0
		for result := range wmnResults {
			if result.Status == "FOUND" {
				wmnCount++
				fmt.Println("Found:", result.Name, "at", result.URL)
			}
		}
		fmt.Printf("WhatsMyName database completed - Found %d matches\n", wmnCount)
		totalFoundCount += wmnCount

		// Reset WaitGroup for next database
		wg = sync.WaitGroup{}

		// Run Maigret database third
		fmt.Println("\n=== Starting Maigret database enumeration ===")

		maigretJobs := make(chan sites.MaigretSiteData, buffersize)
		maigretResults := make(chan tools.ReturnData, buffersize)
		maigretSiteNames := make(chan string, buffersize)

		// Start worker pool for Maigret
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for site := range maigretJobs {
					siteName := <-maigretSiteNames
					result := tools.MaigretCheckURL(username, site, siteName)
					maigretResults <- result
				}
			}()
		}

		// Send Maigret jobs
		go func() {
			for siteName, site := range sites.MaigretSites {
				maigretJobs <- site
				maigretSiteNames <- siteName
			}
			close(maigretJobs)
			close(maigretSiteNames)
		}()

		// Wait and close results
		go func() {
			wg.Wait()
			close(maigretResults)
		}()

		// Collect and print Maigret results
		maigretCount := 0
		for result := range maigretResults {
			if result.Status == "FOUND" {
				maigretCount++
				fmt.Println("Found:", result.Name, "at", result.URL)
			}
		}
		fmt.Printf("Maigret database completed - Found %d matches\n", maigretCount)
		totalFoundCount += maigretCount

		elapsedTime := time.Since(startTime)
		fmt.Printf("\n=== Summary ===\n")
		fmt.Printf("Execution completed in %s\n", elapsedTime)
		fmt.Printf("Total found username on %d sites across all databases\n", totalFoundCount)
		fmt.Printf("  - Sherlock: %d matches\n", sherlockCount)
		fmt.Printf("  - WhatsMyName: %d matches\n", wmnCount)
		fmt.Printf("  - Maigret: %d matches\n", maigretCount)

	default:
		fmt.Printf("Unknown database: %s\n", *databaseFlag)
		fmt.Println("Available databases: sherlock, whatsmyname, maigret")
		return
	}
}
