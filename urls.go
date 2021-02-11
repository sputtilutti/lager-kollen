package main

import (
	"fmt"
	"sync"
)

var urls = make([]string, 0, 50) // 50 should be enough
var urlLock = &sync.Mutex{}

// Add URL to list of URLs to monitor
// Returns true if URL was added, false if already added or not valid URL.
func AddUrl(url string) (bool, error) {
	if ContainsUrl(url) {
		return false, fmt.Errorf("URL '%s' has already been added", url)
	}
	if !IsValidUrl(url) {
		return false, fmt.Errorf("URL '%s' is not a valid URL", url)
	}
	urlLock.Lock()
	defer urlLock.Unlock()
	urls = append(urls, url)
	return true, nil
}

func RemoveUrl(url string) {
	urlLock.Lock()
	defer urlLock.Unlock()
	for i := 0; i < len(urls); i++ {
		if urls[i] == url {
			urls[i] = ""
		}
	}
}

func ContainsUrl(url string) bool {
	urlLock.Lock()
	defer urlLock.Unlock()
	for _, u := range urls {
		if u == url {
			return true
		}
	}
	return false
}

func UrlSize() int {
	urlLock.Lock()
	defer urlLock.Unlock()
	i := 0
	for _, u := range urls {
		if u != "" {
			i++
		}
	}
	return i
}

func GetUrls() []string {
	urlLock.Lock()
	defer urlLock.Unlock()
	return urls[:]
}
