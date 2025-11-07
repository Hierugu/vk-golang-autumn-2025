package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	staticHTMLHandler := http.FileServer(http.Dir("web/html"))
	staticHandler := http.FileServer(http.Dir("web"))

	staticMux := http.NewServeMux()
	staticMux.Handle("/", staticHTMLHandler)
	staticMux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	apiMux := http.NewServeMux()

	siteMux := http.NewServeMux()
	siteMux.Handle("/", staticMux)
	siteMux.Handle("/", apiMux)

	log.Println("Starting server on :8080")
	server := &http.Server{
		Addr: ":8080",
		//Handler:      stripTrailingSlash(mux),
		Handler:      siteMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
