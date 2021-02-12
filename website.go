package main

import (
	"encoding/json"
	"sort"
	"sync"
)

var cache = make(map[string]Website) // Cache keep track of the webpages visited
var cacheLock = &sync.Mutex{}

// Website represents a page that the Scrapper has visisted
type Website struct {
	URL            string `json:"url"`
	Product        string `json:"product"`
	LastVisit      string `json:"lastVisit"`
	LastStatusText string `json:"lastStatusText"`
	HasItemInStock bool   `json:"hasItemInStock"`
}

// Sorting of Website
type ByDomain []Website

func (a ByDomain) Len() int           { return len(a) }
func (a ByDomain) Less(i, j int) bool { return DomainFromURL(a[i].URL) < DomainFromURL(a[j].URL) }
func (a ByDomain) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// GetSiteFromCache checks if a site from cache or nil if it does not exist
func GetSiteFromCache(url string) Website {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	return cache[url]
}

func GetAllWebsites() []Website {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	sites := make([]Website, 0, len(cache))

	for _, website := range cache {
		sites = append(sites, website)
	}
	sort.Sort(ByDomain(sites))
	return sites
}

/*
// RemoveSiteFromCache removes a site from cache
func RemoveSiteFromCache(url string) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	delete(cache, url)
}

// IsSiteInCache checks if site exist in cache
func IsSiteInCache(url string) bool {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	_, exists := cache[url]
	return exists
}
*/
// AddSiteToCache adds a Website to cache, or updates existing entry
func AddSiteToCache(s Website) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	cache[s.URL] = s
}

func (w *Website) String() string {
	out, err := json.Marshal(w)
	if err != nil {
		return ""
	}
	return string(out)
}

func (w *Website) Domain() string {
	return DomainFromURL(w.URL)
}
