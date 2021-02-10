package main

import (
	neturl "net/url"
	"strings"
)

func DomainFromURL(url string) string {
	u, err := neturl.Parse(url)
	if err != nil {
		return ""
	}

	return strings.TrimPrefix(u.Hostname(), "www.") // www.power.se -> power.se
}
