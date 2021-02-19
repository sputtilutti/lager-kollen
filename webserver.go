package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func readCloseToString(rc io.ReadCloser) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(rc)
	return buf.String()
}

// POST /debug/scraper/{domain}
// Method for test/debug.
// POST html content to make the scraper for {domain} scrape the page
func debugScraperPost(w http.ResponseWriter, r *http.Request) {

	v := mux.Vars(r)
	domain := v["domain"]
	scraper, err := GetScraperByDomain(domain)
	if err != nil {
		fmt.Fprintf(w, "Failed to get scraper for domain '%s'. %s", domain, err.Error())
		return
	}
	site, _ := NewSite("http://debug.lager-kollen.se") // any valid URL works

	html := readCloseToString(r.Body)
	defer r.Body.Close()

	scrapeSite(scraper, site, html)
	fmt.Fprintf(w, "Scraped: %s", site.ToString())
}

/*
func getQueryParam(r *http.Request, key string) (string, error) {
	params, provided := r.URL.Query()[key]
	if provided {
		return params[0], nil
	}
	return "", fmt.Errorf("Missing query parameter: '%s'", key)
}
*/

// GET /api/urls/add/{url}
func apiAddURL(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	url := v["url"]

	if err := StartScrapingURL(url); err != nil {
		fmt.Fprintf(w, "Failed to add URL. %s", err.Error())
		return
	}
	fmt.Fprint(w, "URL added")
}

// GET /api/urls/remove/{url}
func apiRemoveURL(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	url := v["url"]

	if err := StopScrapingURL(url); err != nil {
		fmt.Fprintf(w, "Failed to remove URL. %s", err.Error())
	}
	fmt.Fprint(w, "URL removed")
}

// display overview of scraped pages
func overview(w http.ResponseWriter, r *http.Request) {
	html := `<html>
    <head>
        <title>Lager Kollen</title>
        <style>
            table {
                font-family: arial, sans-serif;
                border-collapse: collapse;
                width: 100%;
            }
            
            td, th {
                border: 1px solid #dddddd;
                text-align: left;
                padding: 8px;
            }            
        </style>
    </head>
    <body>
        <h1>Status of scraped websites</h1>
        <table>
            <tr>
                <th>Site</th>
                <th>Product</th>
                <th>Last scraped</th>
                <th>Has item in stock?</th>
            </tr>
    `
	for _, ws := range GetAllSitesFromCache() {
		html += "<tr>"
		html += "<td>" + DomainFromURL(ws.URL) + "</td>"
		html += "<td><a href=\"" + ws.URL + "\">" + ws.Product() + "</a></td>"
		html += "<td>" + ws.LastScraped() + "</td>"
		if ws.HasItemInStock() {
			html += "<td style=\"background-color:MediumSeaGreen;\">YES!!</td>"
		} else {
			html += "<td style=\"background-color:Tomato;\">no</td>"
		}
		html += "</tr>"
	}
	html += "</table>"
	html += "</body></html>"
	fmt.Fprint(w, html)
}

func createWebServer(listenAddress string) {
	r := mux.NewRouter()
	r.HandleFunc("/", overview)

	apiRoute := r.PathPrefix("/api").Subrouter()
	apiRoute.HandleFunc("/urls/add/{url}", apiAddURL).Methods(http.MethodGet)
	apiRoute.HandleFunc("/urls/remove/{url}", apiRemoveURL).Methods(http.MethodGet)

	debugRoute := r.PathPrefix("/debug").Subrouter()
	debugRoute.HandleFunc("/scraper/{domain}", debugScraperPost).Methods(http.MethodPost)

	log.Println("Webserver starting. Listening on", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, r))
}
