package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/ppmpreetham/vesper/sites"
	"github.com/ppmpreetham/vesper/tools"
)

func main() {
	// Parse command-line flags
	flag.Parse()
	flag.Usage = func() {
		fmt.Println("Usage: vesper <username> [options]")
		fmt.Println("Options:")
		fmt.Println("  -h, --help\t\tShow this help message")
		fmt.Println("  -v, --version\t\tShow version information")
		fmt.Println("  -d, --database\tEnumerate on a specific database (default: all)\nList of databases:\n\t- sherlock\n\t- maigret\n\t- whatsmyname")
	}

	flag.Bool("help", false, "Show help message")
	flag.Bool("version", false, "Show version information")
	flag.String("database", "", "Enumerate on a specific database (default: all)")

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

	// DATABASE arg
	database := flag.Lookup("database").Value.String()
	if database != "" {
		fmt.Println("Enumerating on database:", database)
		if database != "whatsmyname" {
			fmt.Println("Error: Currently only 'whatsmyname' database is supported")
			return
		}
	} else {
		fmt.Println("Enumerating on all databases")
		database = "whatsmyname" // Default to whatsmyname if no database specified
	}

	fmt.Println("Starting enumeration for username:", username)

	startTime := time.Now()

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
