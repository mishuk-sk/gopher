package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mishuk-sk/gopher/ChooseYourOwnAdventure/htmlbuilder"
	"github.com/mishuk-sk/gopher/ChooseYourOwnAdventure/htmlbuilder/storyparser"
)

func main() {
	addr := flag.String("host", "localhost:4000", "Used to provide host for server")
	storyFile := flag.String("story", "gopher.json", "Used to provide story json file")
	flag.Parse()
	idleConnsClosed := make(chan struct{})
	f, err := os.Open(*storyFile)
	defer f.Close()
	if err != nil {
		log.Fatalf("Can't open story file with name %s. %s\n", *storyFile, err)
	}
	storyBytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(storyBytes))
	story, err := storyparser.MappedStory(storyBytes)
	if err != nil {
		log.Fatal(err)
	}

	router := createHandlers(story)

	server := http.Server{
		Addr:         *addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      router,
	}

	if err := startServer(&server, idleConnsClosed); err != nil {
		panic(err)
	}
	<-idleConnsClosed
}

func createHandlers(story map[string]storyparser.Arc) *http.ServeMux {
	router := http.NewServeMux()
	for path, arc := range story {
		router.HandleFunc(fmt.Sprint("/", path), createHandler(arc))
	}
	return router
}

func createHandler(arc storyparser.Arc) http.HandlerFunc {
	page := htmlbuilder.GetPage(arc)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/html")
		w.Write(page)
	}
}

func startServer(srv *http.Server, waitShutdown chan struct{}) error {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(waitShutdown)
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
