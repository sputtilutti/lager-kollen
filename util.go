package main

import (
	neturl "net/url"
	"os"
	"strings"
	"time"
)

func DomainFromURL(url string) string {
	u, err := neturl.Parse(url)
	if err != nil {
		return ""
	}

	return strings.TrimPrefix(u.Hostname(), "www.") // www.power.se -> power.se
}

func IsPathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// https://stackoverflow.com/questions/60128401/how-to-check-if-a-file-is-executable-in-go
func IsPathExecutable(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode()&0111 != 0
}

// https://golangcode.com/how-to-check-if-a-string-is-a-url/
func IsValidUrl(url string) bool {
	_, err := neturl.ParseRequestURI(url)
	if err != nil {
		return false
	}

	u, err := neturl.Parse(url)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func NowTimeFormatted() string {
	return time.Now().Format(time.RFC1123)
}
