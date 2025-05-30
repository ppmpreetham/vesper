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

const USERAGENT = "Mozilla/5.0 (X11; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0"

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
	req.Header.Set("User-Agent", USERAGENT)

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

	// Use URLProbe if available, otherwise use URL
	checkURL := result.URL
	if site.URLProbe != "" {
		checkURL = fmt.Sprintf(site.URLProbe, username)
	}

	// Determine HTTP method
	method := "GET"
	if site.RequestMethod != "" {
		method = site.RequestMethod
	}

	// Prepare request body for POST requests
	var requestBody io.Reader
	if method == "POST" && site.RequestPayload != nil {
		if payloadStr, ok := site.RequestPayload.(string); ok {
			requestBody = strings.NewReader(fmt.Sprintf(payloadStr, username))
		}
	}

	// Prepare custom HTTP request
	req, err := http.NewRequest(method, checkURL, requestBody)
	if err != nil {
		result.Status = "ERROR"
		return result
	}

	// Set default User-Agent
	req.Header.Set("User-Agent", USERAGENT)

	// Set custom headers if provided
	if site.Headers != nil {
		for key, value := range site.Headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		result.Status = "ERROR"
		return result
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode >= 500 {
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
		// Default error code is 404 if not specified
		errorCode := 404
		if site.ErrorCode != 0 {
			errorCode = site.ErrorCode
		}

		if resp.StatusCode == errorCode {
			result.Status = "NOT FOUND"
		} else {
			result.Status = "FOUND"
		}
	case "message":
		foundError := false
		// Convert ErrorMessage to []string for checking
		errorMessages := []string(site.ErrorMsg)
		for _, errMsg := range errorMessages {
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
