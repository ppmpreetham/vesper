package tools

import (
	"fmt"
	"io"
	"net"
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

const USERAGENT = "Mozilla/5.0 (X11; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0"

var httpTimeout = 7 * time.Second


func SetHTTPTimeout(timeout time.Duration) {
	httpTimeout = timeout
}


func NewHTTPClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   httpTimeout,
		KeepAlive: 30 * time.Second,
		Resolver: &net.Resolver{
			PreferGo: true,
		},
	}

	return &http.Client{
		Timeout: httpTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 50,
			IdleConnTimeout:     30 * time.Second,
			DialContext:         dialer.DialContext,
			DisableKeepAlives:   false,
			ForceAttemptHTTP2:   false,
		},
	}
}

var httpClient = NewHTTPClient()

func ResetHTTPClient() {
	if httpClient != nil {
		httpClient.CloseIdleConnections()
	}
	httpClient = NewHTTPClient()
}

func executeRequest(req *http.Request) (*http.Response, string, error) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", USERAGENT)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return nil, "", fmt.Errorf("server error: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return resp, string(bodyBytes), nil
}


func WhatsMyNameCheckURL(username string, site sites.WhatsmynameSiteData) ReturnData {
	result := ReturnData{
		Name:     site.Name,
		URL:      fmt.Sprintf(site.URICheck, username),
		Status:   "NOT FOUND",
		Metadata: make(map[string]string),
	}

	req, err := http.NewRequest("GET", result.URL, nil)
	if err != nil {
		result.Status = "ERROR"
		return result
	}

	resp, bodyStr, err := executeRequest(req)
	if err != nil {
		result.Status = "ERROR"
		return result
	}

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

	
	if site.RegexCheck != "" {
		matched, err := regexp.MatchString(site.RegexCheck, username)
		if err != nil || !matched {
			result.Status = "USERNAME CAN'T BE MADE"
			return result
		}
	}

	checkURL := result.URL
	if site.URLProbe != "" {
		checkURL = fmt.Sprintf(site.URLProbe, username)
	}

	method := "GET"
	if site.RequestMethod != "" {
		method = site.RequestMethod
	}

	var requestBody io.Reader
	if method == "POST" && site.RequestPayload != nil {
		if payloadStr, ok := site.RequestPayload.(string); ok {
			requestBody = strings.NewReader(fmt.Sprintf(payloadStr, username))
		}
	}

	req, err := http.NewRequest(method, checkURL, requestBody)
	if err != nil {
		result.Status = "ERROR"
		return result
	}

	if site.Headers != nil {
		for key, value := range site.Headers {
			req.Header.Set(key, value)
		}
	}

	resp, bodyStr, err := executeRequest(req)
	if err != nil {
		result.Status = "ERROR"
		return result
	}

	switch site.ErrorType {
	case "status_code":
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

func MaigretCheckURL(username string, site sites.MaigretSiteData, siteName string) ReturnData {
	result := ReturnData{
		Name:     siteName,
		Status:   "NOT FOUND",
		Metadata: make(map[string]string),
	}

	if site.Disabled != nil && *site.Disabled {
		result.Status = "DISABLED"
		return result
	}

	if site.URL == nil {
		result.Status = "ERROR"
		return result
	}

	url := strings.ReplaceAll(*site.URL, "{username}", username)
	result.URL = url

	if site.RegexCheck != nil {
		matched, err := regexp.MatchString(*site.RegexCheck, username)
		if err != nil || !matched {
			result.Status = "USERNAME CAN'T BE MADE"
			return result
		}
	}

	checkURL := url
	if site.URLProbe != nil {
		checkURL = strings.ReplaceAll(*site.URLProbe, "{username}", username)
	}

	req, err := http.NewRequest("GET", checkURL, nil)
	if err != nil {
		result.Status = "ERROR"
		return result
	}

	resp, bodyStr, err := executeRequest(req)
	if err != nil {
		result.Status = "ERROR"
		return result
	}

	checkType := "status_code"
	if site.CheckType != nil {
		checkType = *site.CheckType
	}

	switch checkType {
	case "status_code":
		if resp.StatusCode == 200 {
			result.Status = "FOUND"
		} else {
			result.Status = "NOT FOUND"
		}
	case "message":
		foundPresence := false

		for _, presStr := range site.PresenseStrs {
			if strings.Contains(bodyStr, presStr) {
				foundPresence = true
				break
			}
		}

		if !foundPresence {
			for _, presStr := range site.PresenceStrs {
				if strings.Contains(bodyStr, presStr) {
					foundPresence = true
					break
				}
			}
		}

		foundAbsence := false
		for _, absStr := range site.AbsenceStrs {
			if strings.Contains(bodyStr, absStr) {
				foundAbsence = true
				break
			}
		}

		if foundPresence && !foundAbsence {
			result.Status = "FOUND"
		} else if foundAbsence {
			result.Status = "NOT FOUND"
		} else {
			if resp.StatusCode == 200 {
				result.Status = "FOUND"
			} else {
				result.Status = "NOT FOUND"
			}
		}
	default:
		result.Status = "ERROR"
	}

	return result
}