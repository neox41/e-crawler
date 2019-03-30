package cmd

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func parseObEmail(page string) {
	// mr [at] mattiareggiani.com
	re := regexp.MustCompile(`[a-z0-9-]{1,30} \[at\] [a-z0-9-]{1,65}.[a-z]{1,}`)
	emails := re.FindAllString(page, -1)
	for e := range emails {
		ObEmails = append(ObEmails, strings.Replace(emails[e], " [at] ", "@", 1))
	}
	// mr[at]mattiareggiani.com
	re = regexp.MustCompile(`[a-z0-9-]{1,30}\[at\][a-z0-9-]{1,65}.[a-z]{1,}`)
	emails = re.FindAllString(page, -1)
	for e := range emails {
		ObEmails = append(ObEmails, strings.Replace(emails[e], "[at]", "@", 1))
	}
	// mr [at] mattiareggiani [dot] com
	re = regexp.MustCompile(`[a-z0-9-]{1,30} \[at\] [a-z0-9-]{1,65} \[dot\] [a-z]{1,}`)
	emails = re.FindAllString(page, -1)
	for e := range emails {
		ObEmails = append(ObEmails, strings.Replace(strings.Replace(emails[e], " [at] ", "@", 1), " [dot] ", ".", 1))
	}
	// mr[at]mattiareggiani[dot]com
	re = regexp.MustCompile(`[a-z0-9-]{1,30}\[at\][a-z0-9-]{1,65}\[dot\][a-z]{1,}`)
	emails = re.FindAllString(page, -1)
	for e := range emails {
		ObEmails = append(ObEmails, strings.Replace(strings.Replace(emails[e], "[at]", "@", 1), "[dot]", ".", 1))
	}
	// TODO add other
}
func ParseEmail(page string) {
	defer WgParser.Done()
	re := regexp.MustCompile(`[a-z0-9-._%+\-]+@[a-z0-9-.\-]+\.[a-z]{1,}`)
	emails := re.FindAllString(page, -1)
	for e := range emails {
		Emails = append(Emails, emails[e])
	}
	parseObEmail(page)
}
func ParseLink(currentUrl string, page io.Reader) {
	z := html.NewTokenizer(page)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data == "a" || t.Data == "link" || t.Data == "script" || t.Data == "iframe" {
				ok, url := getValue(t)
				if !ok || !checkForException(url) {
					continue
				}
				if !strings.HasPrefix(url, "http") {
					urlA := strings.Split(currentUrl, "/")
					urlA = urlA[:len(urlA)-1]
					baseURL := strings.Join(urlA, "/")
					if strings.HasPrefix(url, "/") {
						baseURL = GetRootPath(currentUrl)
					} else {
						baseURL += "/"
					}
					url = fmt.Sprintf("%s%s", baseURL, url)
				}
				if IsInScope(url) {
					LinksLock.Lock()
					if _, exist := Links[url]; !exist {
						Links[url] = false
					}
					LinksLock.Unlock()
				}
			}
		}
	}
}
func getValue(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" || a.Key == "src" {
			href = strings.Split(a.Val, "#")[0]
			ok = true
		}
	}
	return
}
func checkForException(url string) bool {
	if strings.HasPrefix(url, "javascript:") || strings.HasSuffix(url, ".pdf") || strings.HasSuffix(url, ".png") || strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".ico") || strings.HasPrefix(url, "callto://") || strings.HasPrefix(url, "mailto://") || strings.HasPrefix(url, "callto:") || strings.HasPrefix(url, "mailto:") || url == "#" || url == "" {
		return false
	} else {
		return true
	}
}
