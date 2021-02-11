package main

import (
	"encoding/json"
	"fmt"
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

// GetSiteFromCache checks if a site from cache or nil if it does not exist
func GetSiteFromCache(url string) Website {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	return cache[url]
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

func (w *Website) Name() string {
	return fmt.Sprintf("%s - %s", w.Domain(), w.Product)
}
