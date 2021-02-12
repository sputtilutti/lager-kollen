package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type Scraper interface {
	ScrapeSite(url string, doc *goquery.Document) (Website, error)
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

func (s *powerScraper) ScrapeSite(url string, doc *goquery.Document) (Website, error) {
	productIntro := doc.Find("#product-intro")
	if productIntro == nil {
		return Website{}, fmt.Errorf("Could not find id '#product-intro' on URL '%s'", url)
	}
	productName := productIntro.Find("h1").Text()

	mainCard := doc.Find(".product-main-card")
	if mainCard == nil {
		return Website{}, fmt.Errorf("Could not find class '.product-main-card' on URL '%s'", url)
	}

	status := mainCard.Find("pwr-product-stock-label").Text()
	hasItemInStock := status != "Inte i lager"

	return Website{
		URL:            url,
		Product:        productName,
		LastStatusText: status,
		HasItemInStock: hasItemInStock,
		LastVisit:      NowTimeFormatted()}, nil
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

func (s *komplettScraper) ScrapeSite(url string, doc *goquery.Document) (Website, error) {
	return Website{}, nil
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
