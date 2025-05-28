package sites

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ErrorMessage []string

type SherlockSiteData struct {
	ErrorMsg        ErrorMessage      `json:"errorMsg,omitempty"`
	ErrorType       string            `json:"errorType"`
	RegexCheck      string            `json:"regexCheck,omitempty"`
	URL             string            `json:"url"`
	URLMain         string            `json:"urlMain"`
	URLProbe        string            `json:"urlProbe,omitempty"`
	UsernameClaimed string            `json:"username_claimed"`
	IsNSFW          bool              `json:"isNSFW,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	RequestMethod   string            `json:"request_method,omitempty"`
	RequestPayload  interface{}       `json:"request_payload,omitempty"`
	ErrorCode       int               `json:"errorCode,omitempty"`
	ErrorURL        string            `json:"errorUrl,omitempty"`
}

// LoadSherlockSites reads the sherlock JSON file and returns a map of site data
func LoadSherlockSites(jsonFilePath string) (map[string]SherlockSiteData, error) {
	// Read the JSON file
	jsonFile, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	// Create a map to hold the unmarshaled data
	var sitesMap map[string]SherlockSiteData

	// Unmarshal the JSON data into the map
	err = json.Unmarshal(jsonFile, &sitesMap)
	if err != nil {
		return nil, err
	}

	return sitesMap, nil
}

// InitSherlockSites initializes the SherlockSites map from the default JSON file location
func InitSherlockSites() error {
	// Define default path - adjust as needed
	defaultPath := filepath.Join(".", "sherlock_fix.json")

	sites, err := LoadSherlockSites(defaultPath)
	if err != nil {
		return err
	}

	// Assign to the package-level variable
	SherlockSites = sites
	return nil
}
