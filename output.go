package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getFolderName(urlStr string) (*string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// Extract the base domain
	baseDomain := parsedURL.Hostname()

	// Normalize the path to use as a folder name
	path := strings.Trim(parsedURL.Path, "/")
	/*if path == "" {
		path = "index"
	} else {*/
		path = strings.ReplaceAll(path, "/", "_")
	//}

	baseDir := "output"
	directoryPath := filepath.Join(baseDir, baseDomain, path)
	err = os.MkdirAll(directoryPath, 0755)
	if err != nil {
		return nil, err
	}
	log.Printf("getFolderName: %v", directoryPath)
	return &directoryPath, nil
}
func saveDOMToFile(urlStr string, baseDirectory string, domContent string, meta PageMetadata, loadTime time.Duration) error {
	log.Printf("saveDOMToFile %s", urlStr)

	// Parse the URL to extract path segments
	u, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	// Normalize path and split into segments
	pathSegments := strings.Split(strings.Trim(u.Path, "/"), "/")
	var segs = []string{ baseDirectory}
	for _, ps := range pathSegments {
		segs = append(segs, ps)
	}
	directoryPath := filepath.Join(segs...)

	// Create directory structure
	if err := os.MkdirAll(directoryPath, 0755); err != nil {
		return err
	}

	// Determine file names based on whether it's a directory or specific file
	var htmlFileName, jsonFileName string
	if filepath.Ext(u.Path) == "" {
		// No file extension in the URL, assume directory and use 'index.html'
		htmlFileName = "index.html"
	} else {
		// Specific file
		htmlFileName = filepath.Base(u.Path)
	}

	// Timestamp for uniqueness in JSON file name
	timeStamp := time.Now().Format("20060102-150405")
	jsonFileName = fmt.Sprintf("output-%s.json", timeStamp)

	// Full paths for the files
	htmlFilePath := filepath.Join(directoryPath, htmlFileName)
	jsonFilePath := filepath.Join(directoryPath, jsonFileName)

	// Create and write the HTML content
	htmlFile, err := os.Create(htmlFilePath)
	if err != nil {
		return err
	}
	defer htmlFile.Close()

	_, err = htmlFile.WriteString(domContent)
	if err != nil {
		return err
	} else {
		log.Printf("Wrote HTML file %s", htmlFilePath)
	}

	// Write metadata to JSON file
	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonEncoder := json.NewEncoder(jsonFile)
	jsonEncoder.SetIndent("", "    ")
	err = jsonEncoder.Encode(meta)
	if err != nil {
		return err
	}

	log.Printf("Wrote JSON file %s", jsonFilePath)
	return nil
}
