package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func debug(w http.ResponseWriter, r *http.Request) {
}

func getQueryParam(r *http.Request, key string) (string, error) {
	params, provided := r.URL.Query()[key]
	if provided {
		return params[0], nil
	}
	return "", fmt.Errorf("Missing query parameter: '%s'", key)
}

// GET /api?action={action}&url={url}
func api(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		action, err := getQueryParam(r, "action")
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		url, err := getQueryParam(r, "url")
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		res := handleApiAction(action, url)
		fmt.Fprint(w, res)
	} else {
		// default, not found
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleApiAction(action string, url string) string {
	if strings.EqualFold(action, "add") {
		if err := StartScrapingURL(url); err != nil {
			return err.Error()
		}
		return "URL added"
	} else if strings.EqualFold(action, "remove") || strings.EqualFold(action, "delete") {
		if err := StopScrapingURL(url); err != nil {
			return err.Error()
		}
		return "URL removed"
	}

	return fmt.Sprintf("Unknown or unsupported action: '%s'", action)
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
		html += "<td><a href=\"" + ws.URL + "\">" + ws.Product + "</a></td>"
		html += "<td>" + ws.LastScraped + "</td>"
		if ws.HasItemInStock {
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
	http.HandleFunc("/", overview)
	http.HandleFunc("/api", api)
	http.HandleFunc("/debug", debug)
	log.Println("Webserver starting. Listening on", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
