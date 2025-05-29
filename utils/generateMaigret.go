package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// SiteConfig represents the structure of each site in maigret_fix.json
type SiteConfig struct {
	Tags              []string          `json:"tags,omitempty"`
	Engine            *string           `json:"engine,omitempty"`
	AlexaRank         *int              `json:"alexaRank,omitempty"`
	URLMain           *string           `json:"urlMain,omitempty"`
	URL               *string           `json:"url,omitempty"`
	URLSubpath        *string           `json:"urlSubpath,omitempty"`
	URLProbe          *string           `json:"urlProbe,omitempty"`
	UsernameClaimed   *string           `json:"usernameClaimed,omitempty"`
	UsernameUnclaimed *string           `json:"usernameUnclaimed,omitempty"`
	CheckType         *string           `json:"checkType,omitempty"`
	PresenceStrs      []string          `json:"presenceStrs,omitempty"`
	PresenseStrs      []string          `json:"presenseStrs,omitempty"` // Note: keeping the typo as it exists in data
	AbsenceStrs       []string          `json:"absenceStrs,omitempty"`
	RegexCheck        *string           `json:"regexCheck,omitempty"`
	Errors            map[string]string `json:"errors,omitempty"`
	Disabled          *bool             `json:"disabled,omitempty"`
}

