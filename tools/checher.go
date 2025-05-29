package tools

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ppmpreetham/vesper/sites"
)

type ReturnData struct {
	Name     string
	URL      string
	Status   string
	Metadata map[string]string
}

var httpClient = &http.Client{
	Timeout: 7 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 200,
		IdleConnTimeout:     90 * time.Second,
	},
}

func WhatsMyNameCheckURL(username string, site sites.WhatsmynameSiteData) ReturnData {
	result := ReturnData{
		Name:     site.Name,
		URL:      fmt.Sprintf(site.URICheck, username),
		Status:   "NOT FOUND",
		Metadata: make(map[string]string),
	}

	// Prepare custom HTTP request with User-Agent
	req, err := http.NewRequest("GET", result.URL, nil)
	if err != nil {
		result.Status = "ERROR"
		return result
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		result.Status = "ERROR"
		return result
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode >= 500 {
		resp.Body.Close()
		result.Status = "ERROR"
		return result
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Status = "ERROR"
		return result
	}
	bodyStr := string(bodyBytes)

	// Match logic
	if strings.Contains(bodyStr, site.EString) && site.ECode == resp.StatusCode {
		if !strings.Contains(bodyStr, site.MString) {
			mCodeCondition := site.MCode == site.ECode || site.MCode != resp.StatusCode
			if mCodeCondition {
				result.Status = "FOUND"
			}
		}
	}

	return result
}

func SherlockCheckURL(username string, site sites.SherlockSiteData, Sitename string) ReturnData {
	result := ReturnData{
		Name:     Sitename,
		URL:      fmt.Sprintf(site.URL, username),
		Status:   "NOT FOUND",
		Metadata: make(map[string]string),
	}

	// Regex Check username before passing to URL
	if site.RegexCheck != "" {
		matched, err := regexp.MatchString(site.RegexCheck, username)
		if err != nil || !matched {
			result.Status = "USERNAME CAN'T BE MADE"
			return result
		}
	}

	// Prepare custom HTTP request with User-Agent
	req, err := http.NewRequest("GET", result.URL, nil)
	if err != nil {
		result.Status = "ERROR"
		return result
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		result.Status = "ERROR"
		return result
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode >= 500 {
		resp.Body.Close()
		result.Status = "ERROR"
		return result
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Status = "ERROR"
		return result
	}
	bodyStr := string(bodyBytes)

	// Match logic
	switch site.ErrorType {
	case "status_code":
		if resp.StatusCode == site.ErrorCode {
			result.Status = "NOT FOUND"
		} else {
			result.Status = "FOUND"
		}
	case "message":
		foundError := false
		for _, errMsg := range site.ErrorMsg {
			if strings.Contains(bodyStr, errMsg) {
				foundError = true
				break
			}
		}
		if foundError {
			result.Status = "NOT FOUND"
		} else {
			result.Status = "FOUND"
		}
	case "response_url":
		if resp.Request.URL.String() == site.ErrorURL {
			result.Status = "NOT FOUND"
		} else {
			result.Status = "FOUND"
		}
	default:
		result.Status = "ERROR"
	}
	return result
}
