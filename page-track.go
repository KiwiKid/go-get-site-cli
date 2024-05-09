package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// readURLs reads up to 20 URLs from a file
func readURLs(filePath string, seedUrl string) ([]string, error) {

	log.Printf("filePath: %s", filePath)

	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Printf("Unable to create directory: %v", err)
			return nil, err
		}
	}
	
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("readURLs error: %v", err)
		return nil, err
	}
	defer file.Close()

	// Check if the file is empty by seeking the end to get its size
	if stat, err := file.Stat(); err == nil && stat.Size() == 0 && seedUrl != "" {
		log.Printf("Seeding new file with URL: %s", seedUrl)
		if _, err := file.WriteString(seedUrl + "\n"); err != nil {
			log.Printf("Failed to seed file: %v", err)
			return nil, err
		}
	}

	// Rewind the file pointer to the beginning after seeding
	if _, err := file.Seek(0, 0); err != nil {
		log.Printf("Error seeking file: %v", err)
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		if len(lines) >= 20 {
			break
		}
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// writeURLs writes URLs to a file
func writeURLs(filePath string, urls []string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, url := range urls {
		if _, err := file.WriteString(url + "\n"); err != nil {
			return err
		}
	}
	return nil
}

// moveProcessedURLs moves processed URLs from one file to another
func moveProcessedURLs(srcPath, dstPath string, processedURLs []string) error {
	remainingURLs, err := readURLs(srcPath, "")
	if err != nil {
		return err
	}

	// Write processed URLs to the destination file
	if err := writeURLs(dstPath, processedURLs); err != nil {
		return err
	}

	// Re-write the remaining unprocessed URLs back to the source file
	file, err := os.Create(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, url := range remainingURLs {
		if !strings.Contains(strings.Join(processedURLs, "\n"), url) {
			if _, err := file.WriteString(url + "\n"); err != nil {
				return err
			}
		}
	}
	return nil
}
