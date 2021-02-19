package main

import (
	"fmt"
	"sort"
	"sync"
)

var cache = make(map[string]*Website) // Cache keep track of the webpages visited
var cacheLock = &sync.Mutex{}

// Website represents a page that the Scrapper has visisted
type Website struct {
	URL            string
	product        string
	lastScraped    string
	lastStatusText string
	hasItemInStock bool
	sync.Mutex
}

// Sorting of Website
type ByURL []*Website

func (a ByURL) Len() int           { return len(a) }
func (a ByURL) Less(i, j int) bool { return a[i].URL < a[j].URL }
func (a ByURL) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Instantiate a empty Website struct
func NewSite(url string) (*Website, error) {
	if !IsValidUrl(url) {
		return nil, fmt.Errorf("'%s' is not a valid URL", url)
	}
	return &Website{URL: url}, nil
}

/*
// Get Website from cache.
// If it does not exist, return a new empy Website struct (same as NewSite())
func GetSiteFromCache(url string) (*Website, error) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	ws, exists := cache[url]
	if exists {
		return ws, nil
	}

	ws, err := NewSite(url)
	if err != nil {
		return nil, err
	}
	cache[url] = ws
	return ws, nil
}
*/

func GetAllSitesFromCache() []*Website {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	sites := make([]*Website, 0, len(cache))
	for _, s := range cache {
		sites = append(sites, s)
	}
	sort.Sort(ByURL(sites))
	return sites
}

// RemoveSiteFromCache removes a site from cache
func RemoveSiteFromCache(url string) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	delete(cache, url)
}

/*
// IsSiteInCache checks if site exist in cache
func IsSiteInCache(url string) bool {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	_, exists := cache[url]
	return exists
}
*/

// AddSiteToCache adds a Website to cache, or updates existing entry
func AddSiteToCache(s *Website) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	cache[s.URL] = s
}

func (w *Website) ToString() string {
	w.Lock()
	defer w.Unlock()
	return fmt.Sprintf("url=%s, product=%s, lastStatus=%s, inStock=%v",
		w.URL, w.product, w.lastStatusText, w.hasItemInStock)
}

func (w *Website) Domain() string {
	return DomainFromURL(w.URL)
}

func (w *Website) HasItemInStock() bool {
	w.Lock()
	defer w.Unlock()
	return w.hasItemInStock
}

func (w *Website) IsScraped() bool {
	w.Lock()
	defer w.Unlock()
	return w.lastScraped != ""
}

func (w *Website) LastScraped() string {
	w.Lock()
	defer w.Unlock()
	return w.lastScraped
}

func (w *Website) SetProduct(p string) {
	w.Lock()
	defer w.Unlock()
	w.product = p
}

func (w *Website) Product() string {
	w.Lock()
	defer w.Unlock()
	return w.product
}

func (w *Website) Update(statusText string, hasItemInStock bool) {
	w.Lock()
	defer w.Unlock()
	w.lastStatusText = statusText
	w.hasItemInStock = hasItemInStock
	w.lastScraped = NowTimeFormatted()
}
