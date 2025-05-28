// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// )

// // Custom unmarshal type to handle both string and []string for errorMsg
// type ErrorMessage []string

// func (e *ErrorMessage) UnmarshalJSON(data []byte) error {
// 	// Try to unmarshal as string first
// 	var s string
// 	if err := json.Unmarshal(data, &s); err == nil {
// 		*e = []string{s}
// 		return nil
// 	}

// 	// If that fails, try to unmarshal as []string
// 	var arr []string
// 	if err := json.Unmarshal(data, &arr); err != nil {
// 		return err
// 	}
// 	*e = arr
// 	return nil
// }

// // Modified SherlockSiteData to use our custom ErrorMessage type
// type SherlockSiteData struct {
// 	ErrorMsg        ErrorMessage      `json:"errorMsg,omitempty"`
// 	ErrorType       string            `json:"errorType"`
// 	RegexCheck      string            `json:"regexCheck,omitempty"`
// 	URL             string            `json:"url"`
// 	URLMain         string            `json:"urlMain"`
// 	URLProbe        string            `json:"urlProbe,omitempty"`
// 	UsernameClaimed string            `json:"username_claimed"`
// 	IsNSFW          bool              `json:"isNSFW,omitempty"`
// 	Headers         map[string]string `json:"headers,omitempty"`
// 	RequestMethod   string            `json:"request_method,omitempty"`
// 	RequestPayload  interface{}       `json:"request_payload,omitempty"`
// 	ErrorCode       int               `json:"errorCode,omitempty"`
// 	ErrorURL        string            `json:"errorUrl,omitempty"`
// }

// func main() {
// 	if len(os.Args) < 2 {
// 		fmt.Println("Usage: go run generateSherlockSites.go sherlock_fix.json > sites/sherlockSites.go")
// 		return
// 	}

// 	jsonFilePath := os.Args[1]

// 	// Read the JSON file
// 	jsonData, err := os.ReadFile(jsonFilePath)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error reading JSON file: %v\n", err)
// 		return
// 	}

// 	// Parse JSON into our struct
// 	var sitesMap map[string]SherlockSiteData
// 	if err := json.Unmarshal(jsonData, &sitesMap); err != nil {
// 		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
// 		return
// 	}

// 	// Generate Go code
// 	fmt.Println("// Code generated from JSON file; DO NOT EDIT.")
// 	fmt.Println("package sites")
// 	fmt.Println("")
// 	fmt.Println("// SherlockSites maps site names to their data")
// 	fmt.Println("var SherlockSites = map[string]SherlockSiteData{")

// 	// Create a slice of site names to iterate through
// 	siteNames := make([]string, 0, len(sitesMap))
// 	for siteName := range sitesMap {
// 		siteNames = append(siteNames, siteName)
// 	}

// 	// Process each site in the map
// 	for i, siteName := range siteNames {
// 		siteData := sitesMap[siteName]
// 		fmt.Printf("\t%q: {", siteName)

// 		// Only add fields that have values
// 		fieldsAdded := false

// 		// Handle ErrorMsg array
// 		if len(siteData.ErrorMsg) > 0 {
// 			if fieldsAdded {
// 				fmt.Println(",")
// 			} else {
// 				fmt.Println("")
// 			}
// 			fieldsAdded = true
// 			fmt.Println("\t\tErrorMsg: []string{")
// 			for _, msg := range siteData.ErrorMsg {
// 				fmt.Printf("\t\t\t%q,\n", msg)
// 			}
// 			fmt.Print("\t\t}")
// 		}

// 		// Handle other fields in a similar way
// 		if siteData.ErrorType != "" {
// 			if fieldsAdded {
// 				fmt.Println(",")
// 			} else {
// 				fmt.Println("")
// 			}
// 			fieldsAdded = true
// 			fmt.Printf("\t\tErrorType: %q", siteData.ErrorType)
// 		}

// 		if siteData.RegexCheck != "" {
// 			if fieldsAdded {
// 				fmt.Println(",")
// 			} else {
// 				fmt.Println("")
// 			}
// 			fieldsAdded = true
// 			fmt.Printf("\t\tRegexCheck: %q", siteData.RegexCheck)
// 		}

