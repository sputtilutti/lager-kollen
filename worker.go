package main

import (
	"log"
	neturl "net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// https://golangcode.com/how-to-check-if-a-string-is-a-url/
func isValidUrl(url string) bool {
	_, err := neturl.ParseRequestURI(url)
	if err != nil {
		return false
	}

	u, err := neturl.Parse(url)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func scrapeSite(s Scraper, url string, html string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	// get the previous scaped site, will return default (empty) Website struct first time
	oldSite := GetSiteFromCache(url)

	// do the scraping
	newSite, err := s.ScrapeSite(url, doc)
	if err != nil {
		log.Printf("Failed to scrape url '%s'. %s", url, err)
		return
	}

	if verboseLogging {
		log.Printf("Scraped: %v", newSite.String())
	}

	// if we have the item in stock now, but we did not have it in stock before; send notification
	if newSite.HasItemInStock && !oldSite.HasItemInStock {
		log.Printf("Site %s has the item in stock!", newSite.Name())
		// TODO: notify
	}

	// lastly update our cache
	AddSiteToCache(newSite)
}

func worker(urlQueue *chan string) {
	// take a URL from queue
	for url := range *urlQueue {
		if !isValidUrl(url) {
			continue
		}

		if verboseLogging {
			log.Printf("Processing URL=%s", url)
		}

		// figure out which scraper to use, if we do not have one we can just continue
		scraper, err := GetScraperByURL(url)
		if err != nil {
			log.Printf("Failed to find a Scraper for url '%s'", url)
			continue
		}

		// get website content
		html := DownloadWebContent(url)

		scrapeSite(scraper, url, html)

	}
}
