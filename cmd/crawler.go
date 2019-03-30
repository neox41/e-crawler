package cmd

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func Crawler(sem chan struct{}, urlToGet string) {
	defer recoverCrawler()
	defer WgCrawler.Done()

	if IsInScope(urlToGet) {
		req, err := http.NewRequest("GET", urlToGet, nil)
		if err != nil {
			log.Panic(err)
		}
		req.Header.Set("User-Agent", UA)
		resp, err := HttpClient.Do(req)
		if err != nil {
			log.Println("Failed to crawl \"" + urlToGet + "\"")
			LinksLock.Lock()
			Links[urlToGet] = true
			LinksLock.Unlock()
			return
		}
		if Verbose {
			log.Println("Analysing " + urlToGet)
		}
		defer resp.Body.Close()
		bodyContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panic(err)
		}
		// Create 2 readers for parser and crawler
		reader1 := bytes.NewReader(bodyContent)
		reader2 := bytes.NewReader(bodyContent)

		body, err := ioutil.ReadAll(reader1)
		WgParser.Add(1)
		go ParseEmail(string(body))
		ParseLink(urlToGet, reader2)
		LinksLock.Lock()
		Links[urlToGet] = true
		LinksLock.Unlock()
		<-sem
	}
}
func recoverCrawler() {
	if r := recover(); r != nil {
		if Verbose {
			log.Println("Panic Error Recovered from ", r)
		}
	}
}
