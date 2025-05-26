package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	// "github.com/ppmpreetham/vesper/sites"
)

func main() {

	// colors
	info := color.New(color.FgCyan).SprintFunc()
	warning := color.New(color.FgYellow).SprintFunc()
	errorMsg := color.New(color.FgRed).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	// flags
	timeout := flag.Int("timeout", 10, "Set custom timeout for requests (in seconds)")
	threads := flag.Int("threads", 10, "Set number of concurrent threads")
	rateLimit := flag.Int("rate-limit", 1, "Set rate limit between requests (in seconds)")
	verbose := flag.Bool("verbose", false, "Enable verbose output")

	// vesper usage
	flag.Usage = func() {
		fmt.Println(success("Vesper: A tool to check username availability across multiple sites."))
		fmt.Println(info("USAGE:"), "vesper", bold("<username>"), "[flags]")
		fmt.Println(info("Flags:"))
		flag.PrintDefaults()
	}

	flag.Parse()

	// Ensure username is provided
	if flag.NArg() < 1 {
		fmt.Println(errorMsg("ERROR:"), "Username is required.")
		flag.Usage()
		os.Exit(1)
	}

	username := flag.Arg(0)
	fmt.Println(success("Checking availability for username:"), bold(username))

	// Optional: simulate using the flags
	fmt.Println(info("Vesper Server is starting..."))
	fmt.Printf("%s Timeout: %d seconds\n", info("-"), *timeout)
	fmt.Printf("%s Threads: %d\n", info("-"), *threads)
	fmt.Printf("%s Rate Limit: %d second(s)\n", info("-"), *rateLimit)

	if *verbose {
		fmt.Println(warning("Verbose mode is enabled."))
	}

	fmt.Println("Work done.")
}
