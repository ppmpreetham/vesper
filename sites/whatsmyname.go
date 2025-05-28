package sites

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// WhatsmynameSiteData represents a single site definition from whatsmyname_fix.json
type WhatsmynameSiteData struct {
	Name         string            `json:"name"`
	URICheck     string            `json:"uri_check"`
	URIPretty    string            `json:"uri_pretty,omitempty"`
	PostBody     string            `json:"post_body,omitempty"`
	ECode        int               `json:"e_code"`
	EString      string            `json:"e_string"`
	MString      string            `json:"m_string"`
	MCode        int               `json:"m_code"`
	Known        []string          `json:"known"`
	Cat          string            `json:"cat"`
	StripBadChar string            `json:"strip_bad_char,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	Protection   []string          `json:"protection,omitempty"`
}

// WhatsmynameData represents the root structure of whatsmyname_fix.json
type WhatsmynameData struct {
	License    []string              `json:"license"`
	Authors    []string              `json:"authors"`
	Categories []string              `json:"categories"`
	Sites      []WhatsmynameSiteData `json:"sites"`
}

// LoadWhatsmynameSites reads the whatsmyname JSON file and returns the site data
func LoadWhatsmynameSites(jsonFilePath string) ([]WhatsmynameSiteData, error) {
	// Read the JSON file
	jsonFile, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	// Create a struct to hold the unmarshaled data
	var whatsmynameData WhatsmynameData

	// Unmarshal the JSON data
	err = json.Unmarshal(jsonFile, &whatsmynameData)
	if err != nil {
		return nil, err
	}

	return whatsmynameData.Sites, nil
}

// InitWhatsmynameSites initializes the WhatsmynameSites slice from the default JSON file location
func InitWhatsmynameSites() error {
	// Define default path - adjust as needed
	defaultPath := filepath.Join(".", "whatsmyname_fix.json")

	sites, err := LoadWhatsmynameSites(defaultPath)
	if err != nil {
		return err
	}

	// Assign to the package-level variable
	WhatsmynameSites = sites
	return nil
}
