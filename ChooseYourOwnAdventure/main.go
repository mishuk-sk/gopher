package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	addr := flag.String("host", "localhost", "Used to provide host for server")
	storyFile := flag.String("story", "gopher.json", "Used to provide story json file")
	flag.Parse()
	idleConnsClosed := make(chan struct{})
	server := http.Server{
		Addr: *addr,
	}
	startServer(server, idleConnsClosed)
	<-idleConnsClosed
}

func startServer(srv http.Server, waitShutdown chan struct{}) error {
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
