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

func saveDOMToFile(directoryPath string, domContent string, meta PageMetadata, loadTime time.Duration) error {
	// Parse the URL
	log.Printf("saveDOMToFile %s", directoryPath)
	// Define the file names
	timeStamp := time.Now().Format("20060102-150405")
	basePageFileName = meta.Url
	htmlFileName := fmt.Sprintf("%s-%s.html", meta.Url, timeStamp)
	jsonFileName := fmt.Sprintf("output-%s.json", timeStamp)
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
	}else{
		log.Printf("wrote htmlFile %s", htmlFilePath)
	}

	// Write metadata to JSON file
	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		return err
	}else{
		log.Printf("wrote jsonFile %s", jsonFilePath)
	}
	defer jsonFile.Close()

	jsonEncoder := json.NewEncoder(jsonFile)
	jsonEncoder.SetIndent("", "    ")
	err = jsonEncoder.Encode(meta)
	if err != nil {
		return err
	}

	return nil
}
