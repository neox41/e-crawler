package cmd

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"sort"
	"strings"
)

//https://github.com/indraniel/go-learn/blob/master/09-sort-map-keys-by-values.go
func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func IsInScope(urlToGet string) bool {
	host := getHost(urlToGet)
	for _, d := range Domains {
		if d == host {
			return true
		}
	}
	return false
}
func AddToScope(urlToAdd string) {
	host := getHost(urlToAdd)
	Domains = append(Domains, host)
}
func GetRootPath(urlRaw string) string {
	url, err := url.Parse(urlRaw)
	if err != nil {
		log.Panic("Invalid URL")
	}
	return fmt.Sprintf("%s://%s", url.Scheme, getHost(urlRaw))
}
func getHost(urlRaw string) string {
	defer recoverHost()
	var host string
	url, err := url.Parse(urlRaw)
	if err != nil {
		log.Panic("Invalid URL")
	}
	if strings.Contains(url.Host, ":") {
		host, _, _ = net.SplitHostPort(url.Host)
	} else {
		host = url.Host
	}
	return host
}
func UniqueEmails() (uniqueEmailsSorted, uniqueObEmailsSorted PairList) {
	uniqueEmails := make(map[string]int)
	for _, email := range Emails {
		if _, value := uniqueEmails[email]; !value {
			uniqueEmails[email] = 1
		} else {
			uniqueEmails[email]++
		}
	}
	uniqueEmailsSorted = make(PairList, len(uniqueEmails))
	i := 0
	for k, v := range uniqueEmails {
		uniqueEmailsSorted[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(uniqueEmailsSorted))

	uniqueObEmails := make(map[string]int)
	for _, email := range ObEmails {
		if _, value := uniqueObEmails[email]; !value {
			uniqueObEmails[email] = 1
		} else {
			uniqueObEmails[email]++
		}
	}
	uniqueObEmailsSorted = make(PairList, len(uniqueObEmails))
	i = 0
	for k, v := range uniqueObEmails {
		uniqueObEmailsSorted[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(uniqueObEmailsSorted))
	return
}
func recoverHost() {
	if r := recover(); r != nil {
		if Verbose {
			log.Println("Panic Error Recovered from ", r)
		}
	}
}
