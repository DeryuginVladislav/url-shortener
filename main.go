package main

import (
	"log"
	"net/http"
	"url-shortener/handlers"
)

func main() {
	http.HandleFunc("/shorten", handlers.HandleShorten)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("server is running on localhost")

}
