package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/ppmpreetham/vesper/sites"
	"github.com/ppmpreetham/vesper/tools"
)

// avoid synchronization issues
type SherlockJob struct {
	Name string
	Data sites.SherlockSiteData
}

type MaigretJob struct {
	Name string
	Data sites.MaigretSiteData
}

func main() {
	tools.PrintLogo() // logo goes here

	// both short and long flags
	helpFlag := flag.Bool("help", false, "Show help message")
	flag.BoolVar(helpFlag, "h", false, "Show help message")

	versionFlag := flag.Bool("version", false, "Show version information")
	flag.BoolVar(versionFlag, "v", false, "Show version information")

	databaseFlag := flag.String("database", "", "Enumerate on a specific database (default: all)")
	flag.StringVar(databaseFlag, "d", "", "Enumerate on a specific database (default: all)")

	timeoutFlag := flag.Int("timeout", 7, "HTTP request timeout in seconds (default: 7)")
	flag.IntVar(timeoutFlag, "t", 7, "HTTP request timeout in seconds (default: 7)")

	// command-line flags
	flag.Parse()
	flag.Usage = func() {
		fmt.Println("Usage: vesper <username> [options]")
		fmt.Println("Options:")
		fmt.Println("  -h, --help\t\tShow this help message")
		fmt.Println("  -v, --version\t\tShow version information")
		fmt.Println("  -d, --database\tEnumerate on a specific database (default: all)")
		fmt.Println("  -t, --timeout\t\tHTTP request timeout in seconds (default: 7) (high for better results, more time)")
		fmt.Println("\nList of databases:\n\t- whatsmyname\n\t- maigret\n\t- sherlock")
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
		tools.Red("Error: Username is required\n")
		flag.Usage()
		return
	}

	// Validate timeout
	if *timeoutFlag < 1 || *timeoutFlag > 60 {
		fmt.Println("Error: Timeout must be between 1 and 60 seconds")
		return
	}

	// Set the timeout in tools package
	tools.SetHTTPTimeout(time.Duration(*timeoutFlag) * time.Second)

	fmt.Printf("Starting enumeration for username: %s (timeout: %ds)\n", username, *timeoutFlag)

	startTime := time.Now()

	var wg sync.WaitGroup
	buffersize := 1000

	// Check which database to use
	switch *databaseFlag {
	case "sherlock":
		fmt.Println("Using Sherlock database for enumeration...")

		jobs := make(chan SherlockJob, buffersize)
		results := make(chan tools.ReturnData, buffersize)

		// Start worker pool for Sherlock
		numWorkers := 250
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for job := range jobs {
					result := tools.SherlockCheckURL(username, job.Data, job.Name)
					results <- result
				}
			}()
		}

		// Send Sherlock jobs using the static map
		go func() {
			for siteName, site := range sites.SherlockSites {
				jobs <- SherlockJob{
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
				tools.Green("Found: ")
				fmt.Print(result.Name, " at ")
				tools.BoldGreen(result.URL)
				fmt.Print("\n")
			}
		}

		elapsedTime := time.Since(startTime)
		fmt.Printf("\nExecution completed in %s\n", elapsedTime)
		fmt.Printf("Found username on %d sites\n", foundCount)

	case "maigret":
		fmt.Println("Using Maigret database for enumeration...")

		jobs := make(chan MaigretJob, buffersize)
		results := make(chan tools.ReturnData, buffersize)

		// Start worker pool for Maigret
		numWorkers := 250
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for job := range jobs {
					result := tools.MaigretCheckURL(username, job.Data, job.Name)
					results <- result
				}
			}()
		}

		// Send Maigret jobs using the static map
		go func() {
			for siteName, site := range sites.MaigretSites {
				jobs <- MaigretJob{
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

		elapsedTime := time.Since(startTime)
		fmt.Printf("\nExecution completed in %s\n", elapsedTime)
		fmt.Printf("Found username on %d sites\n", foundCount)

	case "": // Default case - run all databases
		fmt.Println("Using all databases for enumeration...")

		totalFoundCount := 0

		// Run Sherlock database first
		fmt.Print("\n=== Starting ")
		tools.BoldOrange("Sherlock")
		fmt.Println(" database enumeration ===")

		tools.ResetHTTPClient() // reset before each database

		sherlockJobs := make(chan SherlockJob, buffersize)
		sherlockResults := make(chan tools.ReturnData, buffersize)

		// Start worker pool for Sherlock
		numWorkers := 250
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for job := range sherlockJobs {
					result := tools.SherlockCheckURL(username, job.Data, job.Name)
					sherlockResults <- result
				}
			}()
		}

		// Send Sherlock jobs
		go func() {
			for siteName, site := range sites.SherlockSites {
				sherlockJobs <- SherlockJob{
					Name: siteName,
					Data: site,
				}
			}
			close(sherlockJobs)
		}()

		// Wait and close results
		go func() {
			wg.Wait()
			close(sherlockResults)
		}()

		// Collect and print Sherlock results
		sherlockCount := 0
		for result := range sherlockResults {
			if result.Status == "FOUND" {
				sherlockCount++
				tools.Green("Found: ")
				fmt.Print(result.Name, " at ")
				tools.BoldGreen(result.URL)
				fmt.Print("\n")
			}
		}
		fmt.Printf("Sherlock database completed - Found %d matches\n", sherlockCount)
		totalFoundCount += sherlockCount

		// reset WaitGroup and HTTP client for next database
		wg = sync.WaitGroup{}
		tools.ResetHTTPClient() // reset between databases

		// Run WhatsMyName database second
		fmt.Print("\n=== Starting ")
		tools.BoldOrange("WhatsMyName")
		fmt.Println(" database enumeration ===")
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
				tools.Green("Found: ")
				fmt.Print(result.Name, " at ")
				tools.BoldGreen(result.URL)
				fmt.Print("\n")
			}
		}
		fmt.Printf("WhatsMyName database completed - Found %d matches\n", wmnCount)
		totalFoundCount += wmnCount

		// reset WaitGroup and HTTP client for next database
		wg = sync.WaitGroup{}
		tools.ResetHTTPClient() // reset between databases

		// Run Maigret database third
		fmt.Print("\n=== Starting ")
		tools.BoldOrange("Maigret")
		fmt.Println(" database enumeration ===")

		maigretJobs := make(chan MaigretJob, buffersize)
		maigretResults := make(chan tools.ReturnData, buffersize)

		// Start worker pool for Maigret
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for job := range maigretJobs {
					result := tools.MaigretCheckURL(username, job.Data, job.Name)
					maigretResults <- result
				}
			}()
		}

		// Send Maigret jobs
		go func() {
			for siteName, site := range sites.MaigretSites {
				maigretJobs <- MaigretJob{
					Name: siteName,
					Data: site,
				}
			}
			close(maigretJobs)
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
				tools.Green("Found: ")
				fmt.Print(result.Name, " at ")
				tools.BoldGreen(result.URL)
				fmt.Print("\n")
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
