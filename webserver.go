package main

import (
	"io"
	"log"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] (%s) %s", r.Method, r.RemoteAddr, r.URL)

	_, err := io.WriteString(w, "hello world")
	if err != nil {
		log.Println("Failed to respond to HTTP request", err)
	}
}

func createWebServer(listenAddress string) {
	http.HandleFunc("/", handleRequest)
	log.Println("Server started, listening on", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
