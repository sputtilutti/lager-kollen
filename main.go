package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func startWebServer(listenAddress string) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go createWebServer(listenAddress)
	<-sigs
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
		log.Printf("Could not open urls file (%s) for reading", urlsFile)
		return
	}
	defer fin.Close()

	scanner := bufio.NewScanner(fin)
	i := 0
	for scanner.Scan() {
		url := scanner.Text()
		err := StartScrapingURL(url)
		if err != nil {
			log.Printf("Error: %s", err)
			continue
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d urls from file", i)
}

var Config struct {
	pollInterval   int
	verboseLogging bool
	dryrun         bool
}

func main() {
	urlsFile := flag.String("urls", "./urls.csv", "Path to csv with URLs to query (./urls.csv)")
	pollInterval := flag.Int("pollInterval", 10, "Frequency in seconds to sleep before next poll (10)")
	listenAddress := flag.String("listenAddress", ":8080", "Listen host:port for webserver (:8080)")
	flag.Parse()

	//store settings in global config
	Config.pollInterval = *pollInterval
	Config.verboseLogging, _ = strconv.ParseBool(os.Getenv("VERBOSE"))
	Config.dryrun, _ = strconv.ParseBool(os.Getenv("DRYRUN"))

	RegisterScraper("power.se", NewPowerScraper)
	RegisterScraper("komplett.se", NewKomplettScraper)

	// optionally, we can load some URLs from file at start
	loadUrlsFile(*urlsFile)

	// Web server will block until application is stopped by SIGINT/SIGTERM
	startWebServer(*listenAddress)

	log.Println("Stopped")
}
