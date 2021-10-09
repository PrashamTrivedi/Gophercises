package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"prashamhtrivedi/cyoa"
	"strings"
)

func main() {
	port := flag.Int("port", 3030, "Port to start server")
	file := flag.String("file", "story.json", "CYOA File")
	flag.Parse()

	fmt.Println(*file)

	fileData, error := os.Open(*file)
	reportError(error)

	story, err := cyoa.JSONStory(fileData)
	if err != nil {
		reportError(err)
	}

	var defaultHandlerTemplate = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <section class="page">
      <h1>{{.Title}}</h1>
      {{range .Paragraphs}}
        <p>{{.}}</p>
      {{end}}
      {{if .Options}}
        <ul>
        {{range .Options}}
          <li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
        </ul>
      {{else}}
        <h3>The End</h3>
      {{end}}
    </section>
    <style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: #FFFCF6;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #777;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #7792a2;
      }
      p {
        text-indent: 1em;
      }
	  @media (prefers-color-scheme: dark) {
    *, #nav h1 a {
      color: #FDFDFD;
    }

    body {
      background: #121212;
    }

	.page {
		background: #121212;
	}

    pre, code {
      background-color: #262626;
    }

    #sub-header, .date {
      color: #BABABA;
    }

    hr {
      background: #EBEBEB;
    }
  }
    </style>
  </body>
</html>`
	temp := template.Must(template.New("").Parse(defaultHandlerTemplate))

	h := cyoa.NewHandler(story, cyoa.WithTemplate(temp), cyoa.WithPathFunc(pathFunc))

	mux := http.NewServeMux()
	mux.Handle("/story/", h)
	fmt.Printf("Starting the server: http://localhost:%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}
func pathFunc(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}
func reportError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
