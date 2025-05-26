package sites

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

func CheckUsername(site Site, username string) (found bool, finalURL string, err error) {

	// Validate username against regex
	re, err := regexp.Compile(site.UsernameRegex)
	if err != nil {
		return false, "", fmt.Errorf("invalid username regex: %w", err)
	}
	// Check if the username matches the regex
	if !re.MatchString(username) {
		return false, "", fmt.Errorf("username '%s' does not match regex '%s'", username, site.UsernameRegex)
	}

	// Check if the site is disabled
	if site.Disabled {
		return false, "", fmt.Errorf("site '%s' is disabled", site.Name)
	}

	// Check if timeout is specified
	var client *http.Client

	if site.TimeoutSeconds != 0 {
		client = &http.Client{
			Timeout: time.Duration(site.TimeoutSeconds) * time.Second,
		}
	} else {
		client = &http.Client{}
	}

	url := fmt.Sprintf(site.UrlTemplate, username)
	req, err := http.NewRequest(site.Method, url, nil)
	if err != nil {
		return false, url, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, url, err
	}
	defer resp.Body.Close()

	finalURL = resp.Request.URL.String()

	// Check status code if set
	if site.ExistStatusCode != 0 && resp.StatusCode == site.ExistStatusCode {
		return true, finalURL, nil
	}
	if site.NotFoundStatusCode != 0 && resp.StatusCode == site.NotFoundStatusCode {
		return false, finalURL, nil
	}

	// Optionally check for error message in body
	if site.ErrorRegex != "" {
		body, _ := io.ReadAll(resp.Body)
		matched, _ := regexp.Match(site.ErrorRegex, body)
		return !matched, finalURL, nil
	}

	return false, finalURL, nil
}
