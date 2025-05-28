package tools

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/ppmpreetham/vesper/sites"
)

type ReturnData struct {
	name     string
	url      string
	status   string
	metadata map[string]string
}

func WhatsMyNameCheckURL(username string, site sites.WhatsmynameSiteData, wg *sync.WaitGroup) ReturnData {
	returnData := ReturnData{
		name:     site.Name,
		url:      "",
		status:   "OK",
		metadata: make(map[string]string),
	}

	formattedURL := fmt.Sprintf(site.URICheck, username)
	// fmt.Println(formattedURL)
	returnData.url = formattedURL

	// Check if the site is reachable
	resp, err := http.Get(formattedURL)
	if err != nil {
		returnData.status = "ERROR"
		return returnData
	}

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		returnData.status = "ERROR"
		return returnData
	}
	bodyStr := string(bodyBytes)

	defer resp.Body.Close()

	// Check if site matches expected conditions
	if strings.Contains(bodyStr, site.EString) && site.ECode == resp.StatusCode {
		// Check if m_string is not in the response body
		if !strings.Contains(bodyStr, site.MString) {
			// Check if (m_code != status_code) when m_code != e_code, otherwise return true
			mCodeCondition := true
			if site.MCode != site.ECode {
				mCodeCondition = site.MCode != resp.StatusCode
			}

			// If both conditions are met, set status to "FOUND"
			if mCodeCondition {
				returnData.status = "FOUND"
				fmt.Println("Found:", site.Name, "at", formattedURL)
			}
		}
	} else {
		returnData.status = "NOT FOUND"
	}

	return returnData
}
