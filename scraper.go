package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Scraper interface {
	ScrapeSite(site *Website, doc *goquery.Document) error
	Domain() string
}

// ScrapperFactory for creating new Scrappers
type ScraperFactory func() (Scraper, error)

var scraperFactories = make(map[string]ScraperFactory)

func RegisterScraper(domain string, f ScraperFactory) {
	scraperFactories[domain] = f
}

// ---- Power.se Scraper ----

type powerScraper struct {
}

func (s *powerScraper) ScrapeSite(site *Website, doc *goquery.Document) error {
	if !site.IsScraped() {
		var productName = ""

		doc.Find("pwr-product-page meta").Each(func(index int, item *goquery.Selection) {
			content, _ := item.Attr("content")
			itemprop, _ := item.Attr("itemprop")
			if itemprop == "name" {
				productName = content
			}
		})
		site.SetProduct(strings.ToLower(productName))
	}

	mainCard := doc.Find(".product-main-card")
	if mainCard == nil {
		return fmt.Errorf("Could not find id '.product-main-card' on '%s'", site.URL)
	}

	status := mainCard.Find("pwr-product-stock-label").Text()
	if len(status) == 0 {
		return fmt.Errorf("Failed to scrape status text from %s", site.URL)
	}
	itemNotInStock := strings.EqualFold("inte i lager", status)
	site.Update(status, !itemNotInStock)

	return nil
}

func (s *powerScraper) Domain() string {
	return "power.se"
}

func NewPowerScraper() (Scraper, error) {
	return &powerScraper{}, nil
}

// ---- Komplett.se Scraper ----

type komplettScraper struct {
}

func (s *komplettScraper) ScrapeSite(site *Website, doc *goquery.Document) error {
	return nil
}

func (s *komplettScraper) Domain() string {
	return "komplett.se"
}

func NewKomplettScraper() (Scraper, error) {
	return &komplettScraper{}, nil
}

// GetScraper tries to figure out which Scraper to use based on the provided URL. Returns nil of we do not have any Scrapper for this URL
func GetScraperByURL(url string) (Scraper, error) {
	return GetScraperByDomain(DomainFromURL(url))
}

func GetScraperByDomain(domain string) (Scraper, error) {
	f, exists := scraperFactories[domain]
	if !exists {
		return nil, fmt.Errorf("No factory defined for domain '%s'", domain)
	}
	return f()
}
