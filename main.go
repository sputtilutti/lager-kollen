package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var verboseLogging, _ = strconv.ParseBool(os.Getenv("VERBOSE"))

func startWebServer(listenAddress string) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go createWebServer(listenAddress)
	<-sigs
}

func startWorkers(nWorkers int, urlQueue chan string) {
	for w := 1; w <= nWorkers; w++ {
		go worker(&urlQueue)
	}
}

func startURLReader(urlsFile string, pollInterval int, urlQueue chan string) {
	log.Println("Reading URLs from", urlsFile)

	go func() {
		for {
			fin, err := os.Open(urlsFile)
			if err != nil {
				log.Fatal("Could not open urls file for reading")
			}

			scanner := bufio.NewScanner(fin)
			for scanner.Scan() {
				urlQueue <- scanner.Text()
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

			fin.Close()
			time.Sleep(time.Duration(pollInterval) * time.Second)
		}
	}()
}

func main() {
	urlsFile := flag.String("urls", "./urls.csv", "Path to csv with URLs to query (./urls.csv)")
	nThreads := flag.Int("nthreads", 1, "Number of scrapper threads to run in parallel. Each scraper is responsible for downloading web content and parsing it (1)")
	pollInterval := flag.Int("pollInterval", 10, "Frequency in seconds to sleep before next poll (10)")
	listenAddress := flag.String("listenAddress", ":8080", "Listen host:port for webserver (:8080)")
	flag.Parse()

	urlQueue := make(chan string, 10) // URLs to query are put to this queue

	RegisterScraper("power.se", NewPowerScraper)
	RegisterScraper("komplett.se", NewKomplettScraper)

	startWorkers(*nThreads, urlQueue)

	// URL Reader reads URLs from file and sends it to workers
	startURLReader(*urlsFile, *pollInterval, urlQueue)

	// Web server will block until application is stopped by SIGINT or SIGTERM
	startWebServer(*listenAddress)

	log.Println("Stopped")
}
