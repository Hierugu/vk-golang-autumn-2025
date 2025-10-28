package main

import (
	"log"
	"net/http"
)

func main() {
	// register the handler from server.go (it must be defined as: func Handler(http.ResponseWriter, *http.Request))
	http.HandleFunc("/", SearchServer)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
