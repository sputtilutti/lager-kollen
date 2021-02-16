package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

var cache = make(map[string]*Website) // Cache keep track of the webpages visited
var cacheLock = &sync.Mutex{}

// Website represents a page that the Scrapper has visisted
type Website struct {
	URL            string `json:"url"`
	Product        string `json:"product"`
	LastScraped    string `json:"lastScraped"`
	LastStatusText string `json:"lastStatusText"`
	HasItemInStock bool   `json:"hasItemInStock"`
	sync.Mutex
}

// Sorting of Website
type ByURL []*Website

func (a ByURL) Len() int           { return len(a) }
func (a ByURL) Less(i, j int) bool { return a[i].URL < a[j].URL }
func (a ByURL) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func NewSite(url string) (Website, error) {
	if !IsValidUrl(url) {
		return Website{}, fmt.Errorf("'%s' is not a valid URL", url)
	}
	return Website{URL: url}, nil
}

/*
func GetSiteFromCache(url string) Website {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	return cache[url]
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
	out, err := json.Marshal(w)
	if err != nil {
		return ""
	}
	return string(out)
}

func (w *Website) Domain() string {
	return DomainFromURL(w.URL)
}

func (w *Website) IsScraped() bool {
	w.Lock()
	defer w.Unlock()
	return w.LastScraped != ""
}

func (w *Website) SetProduct(p string) {
	w.Lock()
	defer w.Unlock()
	w.Product = p
}

func (w *Website) Update(statusText string, hasItemInStock bool) {
	w.Lock()
	defer w.Unlock()
	w.LastStatusText = statusText
	w.HasItemInStock = hasItemInStock
	w.LastScraped = NowTimeFormatted()
}
