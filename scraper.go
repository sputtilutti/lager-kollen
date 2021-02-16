package main

import (
	"fmt"

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
		productIntro := doc.Find("#product-intro")
		if productIntro == nil {
			return fmt.Errorf("Could not find id '#product-intro' on '%s'", site.URL)
		}
		productName := productIntro.Find("h1").Text()
		site.SetProduct(productName)
	}

	mainCard := doc.Find(".product-main-card")
	if mainCard == nil {
		return fmt.Errorf("Could not find id '.product-main-card' on '%s'", site.URL)
	}

	status := mainCard.Find("pwr-product-stock-label").Text()
	hasItemInStock := status != "Inte i lager"
	site.Update(status, hasItemInStock)
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
