package main

import (
	"fmt"
	"net/http"
	"pht/urlshortnerhandler"
)

func main() {
	fmt.Println("Hello")

	mux := defaultMux()
	pathToUrls := map[string]string{
		"/reddit": "https://www.reddit.com/r/UmbrellaAcademy/",
		"/imdb":   "https://www.imdb.com/title/tt1312171/",
	}

	mapHandler := urlshortnerhandler.MapHandler(pathToUrls, mux)

	yaml := `
- path: /saa-notes
  url: https://notes.prashamhtrivedi.in/saa/index.html
- path: /solution-arch
  url: https://prashamhtrivedi.in/passing_aws_saa.html`
	yamlHandler, err := urlshortnerhandler.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
