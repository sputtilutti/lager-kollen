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
var dryrun, _ = strconv.ParseBool(os.Getenv("DRYRUN"))

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

// Read URLs from a file, if file exist.
// Each URL is added to the URL repo/list
func loadUrlsFile(urlsFile string) {
	if !IsPathExists(urlsFile) {
		return
	}

	log.Println("Reading URLs from", urlsFile)

	fin, err := os.Open(urlsFile)
	if err != nil {
		log.Fatalf("Could not open urls file (%s) for reading", urlsFile)
	}
	defer fin.Close()

	scanner := bufio.NewScanner(fin)
	i := 0
	for scanner.Scan() {
		_, err := AddUrl(scanner.Text())
		if err != nil {
			log.Println("Error:", err)
		} else {
			i++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d urls from file", i)
}

// In a go-routine, periodically get the latest list of URLs and
// send each URL to the worker for further processing
func startURLPoller(pollInterval int, urlQueue chan string) {
	go func() {
		for {
			urls := GetUrls()
			for _, url := range urls {
				if url != "" {
					urlQueue <- url
				}
			}
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

	loadUrlsFile(*urlsFile)

	// URL Poller polls URLs and sends it to workers
	startURLPoller(*pollInterval, urlQueue)

	// Web server will block until application is stopped by SIGINT/SIGTERM
	startWebServer(*listenAddress)

	log.Println("Stopped")
}
