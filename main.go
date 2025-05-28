package main

import (
	"sync"

	"github.com/ppmpreetham/vesper/sites"
	"github.com/ppmpreetham/vesper/tools"
)

func main() {
	username := "ppmpreetham"
	var wg sync.WaitGroup

	for _, site := range sites.WhatsmynameSites {
		wg.Add(1)
		go tools.WhatsMyNameCheckURL(username, site, &wg)
	}

	wg.Wait()
}
