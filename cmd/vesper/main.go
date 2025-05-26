package main

import (
	"flag"
	"fmt"
	"os"

	"sync"

	"github.com/fatih/color"
	"github.com/ppmpreetham/vesper/sites"
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

	// Optional: simulate using the flags
	fmt.Println(info("Vesper is starting..."))
	fmt.Printf("%s Timeout: %d seconds\n", info("-"), *timeout)
	fmt.Printf("%s Threads: %d\n", info("-"), *threads)
	fmt.Printf("%s Rate Limit: %d second(s)\n", info("-"), *rateLimit)

	if *verbose {
		fmt.Println(warning("Verbose mode is enabled."))
	}

	username := flag.Arg(0)
	fmt.Println(success("Checking availability for username:"), bold(username))

	// Set global timeout for all sites
	sites.GlobalTimeoutSeconds = *timeout

	var wg sync.WaitGroup

	for _, site := range sites.DefaultSites {
		wg.Add(1)
		go func(site sites.Site) {
			defer wg.Done()
			found, url, err := sites.CheckUsername(site, username)
			if err != nil {
				fmt.Printf("%s %s: error: %v\n", warning("!"), site.Name, err)
				return
			}
			if found {
				fmt.Printf("%s %s: %s -> %s\n", success("+"), site.Name, success("FOUND"), url)
			} else {
				fmt.Printf("%s %s: %s\n", errorMsg("-"), site.Name, warning("NOT FOUND"))
			}
		}(site)
	}
	wg.Wait()

	fmt.Println("Work done.")
}
