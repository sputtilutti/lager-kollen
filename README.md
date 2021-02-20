# Lager Kollen

A solution to monitor specific item's stock/inventory status on selected websites in order to get alert when the item is in stock.

Note: by 'solution' I actually mean Proof of Concept solution - the code works but does not provide a full solution (e.g. no notification). See background below.

# Background

I pre-ordered a Playstation 5 back in May 2020 and due to the insane demand of the product together with a low distribution rate from Sony, I realised I probably had to wait a long time to get my PS5. The distributors (in Sweden) seemed to get a delivery roughly every month and some of them kept a couple of items to be sold directly on their webpage and not via a pre-ordering system. Problem is that as soon as they announced it in stock, it was sold out within minutes. There are sites that allow you to track inventory/stock status, but the free versions are fairly limited (to nbr pages and refresh rate) and I did not want to pay for premium - so I figured I could develop my own solution that can run on a box and notify me when the item is in stock.
Also gave me an opportunity to practice coding Golang.

Enter `lager-kollen`, a Go package that periodically downloads the webpage content of selected URLs, parses it and tries to find a text/string to figure out if the item is in stock or not. If found, it will send an email to notify that website X has item Y in store.

I started coding on this in a fairly slow pace, figured I had some time before next batch would arrive to the distributors. After completing a first draft version, I was notififed that the place I pre-ordered by PS5 from (Webhallen.se) expect to have my order by March 1 2021 !! 
So, with that I do not feel the need to complete the ideas I had for this project. The code available works and serves as a POC for how it could be done, but there are a bunch of improvements that I wanted to add, such as add a nice REST API to poll status, add new URLs, etc.

# Some notes on architecture

Scraping dynamic (Javascript) generated webpage was not so easy as I thought. I was first planning to do a simple Python script but realised that to actually have the webpage rendered with all JS stuff, I required a web browser engine or similar. 

A *headless* browser would probably do the trick, which is what is mostly used when doing UI testing. I tried out different Python libraries and also some Golang libraries but I could not get the to work without installing additional software. I wanted to be able to run this on a small linux server on AWS without any extra stuff installed. I eventually found [PhantomJS](https://phantomjs.org/), which you can run as a small binary (on any/most OS) and it will render dynamic Javascript generated HTML website without any extra browser client installed. This was perfect. So `lager-kollen` is Go code that executes `phantomjs` binary to download a webpage, output the HTML content to stdout and read/parse it. Works smoothly.

For each webpage scraped, dedicated code needs to be written to figure out what text to extract from it. I have thus far only implemented support for [power.se](https://www.power.se/) (which seems to work for any item on their page). Future sites will require you to add support for that specific site. The parsing of HTML code is done with [goquery](https://github.com/PuerkitoBio/goquery), which is similart to JQuery and works really nice.

## What the script does

You add URLs to the script. For each URL added, a go-routine is started that continously requests that webpage content of that URL. A [scraper.go](Scraper) is used to scrape the content and figure out if the item is in stock (compared to last time the page was scraped). If the item is in stock, a notification is (meant to be) sent alerting that it is time to unleash your big wallet.

URLs can be removed, i.e. set to no longer be moitored. Adding/removing URLs are done via REST API, see below.

## PhantomJS wrapper script

`lager-kollen` will generate a Javascript wrapper script used to genrate and download webpage HTML content. It assumes you have `PhantomJS` installed already.

Below is the sccript that it generates

    >> cat /tmp/phantomjs-wrapper.js

    var system = require('system');
    url = system.args[1];
    var page = require("webpage").create();
    page.open(url, function () {
        console.log(page.content);
        phantom.exit();
    });


# Pre-reqs

You need PhantomJS binary installed. `lager-kollen` will execute it for you and will try to find it automatically by looking at the following places:

    Environment variable "PHANTOMJS_BIN"
    /usr/local/bin/phantomjs
    ./phantomjs


# Usage

Run it with

    go run *.go 

It will look for a text file called "urls.csv" where you list the URLs to parse. A bunch of `Worker` go-threads will pull URLs from that file, execute phantomJS to download the page content, parase it and figure out if it needs to notify.
A simple cache keeps track of webpages scraped.

Additional input parameters are available, run with `-help` to list those.

If you set `VERBOSE` env. variable, the script it will output verbose logging. Example

    ❯ VERBOSE=1 go run *.go
    2021/02/10 16:04:48 Reading URLs from ./urls.csv
    2021/02/10 16:04:48 Server started, listening on :8080
    2021/02/10 16:04:48 Processing URL=https://www.power.se/gaming/gamingmoss-och-tangentbord/gaming-mus/cepter-rogue-gamingmus/p-993356/
    2021/02/10 16:04:57 Scraped: {"url":"https://www.power.se/gaming/gamingmoss-och-tangentbord/gaming-mus/cepter-rogue-gamingmus/p-993356/","product":"CEPTER ROGUE GAMINGMUSCEPTER ROGUE GAMINGMUS","lastVisit":"2021-02-10 16:04:57.419487 +0100 CET m=+8.665971460","lastStatusText":"50+ i lager","hasItemInStock":true}
    2021/02/10 16:04:57 Site power.se - CEPTER ROGUE GAMINGMUSCEPTER ROGUE GAMINGMUS has the item in stock!
    2021/02/10 16:04:57 Processing URL=https://www.power.se/gaming/konsol/playstation-5/p-1077687/
    2021/02/10 16:05:02 Scraped: {"url":"https://www.power.se/gaming/konsol/playstation-5/p-1077687/","product":"PLAYSTATION 5PLAYSTATION 5","lastVisit":"2021-02-10 16:05:02.718922 +0100 CET m=+13.965400726","lastStatusText":"Inte i lager","hasItemInStock":false}

# Build

Simply execute 

    go build

# REST API

## Add URL, path=/api/urls/add?url={url}

Add a URL to be monitored by the application.

Example

    curl 'http://localhost:8080/api/urls/add?url=https://www.power.se/gaming/konsol/playstation-5/p-1077687/'

## Remove URL, path=/api/urls/remove?url=url

Remove URL from being monitored

Example

    curl 'http://localhost:8080/api/urls/remove?url=https://www.power.se/gaming/konsol/playstation-5/p-1077687/'

## Debug requests

Debug requests are used to troubleshoot and test the application

### Post HTML content from file and scrape it

A way to more easily test the scraper, is to download the webpage content to file and then POST the file content to the script.

First download the webpage content to .html file

    >> phantomjs /tmp/phantomjs-wrapper.js https://www.power.se/gaming/konsol/playstation-5/p-1077687/ > power.se-PS5.html

Then POST it to the script, you need to help it by naming the URL domain so that the correct scraper is used and returns the output what was scraped:

    >> curl -X POST -d @power.se-ps5.html http://localhost:8080/debug/scraper/power.se
    Scraped: url=http://power.se, product=PLAYSTATION 5, lastStatus=Inte i lager, inStock=false

# Future enhancements

 - Persistent cache (-> sqlite?)
 - Email notification
 - Push notification to mobile device or other apps??
