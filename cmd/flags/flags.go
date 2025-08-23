package flags

import (
	"flag"
	"fmt"

	"github.com/ppmpreetham/vesper/pkg/types"
	"github.com/ppmpreetham/vesper/tools"
)

// Parse parses command-line flags and returns the configuration
func Parse() (types.Config, string, bool) {
	// Default configuration
	config := types.Config{
		Database:   "",
		Timeout:    7,
		NumWorkers: 250,
		BufferSize: 1000,
	}

	// both short and long flags
	helpFlag := flag.Bool("help", false, "Show help message")
	flag.BoolVar(helpFlag, "h", false, "Show help message")

	versionFlag := flag.Bool("version", false, "Show version information")
	flag.BoolVar(versionFlag, "v", false, "Show version information")

	databaseFlag := flag.String("database", "", "Enumerate on a specific database (default: all)")
	flag.StringVar(databaseFlag, "d", "", "Enumerate on a specific database (default: all)")

	timeoutFlag := flag.Int("timeout", 7, "HTTP request timeout in seconds (default: 7)")
	flag.IntVar(timeoutFlag, "t", 7, "HTTP request timeout in seconds (default: 7)")

	// Define usage
	flag.Usage = func() {
		fmt.Println("Usage: vesper <username> [options]")
		fmt.Println("Options:")
		fmt.Println("  -h, --help\t\tShow this help message")
		fmt.Println("  -v, --version\t\tShow version information")
		fmt.Println("  -d, --database\tEnumerate using a specific database (default: whatsmyname)")
		fmt.Println("  -t, --timeout\t\tHTTP request timeout in seconds (default: 7) (high for better results, more time)")
		fmt.Println("\nList of databases:\n\t- whatsmyname (default)\n\t- sherlock\n\t- maigret\n\t- all")
	}

	// Parse flags
	flag.Parse()

	// Check for help and version flags
	if *helpFlag {
		flag.Usage()
		return config, "", true
	}

	if *versionFlag {
		fmt.Println("Vesper version 1.0.0")
		return config, "", true
	}

	// USERNAME arg
	username := flag.Arg(0)
	if username == "" {
		tools.Red("Error: Username is required\n")
		flag.Usage()
		return config, "", true
	}

	// Validate timeout
	if *timeoutFlag < 1 || *timeoutFlag > 60 {
		fmt.Println("Error: Timeout must be between 1 and 60 seconds")
		return config, "", true
	}

	// Update config
	config.Database = *databaseFlag
	config.Timeout = *timeoutFlag

	return config, username, false
}
