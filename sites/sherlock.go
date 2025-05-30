package sites

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ErrorMessage []string

// Custom unmarshaling to handle both string and []string for errorMsg
func (em *ErrorMessage) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as array first
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*em = ErrorMessage(arr)
		return nil
	}

	// If that fails, try as single string
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*em = ErrorMessage([]string{str})
		return nil
	}

	return nil // Return empty slice if both fail
}

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
