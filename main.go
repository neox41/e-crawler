package main

import (
	"crypto/tls"
	"github.com/mattiareggiani/e-crawler/cmd"
	"flag"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	parseArgs()
	bindVars()
}

func main() {
	log.Println("Spidering the target")
	sem := make(chan struct{}, cmd.Threats)
	start := time.Now()
	cmd.WgCrawler.Add(1)
	sem <- struct{}{}
	go cmd.Crawler(sem, cmd.Target)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Interrupted")
		cmd.Report()
		os.Exit(1)
	}()
	cmd.WgCrawler.Wait()
	if len(cmd.Links) > 0 {
		log.Println("Spidering the web pages")
		for {
			pendingLinks := false
			for link := range cmd.Links {
				if !cmd.Links[link] {
					sem <- struct{}{}
					cmd.WgCrawler.Add(1)
					go cmd.Crawler(sem, link)
					pendingLinks = true
				}
			}
			if !pendingLinks {
				break
			}
			cmd.WgCrawler.Wait()
			cmd.WgParser.Wait()
		}
	}
	log.Printf("Done in %v\n", time.Since(start))
	cmd.Report()
}
func parseArgs() {
	flag.StringVar(&cmd.Target, "target", "", "Target website")
	flag.Var(&cmd.Domains, "domain", "Domain(s) in scope for crawling")
	flag.StringVar(&cmd.Proxy, "proxy", "", "Web Proxy (e.g. http://127.0.0.1:8080)")
	flag.StringVar(&cmd.UA, "user-agent", "e-crawler (https://github.com/mattiareggiani/e-crawler)", "User-Agent")
	flag.BoolVar(&cmd.Insecure, "insecure", false, "Skip TLS certificate verification")
	flag.StringVar(&cmd.OutputFile, "output", "", "Output file")
	flag.IntVar(&cmd.Threats, "threats", 10, "Number of threats")
	flag.BoolVar(&cmd.Verbose, "verbose", false, "Display more information")
	flag.Parse()
}

func bindVars() {
	// Target
	if !(len(cmd.Target) > 0) {
		log.Fatalln("Specify the target")
	}
	u, err := url.Parse(cmd.Target)
	if err != nil {
		log.Fatalf("Target error: %v", err)
	}
	if u.Path == "" {
		cmd.Target += "/"
	}

	// TLS Certificate
	if cmd.Insecure {
		cmd.HttpTransport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	} else {
		cmd.HttpTransport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: false}}
	}

	// Proxy
	if len(cmd.Proxy) > 0 {
		proxyUrl, err := url.Parse(cmd.Proxy)
		if err != nil {
			log.Fatalf("Proxy error: %v", err)
		}
		cmd.HttpTransport.Proxy = http.ProxyURL(proxyUrl)
	}

	// HTTP Timeout
	cmd.HttpTransport.TLSHandshakeTimeout = 5 * time.Second
	cmd.HttpTransport.Dial = (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial
	cmd.HttpClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: cmd.HttpTransport,
	}

	// Add the current target domain to the scope
	cmd.AddToScope(cmd.Target)

	// init links map
	cmd.Links = make(map[string]bool)
}
