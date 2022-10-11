# e-crawler

Web crawler to extract e-mails from website

## Usage
```
go get github.com/mattiareggiani/e-crawler/cmd@latest

./e-crawler -h

-domain value
        Domain(s) in scope for crawling
  -insecure
        Skip TLS certificate verification
  -output string
        Output file
  -proxy string
        Web Proxy (e.g. http://127.0.0.1:8080)
  -target string
        Target website
  -threats int
        Number of threats (default 10)
  -user-agent string
        User-Agent (default "https://github.com/mattiareggiani/e-crawler")
  -verbose
        Display more information
```
## Example
```
./e-crawler -target https://evilcorp.local -proxy http://127.0.0.1:8080 -insecure true
2019/03/30 15:27:44 Spidering the target
2019/03/30 15:27:45 Spidering the web pages
2019/03/30 15:27:56 Done in 11.5839335s
2019/03/30 15:27:56 Emails Found:
elliot@evilcorp.local (2)
support@evilcorp.local (2)
angela@evilcorp.local (1)
phillip@evilcorp.local (1)
```
