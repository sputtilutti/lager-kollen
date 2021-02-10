package main

import (
	"log"
	"os"
	"os/exec"
)

var wrapperScriptPath = "/tmp/phantomjs-wrapper.js"

var wrapperScriptContent = `
var system = require('system');

url = system.args[1];

var page = require("webpage").create();
page.open(url, function () {
    console.log(page.content);
    phantom.exit();
});
`

func isPathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// https://stackoverflow.com/questions/60128401/how-to-check-if-a-file-is-executable-in-go
func isPathExecutable(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode()&0111 != 0

}

// Get path to phantomJS wrapper script, which is the script that actually does the download
// of webpage content and writes it to console. The wrapper script is executed by phantomjs binary.
// If the script does not exist, it will be generated
func getPhantomJSwrapperPath() string {
	if !isPathExists(wrapperScriptPath) {
		// wrapper script does not exist, create it
		fout, err := os.Create(wrapperScriptPath)
		if err != nil {
			log.Fatalf("Failed to create phantomJS script (%s). %s", wrapperScriptPath, err)
		}
		defer fout.Close()
		_, err = fout.WriteString(wrapperScriptContent)
		if err != nil {
			log.Fatalf("Failed to create phantomJS script (%s). %s", wrapperScriptPath, err)
		}
	}
	return wrapperScriptPath
}

// Get path to phantomJS binary. Die if not found.
// Will try to find it according to the following
//  os.Getenv("PHANTOMJS_BIN")
//  "/usr/local/bin/phantomjs"
//  "./phantomjs"

func getPhantomJSbinPath() string {
	paths := []string{os.Getenv("PHANTOMJS_BIN"), "/usr/local/bin/phantomjs", "./phantomjs"}

	for _, p := range paths {
		if isPathExists(p) && isPathExecutable(p) {
			return p
		}
	}
	log.Fatalln("Failed to find phantomJS binary")
	return ""
}

// Download the HTML content of the page specified with URL
// This function will execute a Javascript script that uses phantomjs library to download the content.
// I could not figure out how to load and download dynamic html with only Go (without browser client and stuff)
func DownloadWebContent(url string) string {
	bin := getPhantomJSbinPath()
	wrapperScript := getPhantomJSwrapperPath()

	out, err := exec.Command(bin, wrapperScript, url).Output()
	if err != nil {
		log.Fatalln("Failed to exeucute phantomJS binary", err)
	}
	return string(out)
}
