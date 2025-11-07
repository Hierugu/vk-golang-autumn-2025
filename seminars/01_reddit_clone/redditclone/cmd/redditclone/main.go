package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	staticHTMLHandler := http.FileServer(http.Dir("static/html"))
	staticHandler := http.FileServer(http.Dir("static"))

	mux := http.NewServeMux()
	mux.Handle("/", staticHTMLHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	log.Println("Starting server on :8080")
	server := &http.Server{
		Addr: ":8080",
		//Handler:      stripTrailingSlash(mux),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
