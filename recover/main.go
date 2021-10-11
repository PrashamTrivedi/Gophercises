package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", recoverMiddleWare(mux, true)))
}

func recoverMiddleWare(app http.Handler, dev bool) http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				stackMessage := string(stack)
				log.Println(stackMessage)
				if !dev {

					http.Error(rw, "Something went wrong", http.StatusInternalServerError)
				} else {
					rw.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(rw, "<h1>Panic: %v</h1><br><pre>%s</pre>", err, stackMessage)

				}
			}
		}()

		nrw := &respWriter{
			ResponseWriter: rw,
		}
		app.ServeHTTP(nrw, r)
		nrw.flush()
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
