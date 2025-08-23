package types

import "github.com/ppmpreetham/vesper/sites"

// Config holds application configuration
type Config struct {
	Database   string
	Timeout    int
	NumWorkers int
	BufferSize int
}

// JobResult holds the results from a database run
type JobResult struct {
	DatabaseName string
	FoundCount   int
}

// SherlockJob represents a job for the Sherlock database
type SherlockJob struct {
	Name string
	Data sites.SherlockSiteData
}

// MaigretJob represents a job for the Maigret database
type MaigretJob struct {
	Name string
	Data sites.MaigretSiteData
}
