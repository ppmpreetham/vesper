package utils

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// )

// // WhatsmynameSiteData represents a single site definition from whatsmyname_fix.json
// type WhatsmynameSiteData struct {
// 	Name         string            `json:"name"`
// 	URICheck     string            `json:"uri_check"`
// 	URIPretty    string            `json:"uri_pretty,omitempty"`
// 	PostBody     string            `json:"post_body,omitempty"`
// 	ECode        int               `json:"e_code"`
// 	EString      string            `json:"e_string"`
// 	MString      string            `json:"m_string"`
// 	MCode        int               `json:"m_code"`
// 	Known        []string          `json:"known"`
// 	Cat          string            `json:"cat"`
// 	StripBadChar string            `json:"strip_bad_char,omitempty"`
// 	Headers      map[string]string `json:"headers,omitempty"`
// 	Protection   []string          `json:"protection,omitempty"`
// }

// // WhatsmynameData represents the root structure of whatsmyname_fix.json
// type WhatsmynameData struct {
// 	License    []string              `json:"license"`
// 	Authors    []string              `json:"authors"`
// 	Categories []string              `json:"categories"`
// 	Sites      []WhatsmynameSiteData `json:"sites"`
// }

// func main() {
// 	if len(os.Args) < 2 {
// 		fmt.Println("Usage: go run generateWhatsmyname.go whatsmyname_fix.json > sites/whatsmynameSites.go")
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
// 	var whatsmynameData WhatsmynameData
// 	if err := json.Unmarshal(jsonData, &whatsmynameData); err != nil {
// 		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
// 		return
// 	}

// 	// Generate Go code
// 	fmt.Println("// Code generated from JSON file; DO NOT EDIT.")
// 	fmt.Println("package sites")
// 	fmt.Println("")
// 	fmt.Println("// WhatsmynameSites contains all site data for username checking")
// 	fmt.Println("var WhatsmynameSites = []WhatsmynameSiteData{")

// 	// Process each site in the array
// 	for i, site := range whatsmynameData.Sites {
// 		fmt.Printf("\t{")

// 		// Only add fields that have values
// 		fieldsAdded := false

// 		// Name field
// 		if site.Name != "" {
// 			fmt.Printf("Name: %q", site.Name)
// 			fieldsAdded = true
// 		}

// 		// URICheck field
// 		if site.URICheck != "" {
// 			if fieldsAdded {
// 				fmt.Print(", ")
// 			}
// 			fmt.Printf("URICheck: %q", site.URICheck)
// 			fieldsAdded = true
// 		}

// 		// URIPretty field (optional)
// 		if site.URIPretty != "" {
// 			if fieldsAdded {
// 				fmt.Print(", ")
// 			}
// 			fmt.Printf("URIPretty: %q", site.URIPretty)
// 			fieldsAdded = true
// 		}

// 		// PostBody field (optional)
// 		if site.PostBody != "" {
// 			if fieldsAdded {
// 				fmt.Print(", ")
// 			}
// 			fmt.Printf("PostBody: %q", site.PostBody)
// 			fieldsAdded = true
// 		}

// 		// ECode field
// 		if site.ECode != 0 {
// 			if fieldsAdded {
// 				fmt.Print(", ")
// 			}
// 			fmt.Printf("ECode: %d", site.ECode)
// 			fieldsAdded = true
// 		}

// 		// EString field
// 		if site.EString != "" {
// 			if fieldsAdded {
// 				fmt.Print(", ")
// 			}
// 			fmt.Printf("EString: %q", site.EString)
// 			fieldsAdded = true
// 		}

// 		// MString field
// 		if site.MString != "" {
// 			if fieldsAdded {
// 				fmt.Print(", ")
// 			}
// 			fmt.Printf("MString: %q", site.MString)
// 			fieldsAdded = true
// 		}

// 		// MCode field
// 		if site.MCode != 0 {
// 			if fieldsAdded {
// 				fmt.Print(", ")
// 			}
// 			fmt.Printf("MCode: %d", site.MCode)
// 			fieldsAdded = true
// 		}

// 		// Known field
// 		if len(site.Known) > 0 {
// 			if fieldsAdded {
// 				fmt.Print(", ")
// 			}
// 			fmt.Print("Known: []string{")
// 			for j, known := range site.Known {
// 				fmt.Printf("%q", known)
// 				if j < len(site.Known)-1 {
// 					fmt.Print(", ")
// 				}
// 			}
// 			fmt.Print("}")
// 			fieldsAdded = true
// 		}

// 		// Cat field
// 		if site.Cat != "" {
// 			if fieldsAdded {
// 				fmt.Print(", ")
// 			}
// 			fmt.Printf("Cat: %q", site.Cat)
// 			fieldsAdded = true
// 		}

// 		// StripBadChar field (optional)
// 		if site.StripBadChar != "" {
// 			if fieldsAdded {
// 				fmt.Print(", ")
// 			}
// 			fmt.Printf("StripBadChar: %q", site.StripBadChar)
// 			fieldsAdded = true
// 		}

// 		// Headers field (optional)
// 		if len(site.Headers) > 0 {
// 			if fieldsAdded {
// 				fmt.Print(", ")
// 			}
// 			fmt.Print("Headers: map[string]string{")
// 			headerCount := 0
// 			for key, value := range site.Headers {
// 				fmt.Printf("%q: %q", key, value)
// 				headerCount++
// 				if headerCount < len(site.Headers) {
// 					fmt.Print(", ")
// 				}
// 			}
// 			fmt.Print("}")
// 			fieldsAdded = true
// 		}

// 		// Protection field (optional)
// 		if len(site.Protection) > 0 {
// 			if fieldsAdded {
// 				fmt.Print(", ")
// 			}
// 			fmt.Print("Protection: []string{")
// 			for j, protection := range site.Protection {
// 				fmt.Printf("%q", protection)
// 				if j < len(site.Protection)-1 {
// 					fmt.Print(", ")
// 				}
// 			}
// 			fmt.Print("}")
// 			fieldsAdded = true
// 		}

// 		// Add comma after each entry except for the last one
// 		if i < len(whatsmynameData.Sites)-1 {
// 			fmt.Println("},")
// 		} else {
// 			fmt.Println("}")
// 		}
// 	}

// 	fmt.Println("}")

// 	// Also include a type definition for WhatsmynameSiteData
// 	fmt.Println("")
// 	fmt.Println("// WhatsmynameSiteData represents a single site definition from whatsmyname_fix.json")
// 	fmt.Println("type WhatsmynameSiteData struct {")
// 	fmt.Println("\tName         string            `json:\"name\"`")
// 	fmt.Println("\tURICheck     string            `json:\"uri_check\"`")
// 	fmt.Println("\tURIPretty    string            `json:\"uri_pretty,omitempty\"`")
// 	fmt.Println("\tPostBody     string            `json:\"post_body,omitempty\"`")
// 	fmt.Println("\tECode        int               `json:\"e_code\"`")
// 	fmt.Println("\tEString      string            `json:\"e_string\"`")
// 	fmt.Println("\tMString      string            `json:\"m_string\"`")
// 	fmt.Println("\tMCode        int               `json:\"m_code\"`")
// 	fmt.Println("\tKnown        []string          `json:\"known\"`")
// 	fmt.Println("\tCat          string            `json:\"cat\"`")
// 	fmt.Println("\tStripBadChar string            `json:\"strip_bad_char,omitempty\"`")
// 	fmt.Println("\tHeaders      map[string]string `json:\"headers,omitempty\"`")
// 	fmt.Println("\tProtection   []string          `json:\"protection,omitempty\"`")
// 	fmt.Println("}")
// }