func main() {
	// Check command line arguments
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input_json_file> <output_go_file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s ../maigret_fix.json ../maigretSites.go\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Read the input JSON file
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading %s: %v", inputFile, err)
	}

	// Parse the JSON data
	var sites map[string]SiteConfig
	if err := json.Unmarshal(data, &sites); err != nil {
		log.Fatalf("Error parsing JSON from %s: %v", inputFile, err)
	}

	// Generate the Go code
	var output strings.Builder

	// Package declaration and imports
	output.WriteString("// Code generated from JSON file; DO NOT EDIT.\n")
	output.WriteString("package sites\n\n")

	// Generate the struct definition
	output.WriteString("// MaigretSiteData represents configuration for a single site\n")
	output.WriteString("type MaigretSiteData struct {\n")
	output.WriteString("    Tags              []string          `json:\"tags,omitempty\"`\n")
	output.WriteString("    Engine            *string           `json:\"engine,omitempty\"`\n")
	output.WriteString("    AlexaRank         *int              `json:\"alexaRank,omitempty\"`\n")
	output.WriteString("    URLMain           *string           `json:\"urlMain,omitempty\"`\n")
	output.WriteString("    URL               *string           `json:\"url,omitempty\"`\n")
	output.WriteString("    URLSubpath        *string           `json:\"urlSubpath,omitempty\"`\n")
	output.WriteString("    URLProbe          *string           `json:\"urlProbe,omitempty\"`\n")
	output.WriteString("    UsernameClaimed   *string           `json:\"usernameClaimed,omitempty\"`\n")
	output.WriteString("    UsernameUnclaimed *string           `json:\"usernameUnclaimed,omitempty\"`\n")
	output.WriteString("    CheckType         *string           `json:\"checkType,omitempty\"`\n")
	output.WriteString("    PresenceStrs      []string          `json:\"presenceStrs,omitempty\"`\n")
	output.WriteString("    PresenseStrs      []string          `json:\"presenseStrs,omitempty\"` // Note: typo in original data\n")
	output.WriteString("    AbsenceStrs       []string          `json:\"absenceStrs,omitempty\"`\n")
	output.WriteString("    RegexCheck        *string           `json:\"regexCheck,omitempty\"`\n")
	output.WriteString("    Errors            map[string]string `json:\"errors,omitempty\"`\n")
	output.WriteString("    Disabled          *bool             `json:\"disabled,omitempty\"`\n")
	output.WriteString("}\n\n")

	// Generate the map variable
	output.WriteString("// MaigretSites maps site names to their data\n")
	output.WriteString("var MaigretSites = map[string]MaigretSiteData{\n")

	// Add each site to the map
	for siteName, config := range sites {
		output.WriteString(fmt.Sprintf("    \"%s\": {\n", siteName))

		// Add non-nil fields
		if len(config.Tags) > 0 {
			output.WriteString("        Tags: []string{\n")
			for _, tag := range config.Tags {
				output.WriteString(fmt.Sprintf("            \"%s\",\n", tag))
			}
			output.WriteString("        },\n")
		}

		if config.Engine != nil {
			output.WriteString(fmt.Sprintf("        Engine: stringPtr(\"%s\"),\n", *config.Engine))
		}

		if config.AlexaRank != nil {
			output.WriteString(fmt.Sprintf("        AlexaRank: intPtr(%d),\n", *config.AlexaRank))
		}

		if config.URLMain != nil {
			output.WriteString(fmt.Sprintf("        URLMain: stringPtr(\"%s\"),\n", *config.URLMain))
		}

		if config.URL != nil {
			output.WriteString(fmt.Sprintf("        URL: stringPtr(\"%s\"),\n", *config.URL))
		}

		if config.URLSubpath != nil {
			output.WriteString(fmt.Sprintf("        URLSubpath: stringPtr(\"%s\"),\n", *config.URLSubpath))
		}

		if config.URLProbe != nil {
			output.WriteString(fmt.Sprintf("        URLProbe: stringPtr(\"%s\"),\n", *config.URLProbe))
		}

		if config.UsernameClaimed != nil {
			output.WriteString(fmt.Sprintf("        UsernameClaimed: stringPtr(\"%s\"),\n", *config.UsernameClaimed))
		}

		if config.UsernameUnclaimed != nil {
			output.WriteString(fmt.Sprintf("        UsernameUnclaimed: stringPtr(\"%s\"),\n", *config.UsernameUnclaimed))
		}

		if config.CheckType != nil {
			output.WriteString(fmt.Sprintf("        CheckType: stringPtr(\"%s\"),\n", *config.CheckType))
		}

		if len(config.PresenceStrs) > 0 {
			output.WriteString("        PresenceStrs: []string{\n")
			for _, str := range config.PresenceStrs {
				output.WriteString(fmt.Sprintf("            \"%s\",\n", strings.ReplaceAll(str, "\"", "\\\"")))
			}
			output.WriteString("        },\n")
		}

		if len(config.PresenseStrs) > 0 {
			output.WriteString("        PresenseStrs: []string{\n")
			for _, str := range config.PresenseStrs {
				output.WriteString(fmt.Sprintf("            \"%s\",\n", strings.ReplaceAll(str, "\"", "\\\"")))
			}
			output.WriteString("        },\n")
		}

		if len(config.AbsenceStrs) > 0 {
			output.WriteString("        AbsenceStrs: []string{\n")
			for _, str := range config.AbsenceStrs {
				output.WriteString(fmt.Sprintf("            \"%s\",\n", strings.ReplaceAll(str, "\"", "\\\"")))
			}
			output.WriteString("        },\n")
		}

		if config.RegexCheck != nil {
			output.WriteString(fmt.Sprintf("        RegexCheck: stringPtr(\"%s\"),\n", strings.ReplaceAll(*config.RegexCheck, "\"", "\\\"")))
		}

		if len(config.Errors) > 0 {
			output.WriteString("        Errors: map[string]string{\n")
			for key, value := range config.Errors {
				output.WriteString(fmt.Sprintf("            \"%s\": \"%s\",\n", key, strings.ReplaceAll(value, "\"", "\\\"")))
			}
			output.WriteString("        },\n")
		}

		if config.Disabled != nil {
			output.WriteString(fmt.Sprintf("        Disabled: boolPtr(%t),\n", *config.Disabled))
		}

		output.WriteString("    },\n")
	}

	output.WriteString("}\n\n")

	// Add helper functions
	output.WriteString("// Helper functions for creating pointers\n")
	output.WriteString("func stringPtr(s string) *string {\n")
	output.WriteString("    return &s\n")
	output.WriteString("}\n\n")
	output.WriteString("func intPtr(i int) *int {\n")
	output.WriteString("    return &i\n")
	output.WriteString("}\n\n")
	output.WriteString("func boolPtr(b bool) *bool {\n")
	output.WriteString("    return &b\n")
	output.WriteString("}\n")

	// Write the output to file
	err = ioutil.WriteFile(outputFile, []byte(output.String()), 0644)
	if err != nil {
		log.Fatalf("Error writing to %s: %v", outputFile, err)
	}

	// Generate statistics
	enabled := 0
	disabled := 0
	for _, site := range sites {
		if site.Disabled != nil && *site.Disabled {
			disabled++
		} else {
			enabled++
		}
	}

	fmt.Printf("Successfully generated Go map from %s to %s\n", inputFile, outputFile)
	fmt.Printf("Statistics:\n")
	fmt.Printf("  Total sites: %d\n", len(sites))
	fmt.Printf("  Enabled sites: %d\n", enabled)
	fmt.Printf("  Disabled sites: %d\n", disabled)
}
