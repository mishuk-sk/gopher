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
func middleware(handler http.Handler, deb bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(rw http.ResponseWriter) {
			if r := recover(); r != nil {
				log.Printf("Error: %v\n", r)
				stack := debug.Stack()
				log.Printf("Stack trace:\n%s\n", stack)
				rw.WriteHeader(http.StatusInternalServerError)
				if !deb {
					rw.Write([]byte("Something went wrong"))
					return
				}
				fmt.Fprintf(rw, "Error: %v;\nStack trace:%s\n", r, stack)
			}
		}(w)
		handler.ServeHTTP(w, r)
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
