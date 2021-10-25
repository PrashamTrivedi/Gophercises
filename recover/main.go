package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/", sourceCodeHandler)
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", recoverMiddleWare(mux)))
}

func sourceCodeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	lineStr := r.FormValue("line")
	line, err := strconv.Atoi(lineStr)
	if err != nil {
		line = -1
	}
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var lines [][2]int
	if line > 0 {
		lines = append(lines, [2]int{line, line})
	}
	lexer := lexers.Get("go")
	iterator, err := lexer.Tokenise(nil, b.String())
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}
	formatter := html.New(html.TabWidth(2), html.LineNumbersInTable(true), html.LinkableLineNumbers(true, ""), html.HighlightLines(lines), html.WithLineNumbers(true))
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<style>pre { font-size: 1.2em; }</style>")
	formatter.Format(w, style, iterator)
}

func recoverMiddleWare(app http.Handler) http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				stackMessage := string(stack)
				log.Println(stackMessage)
				rw.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(rw, "<h1>Panic: %v</h1><br><pre>%s</pre>", err, parseStack(stackMessage))

			}
		}()

		// nrw := &respWriter{
		// 	ResponseWriter: rw,
		// }
		app.ServeHTTP(rw, r)
		// nrw.flush()
	}
}

type respWriter struct {
	http.ResponseWriter
	writes [][]byte
	status int
}

func (rw *respWriter) Write(b []byte) (int, error) {
	rw.writes = append(rw.writes, b)
	return len(b), nil
}

func (rw *respWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
}

func (rw *respWriter) flush() error {
	if rw.status != 0 {
		rw.ResponseWriter.WriteHeader(rw.status)
	}
	for _, write := range rw.writes {
		_, err := rw.ResponseWriter.Write(write)
		if err != nil {
			return err
		}
	}
	return nil
}

func panicDemo(w http.ResponseWriter, r *http.Request) {

	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}

func parseStack(stack string) string {
	lines := strings.Split(stack, "\n")

	for li, line := range lines {
		if len(line) == 0 || line[0] != '\t' {
			continue
		}
		file := ""
		for index, ch := range line {
			if ch == ':' {
				file = line[1:index]
				break
			}
		}

		lineStrBuilder := strings.Builder{}
		for i := len(file) + 2; i < len(line); i++ {
			if line[i] < '0' || line[i] > '9' {
				break
			}

			lineStrBuilder.WriteByte(line[i])
		}
		lineStr := lineStrBuilder.String()
		v := url.Values{}
		v.Set("path", file)
		v.Set("line", lineStr)
		lines[li] = fmt.Sprintf("\t<a href=\"/debug?%s\">%s:%s</a>%s", v.Encode(), file, lineStr, line[len(file)+2+len(lineStr):])
	}
	return strings.Join(lines, "\n")
}
