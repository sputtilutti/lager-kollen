package main

import (
	"io"
	"log"
	"net/http"
)

func debug(w http.ResponseWriter, r *http.Request) {
}

func api(w http.ResponseWriter, r *http.Request) {
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
            
            tr:nth-child(even) {
                background-color: #dddddd;
            }
        </style>
    </head>
    <body>
        <h1>Status of scraped websites</h1>
        <table>
            <tr>
                <th>Site</th>
                <th>Product</th>
                <th>Last Visit</th>
                <th>Has item in stock?</th>
            </tr>
    `
	for _, ws := range GetAllWebsites() {
		html += "<tr>"
		html += "<td>" + DomainFromURL(ws.URL) + "</td>"
		html += "<td><a href=\"" + ws.URL + "\">" + ws.Product + "</a></td>"
		html += "<td>" + ws.LastVisit + "</td>"
		if ws.HasItemInStock {
			html += "<td style=\"background-color:MediumSeaGreen;\">YES!!</td>"
		} else {
			html += "<td style=\"background-color:Tomato;\">no</td>"
		}
		html += "</tr>"
	}
	html += "</table>"
	html += "</body></html>"
	_, err := io.WriteString(w, html)
	if err != nil {
		log.Println("Failed to respond to HTTP request", err)
	}
}

func createWebServer(listenAddress string) {
	http.HandleFunc("/", overview)
	http.HandleFunc("/api", api)
	http.HandleFunc("/debug", debug)
	log.Println("Server started, listening on", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
