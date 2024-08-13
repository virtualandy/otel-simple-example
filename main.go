package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/rolldice", rolldice)
	log.Println("running @ localhost:8000")

	log.Fatal(http.ListenAndServe(":8000", nil))
}
