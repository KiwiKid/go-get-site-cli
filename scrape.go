package main

import (
	"context"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type PageMetadata struct {
	Title    string   `json:"title"`
	Url 	url.URL `json:"url"`
	Links    []string `json:"links"`
	LocalSiteLinks []string `json:"local_links"`
	MetaTags []string `json:"meta_tags"`
}


func linkCouldBePage(s string, baseUrl string) bool {
	nonPageExtensions := map[string]bool{
		".css":         true,
		".js":          true,
		".ico":         true,
		".png":         true,
		".jpg":         true,
		".jpeg":        true,
		".gif":         true,
		".bmp":         true,
		".tif":         true,
		".tiff":        true,
		".svg":         true,
		".webp":        true,
		".mp3":         true,
		".wav":         true,
		".mp4":         true,
		".mov":         true,
		".avi":         true,
		".mkv":         true,
		".pdf":         true,
		".xml":         true,
		".txt":         true,
		".less":        true,
		".webmanifest": true,
		".zip":         true,
		".rar":         true,
		".tar":         true,
		".gz":          true,
		".json":        true,
		".woff":        true,
		".woff2":       true,
		".ttf":         true,
		".eot":         true,
		".otf":         true,
		".flv":         true,
		".swf":         true,
		".iso":         true,
	}

	// Parse the href to a URL structure
	u, err := url.Parse(s)
	if err != nil {
		log.Printf("failed to parse:check %s %v\n", s, err)
		return false
	}
	if len(u.String()) == 0 {
		log.Printf("failed to parse:check (empty string) %s\n", s)

		return false
	}
	bu, buErr := url.Parse(baseUrl)
	if buErr != nil {
		return false
	}
	// Extract only the path, ignoring query and fragment
	path := u.EscapedPath()

	log.Printf("linkCouldBePage:check %s\n", s)
	// Check the extension
	for ext := range nonPageExtensions {
		if strings.HasSuffix(strings.ToLower(path), ext) {
			return false
		}
	}

	baseDomain := bu.Host
	pathDomain := u.Host

	if baseDomain == pathDomain || strings.HasPrefix(s, "/") {
		log.Print("linkCouldBePage:VALID\n")
		return true
	} else {
		log.Print("linkCouldBePage:INVALID\n")
		return false
	}
}

func scrapeDOM(urltoGet string, baseUrl string) (string, PageMetadata, time.Duration, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	parsedURL, pErr := url.Parse(urltoGet)
	if pErr != nil {
		return "", PageMetadata{}, time.Duration(0), pErr
	}

	var domContent, title string
	var links, metaTags []string
	var meta PageMetadata
	startTime := time.Now()

	err := chromedp.Run(ctx,
		chromedp.Navigate(urltoGet),
		chromedp.Title(&title),
		chromedp.OuterHTML("html", &domContent),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a[href]')).map(a => a.href)`, &links),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('meta')).map(meta => meta.outerHTML)`, &metaTags),
	)
	loadTime := time.Since(startTime)

	if err != nil {
		return "", PageMetadata{}, loadTime, err
	}

	localSiteLinks := []string{}
	for _, l := range links {
		log.Printf("linkCouldBePage(%v, %v)", l, baseUrl)
		if linkCouldBePage(l, baseUrl) {
			localSiteLinks = append(localSiteLinks, l)
		}
	}

	meta = PageMetadata{
		Title:    title,
		Url: *parsedURL,
		Links:    links,
		LocalSiteLinks: localSiteLinks,
		MetaTags: metaTags,
	}
	return domContent, meta, loadTime, nil
}