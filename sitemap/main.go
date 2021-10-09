package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"prashamhtrivedi/link"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

func main() {
	urlFlag := flag.String("url", "http://prashamhtrivedi.in/", "URL to get Sitemap for")
	maxDepth := flag.Int("depth", 4, "Maximum number of depth for sitemap")

	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)

	toXml := urlSet{
		XmlNs: xmlns,
	}
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{page})
	}

	enc := xml.NewEncoder(os.Stdout)
	fmt.Println(xml.Header)
	enc.Indent("", " ")
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}

	fmt.Println()

	// fmt.Println(strings.Join(pages, "\n"))

}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nextQ := map[string]struct{}{
		urlStr: struct{}{},
	}

	for i := 0; i <= maxDepth; i++ {
		q, nextQ = nextQ, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for url, _ := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range getUrl(url) {
				if _, ok := seen[url]; !ok {
					nextQ[link] = struct{}{}
				}
			}
		}
	}
	var ret []string

	for url, _ := range seen {
		ret = append(ret, url)
	}
	return ret
}
func getUrl(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	requestUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: requestUrl.Scheme,
		Host:   requestUrl.Host,
	}

	pages := hrefs(resp.Body, *baseUrl)

	return pages
}

type loc struct {
	Value string `xml:"loc"`
}

type urlSet struct {
	Urls  []loc  `xml:"url"`
	XmlNs string `xml:"xmlns,attr"`
}

func hrefs(r io.Reader, baseUrl url.URL) []string {
	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}
	// fmt.Println(links)

	var hrefs []string

	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "../"):
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, fmt.Sprintf("%s%s", baseUrl.String(), l.Href))
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}

	}
	hrefs = filter(hrefs, func(s string) bool {
		return strings.Contains(s, baseUrl.Host)
	})

	return hrefs
}

func filter(links []string, keepFunc func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFunc(link) {
			ret = append(ret, link)
		}
	}

	return ret
}
