package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
		log.Printf("Site %s has %s in stock!", newSite.Domain(), newSite.Product)
		// TODO: notify
	}

	// lastly update our cache
	AddSiteToCache(newSite)
}

func worker(urlQueue *chan string) {
	// take a URL from queue
	for url := range *urlQueue {
		if verboseLogging {
			log.Printf("Processing URL=%s", url)
		}

		// figure out which scraper to use, if we do not have one we can just continue
		scraper, err := GetScraperByURL(url)
		if err != nil {
			log.Printf("Failed to find a Scraper for url '%s'", url)
			continue
		}

		// If we are doing a dryrun, just stop here
		if dryrun {
			continue
		}

		// get website content
		html := DownloadWebContent(url)

		scrapeSite(scraper, url, html)

	}
}
