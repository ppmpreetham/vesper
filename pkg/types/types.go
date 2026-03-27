package types

import "github.com/ppmpreetham/vesper/sites"

type Config struct {
	Database   string
	Timeout    int
	NumWorkers int
	BufferSize int
}

type JobResult struct {
	DatabaseName string
	FoundCount   int
}

type SherlockJob struct {
	Name string
	Data sites.SherlockSiteData
}

type MaigretJob struct {
	Name string
	Data sites.MaigretSiteData
}
