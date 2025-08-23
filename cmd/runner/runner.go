package runner

import (
	"fmt"

	"github.com/ppmpreetham/vesper/pkg/types"
	"github.com/ppmpreetham/vesper/tools"
)

// RunEnumeration runs the enumeration against the specified databases
func RunEnumeration(username string, config types.Config) []types.JobResult {
	fmt.Printf("Starting enumeration for username: %s (timeout: %ds)\n", tools.OrangeString(username), config.Timeout)

	var results []types.JobResult

	switch config.Database {
	case "sherlock":
		fmt.Println("Using Sherlock database for enumeration...")
		count := RunSherlockDatabase(username, config)
		results = append(results, types.JobResult{DatabaseName: "Sherlock", FoundCount: count})

	case "maigret":
		fmt.Println("Using Maigret database for enumeration...")
		count := RunMaigretDatabase(username, config)
		results = append(results, types.JobResult{DatabaseName: "Maigret", FoundCount: count})

	case "all": // New option to run all databases
		fmt.Println("Using all databases for enumeration...")
		results = RunAllDatabases(username, config)

	case "whatsmyname", "": // Default case - run only WhatsMyName
		fmt.Println("Using WhatsMyName database for enumeration...")
		count := RunWhatsMyNameDatabase(username, config)
		results = append(results, types.JobResult{DatabaseName: "WhatsMyName", FoundCount: count})

	default:
		fmt.Printf("Unknown database: %s\n", config.Database)
		fmt.Println("Available databases: sherlock, whatsmyname, maigret, all")
	}

	return results
}

// RunAllDatabases runs enumeration against all available databases
func RunAllDatabases(username string, config types.Config) []types.JobResult {
	var results []types.JobResult

	// Run Sherlock database first
	fmt.Print("\n=== Starting ")
	tools.BoldOrange("Sherlock")
	fmt.Println(" database enumeration ===")
	tools.ResetHTTPClient()
	sherlockCount := RunSherlockDatabase(username, config)
	fmt.Printf("Sherlock database completed - Found %d matches\n", sherlockCount)
	results = append(results, types.JobResult{DatabaseName: "Sherlock", FoundCount: sherlockCount})

	// Run WhatsMyName database second
	fmt.Print("\n=== Starting ")
	tools.BoldOrange("WhatsMyName")
	fmt.Println(" database enumeration ===")
	tools.ResetHTTPClient()
	wmnCount := RunWhatsMyNameDatabase(username, config)
	fmt.Printf("WhatsMyName database completed - Found %d matches\n", wmnCount)
	results = append(results, types.JobResult{DatabaseName: "WhatsMyName", FoundCount: wmnCount})

	// Run Maigret database third
	fmt.Print("\n=== Starting ")
	tools.BoldOrange("Maigret")
	fmt.Println(" database enumeration ===")
	tools.ResetHTTPClient()
	maigretCount := RunMaigretDatabase(username, config)
	fmt.Printf("Maigret database completed - Found %d matches\n", maigretCount)
	results = append(results, types.JobResult{DatabaseName: "Maigret", FoundCount: maigretCount})

	return results
}