// 		if siteData.URL != "" {
// 			if fieldsAdded {
// 				fmt.Println(",")
// 			} else {
// 				fmt.Println("")
// 			}
// 			fieldsAdded = true
// 			fmt.Printf("\t\tURL: %q", siteData.URL)
// 		}

// 		if siteData.URLMain != "" {
// 			if fieldsAdded {
// 				fmt.Println(",")
// 			} else {
// 				fmt.Println("")
// 			}
// 			fieldsAdded = true
// 			fmt.Printf("\t\tURLMain: %q", siteData.URLMain)
// 		}

// 		if siteData.URLProbe != "" {
// 			if fieldsAdded {
// 				fmt.Println(",")
// 			} else {
// 				fmt.Println("")
// 			}
// 			fieldsAdded = true
// 			fmt.Printf("\t\tURLProbe: %q", siteData.URLProbe)
// 		}

// 		if siteData.UsernameClaimed != "" {
// 			if fieldsAdded {
// 				fmt.Println(",")
// 			} else {
// 				fmt.Println("")
// 			}
// 			fieldsAdded = true
// 			fmt.Printf("\t\tUsernameClaimed: %q", siteData.UsernameClaimed)
// 		}

// 		if siteData.IsNSFW {
// 			if fieldsAdded {
// 				fmt.Println(",")
// 			} else {
// 				fmt.Println("")
// 			}
// 			fieldsAdded = true
// 			fmt.Print("\t\tIsNSFW: true")
// 		}

// 		if len(siteData.Headers) > 0 {
// 			if fieldsAdded {
// 				fmt.Println(",")
// 			} else {
// 				fmt.Println("")
// 			}
// 			fieldsAdded = true
// 			fmt.Println("\t\tHeaders: map[string]string{")
// 			headerKeys := make([]string, 0, len(siteData.Headers))
// 			for key := range siteData.Headers {
// 				headerKeys = append(headerKeys, key)
// 			}
// 			for j, key := range headerKeys {
// 				value := siteData.Headers[key]
// 				fmt.Printf("\t\t\t%q: %q", key, value)
// 				if j < len(headerKeys)-1 {
// 					fmt.Println(",")
// 				} else {
// 					fmt.Println("")
// 				}
// 			}
// 			fmt.Print("\t\t}")
// 		}

// 		if siteData.RequestMethod != "" {
// 			if fieldsAdded {
// 				fmt.Println(",")
// 			} else {
// 				fmt.Println("")
// 			}
// 			fieldsAdded = true
// 			fmt.Printf("\t\tRequestMethod: %q", siteData.RequestMethod)
// 		}

// 		// Handle RequestPayload if present
// 		if siteData.RequestPayload != nil {
// 			if fieldsAdded {
// 				fmt.Println(",")
// 			} else {
// 				fmt.Println("")
// 			}
// 			fieldsAdded = true

// 			// Convert the interface{} to string representation
// 			payloadBytes, err := json.Marshal(siteData.RequestPayload)
// 			if err == nil {
// 				payloadStr := string(payloadBytes)
// 				fmt.Printf("\t\tRequestPayload: %#v", payloadStr)
// 			}
// 		}

// 		if siteData.ErrorCode != 0 {
// 			if fieldsAdded {
// 				fmt.Println(",")
// 			} else {
// 				fmt.Println("")
// 			}
// 			fieldsAdded = true
// 			fmt.Printf("\t\tErrorCode: %d", siteData.ErrorCode)
// 		}

// 		if siteData.ErrorURL != "" {
// 			if fieldsAdded {
// 				fmt.Println(",")
// 			} else {
// 				fmt.Println("")
// 			}
// 			fieldsAdded = true
// 			fmt.Printf("\t\tErrorURL: %q", siteData.ErrorURL)
// 		}

// 		if fieldsAdded {
// 			fmt.Println("")
// 		}

// 		// IMPORTANT: Always add comma after each map entry, except for the last one
// 		if i < len(siteNames)-1 {
// 			fmt.Println("\t},")
// 		} else {
// 			fmt.Println("\t}")
// 		}
// 	}

// 	fmt.Println("}")
// }
