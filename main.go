package main

import (
	"errors"
	"log"
	"strings"

	"github.com/christophberger/start"
)

func main() {
    start.Add(&start.Command{
		Name: "scrape",
		Short: "scrape <url>",
		Long: "Scrape the DOM of the specified URL and save it along with metadata",
		Cmd: scrapeCmd,
	})

	start.Up()
}


func scrapeCmd(cmd *start.Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("failed to get any command line args")
	}
	url := cmd.Args[0]

	realUrl := url

	folder, err := getFolderName(url)
	if err != nil {
		return err
	}

	var unprocessedPath = *folder + "/unprocessed" 
	var processedPath = *folder + "/processed" 

	// Read existing URLs from both unprocessed and processed files
	unprocessedURLs, err := readURLs(unprocessedPath, realUrl)
	if err != nil {
		log.Fatalf("Failed to read unprocessed URLs: %s", err)
	}
	processedURLs, err := readURLs(processedPath, "")
	if err != nil {
		log.Fatalf("Failed to read processed URLs: %s", err)
	}

	// Merge both URL lists to check against
	allExistingURLs := append(unprocessedURLs, processedURLs...)

	for _, url := range unprocessedURLs {
		// Assuming scrapeDOM and saveDOMToFile are defined elsewhere and operate on a single URL
		domContent, meta, loadTime, err := scrapeDOM(url, realUrl)
		if err != nil {
			return err
		}
		log.Printf("Scraped %s, saving data", url)
		err = saveDOMToFile(*folder, domContent, meta, loadTime)
		if err != nil {
			return err
		}

		// This should handle moving the processed URL after successful scraping
		if err := moveProcessedURLs(unprocessedPath, processedPath, []string{url}); err != nil {
			return err
		}
		
		log.Printf("found %d", len(meta.LocalSiteLinks))
		for _, newURL := range meta.LocalSiteLinks {
			
			if !contains(allExistingURLs, newURL) {
				log.Printf("writing new url %s", newURL)
				if err := writeURLs(unprocessedPath, []string{newURL}); err != nil {
					return err
				}else{
					log.Printf("wrote new url to unprocessed %s", newURL)
				}
			}else{
				log.Printf("url already processed %s", newURL)
			}
		}
	}
	return nil
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == strings.TrimSpace(item) {
			return true
		}
	}
	return false
}