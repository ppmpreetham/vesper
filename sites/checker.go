package sites

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

func CheckUsername(site Site, username string) (found bool, err error) {
	client := &http.Client{
		Timeout: time.Duration(site.TimeoutSeconds) * time.Second,
	}
	url := fmt.Sprintf(site.UrlTemplate, username)
	req, err := http.NewRequest(site.Method, url, nil)
	if err != nil {
		return false, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Check status code if set
	if site.ExistStatusCode != 0 && resp.StatusCode == site.ExistStatusCode {
		return true, nil
	}
	if site.NotFoundStatusCode != 0 && resp.StatusCode == site.NotFoundStatusCode {
		return false, nil
	}

	// Optionally check for error message in body
	if site.ErrorRegex != "" {
		body, _ := io.ReadAll(resp.Body)
		matched, _ := regexp.Match(site.ErrorRegex, body)
		return !matched, nil
	}

	return false, nil
}
