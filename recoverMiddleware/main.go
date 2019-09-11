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
	log.Fatal(http.ListenAndServe(":3000", middleware(mux, true)))
}

//FIXME doesn't work with /panic-after/ (Status:  200)
func middleware(handler http.Handler, deb bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Error: %v\n", r)
				stack := debug.Stack()
				log.Printf("Stack trace:\n%s\n", stack)
				w.WriteHeader(http.StatusInternalServerError)
				if !deb {
					w.Write([]byte("Something went wrong"))
					return
				}
				fmt.Fprintf(w, "Error: %v;\nStack trace:%s\n", r, stack)
			}
		}()

		dw := deferedWriter{w, 0, [][]byte{}}
		handler.ServeHTTP(&dw, r)
		flush(dw)
	}
}

type deferedWriter struct {
	http.ResponseWriter
	header  int
	message [][]byte
}

func (dw *deferedWriter) Write(ms []byte) (int, error) {
	dw.message = append(dw.message, ms)
	return len(dw.message), nil
}

func (dw *deferedWriter) WriteHeader(status int) {
	dw.header = status
}

func flush(dw deferedWriter) {
	if dw.header != 0 {
		dw.WriteHeader(dw.header)
	}
	for _, b := range dw.message {
		_, err := dw.ResponseWriter.Write(b)
		if err != nil {
			log.Fatalf("Can't write to original ResponseWriter. Err: %v\n", err)
		}
	}
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
