package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	hn "prashamhtrivedi/quiethn"
)

//go:embed index.gohtml
var templateFile embed.FS

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFS(templateFile, "index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	sc := storyCache{
		duration:   6 * time.Second,
		numStories: numStories,
	}

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		temp := storyCache{
			duration:   6 * time.Second,
			numStories: numStories,
		}
		temp.getStories()
		sc.mutex.Lock()
		sc.cache = temp.cache
		sc.expiration = temp.expiration
		sc.mutex.Unlock()
		<-ticker.C
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := sc.getStories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		data := templateData{
			Stories: stories,
			Time:    time.Since(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

type storyCache struct {
	numStories int
	cache      []item
	expiration time.Time
	duration   time.Duration
	mutex      sync.Mutex
}

func (sc *storyCache) getStories() ([]item, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if time.Since(sc.expiration) < 0 {
		return sc.cache, nil
	}
	fmt.Println("Cache Invalid")
	stories, err := getTopStories(sc.numStories)
	sc.expiration = time.Now().Add(sc.duration)
	if err != nil {
		return nil, err
	}
	sc.cache = stories
	return sc.cache, nil
	// cache = stories
}

func getTopStories(numStories int) ([]item, error) {

	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("failed to load top stories")
	}

	var stories []item
	at := 0
	for len(stories) < numStories {
		need := (numStories - len(stories)) * 5 / 4
		stories = append(stories, getStories(ids[at:at+need])...)
		at += need
	}

	return stories[:numStories], nil
}

func getStories(ids []int) []item {
	numStories := len(ids)
	var stories []item
	type result struct {
		index int
		item  item
		err   error
	}
	resultCh := make(chan result)
	for i := 0; i < numStories; i++ {

		go func(index, id int) {
			var client hn.Client
			hnItem, err := client.GetItem(id)
			if err != nil {
				resultCh <- result{index: index, err: err}
			} else {
				resultCh <- result{index: index, item: parseHNItem(hnItem)}
			}

		}(i, ids[i])

	}

	var results []result
	for i := 0; i < numStories; i++ {
		resultData := <-resultCh

		results = append(results, resultData)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	for _, resultData := range results {
		if resultData.err != nil {
			continue
		} else {
			if isStoryLink(resultData.item) {
				stories = append(stories, resultData.item)
			}
		}
	}
	return stories
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
