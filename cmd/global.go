package cmd

import (
	"net/http"
	"sync"
)

type domainFlags []string
type Pair struct {
	Key   string
	Value int
}
type PairList []Pair

func (i *domainFlags) String() string {
	return "Domains in scope"
}

func (i *domainFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// Global
var (
	Links         map[string]bool
	Emails        []string
	ObEmails      []string
	WgCrawler     sync.WaitGroup
	WgParser      sync.WaitGroup
	LinksLock     sync.RWMutex
	HttpClient    *http.Client
	HttpTransport *http.Transport
)

// Args
var (
	Domains    domainFlags
	Target     string
	Proxy      string
	Insecure   bool
	UA         string
	OutputFile string
	Threats    int
	Verbose    bool
)
