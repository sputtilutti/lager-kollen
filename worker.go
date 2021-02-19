package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var wChannels = make(map[string]chan string, 1)
var wChannelsLock = &sync.Mutex{}
var idx = 0

func StartScrapingURL(url string) error {
	wChannelsLock.Lock()
	defer wChannelsLock.Unlock()

	_, exists := wChannels[url]
	if exists {
		return fmt.Errorf("Cannot start scraping URL '%s'. It is already being monitored/scraped.", url)
	}

	site, err := NewSite(url)
	if err != nil {
		return err
	}

	AddSiteToCache(site)

	// make a channel to communicate with the worker
	ch := make(chan string)
	wChannels[url] = ch
	idx++
	go worker(idx, site, ch)
	return nil
}

func StopScrapingURL(url string) error {
	wChannelsLock.Lock()
	defer wChannelsLock.Unlock()

	_, exists := wChannels[url]
	if !exists {
		return fmt.Errorf("Cannot stop scraping URL '%s'. It is not being monitored/scraped.", url)
	}
	sendStop(wChannels[url])
	RemoveSiteFromCache(url)
	delete(wChannels, url)
	return nil
}

/*
func StopAllScrapers() {
	wChannelsLock.Lock()
	defer wChannelsLock.Unlock()
	for _, ch := range wChannels {
		sendStop(ch)
	}
}
*/

func sendStop(ch chan string) {
	ch <- "stop"
}

func scrapeSite(s Scraper, ws *Website, html string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	hadItemInStock := ws.HasItemInStock()

	// do the scraping
	err = s.ScrapeSite(ws, doc)
	if err != nil {
		log.Printf("Failed to scrape url '%s'. %s", ws.URL, err)
		return
	}

	// if we have the item in stock now, but we did not have it in stock before; send notification
	if ws.HasItemInStock() && !hadItemInStock {
		log.Printf("Site %s has %s in stock!", ws.Domain(), ws.Product())
		// TODO: notify
	}
}

func worker(id int, site *Website, msgs chan string) {
	// take a site from queue

	ticker := time.NewTicker(time.Duration(Config.pollInterval) * time.Second)

	for {
		select {
		case msg := <-msgs:
			if msg == "stop" {
				log.Printf("[%d] Requested to stop monitoring URL '%s'", id, site.URL)
				return
			} else {
				log.Printf("[%d] Worker received msg: %s. Not sure what to do with that.", id, msg)
			}
		case <-ticker.C:
			if Config.verboseLogging {
				log.Printf("[%d] Processing URL=%s", id, site.URL)
			}

			// figure out which scraper to use, if we do not have one we can just continue
			scraper, err := GetScraperByURL(site.URL)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			// If we are doing a dryrun, just stop here
			if Config.dryrun {
				continue
			}

			// get website content
			html := DownloadWebContent(site.URL)

			scrapeSite(scraper, site, html)

			if Config.verboseLogging {
				log.Printf("[%d] Scraped: %v", id, site.ToString())
			}

		}
	}
}
