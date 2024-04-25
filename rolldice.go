package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func rolldice(w http.ResponseWriter, r *http.Request) {
	roll := 1 + rand.Intn(6)

	// Sleep for a few seconds, so we can see longer traces
	// And breakout the slow api call versus other work
	// that might be done in this function, i.e. a DB query
	time.Sleep(time.Millisecond * time.Duration(roll*100))

	log.Println("simulating a call to slow api")
	err := simulateSlowAPI(500 * roll)
	if err != nil {
		log.Printf("Failed to make call to simulate slow api")
	}

	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}

// A function for simulating slowness with an external http call
func simulateSlowAPI(sleepInMilliseconds int) error {
	// Construct the URL with the sleep parameter
	baseURL := "https://fakeresponder.com/"
	sleepParam := "?sleep=" + strconv.Itoa(sleepInMilliseconds)

	// Create the full URL
	fullURL := baseURL + sleepParam

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		log.Println("Error making the GET request:", err)
		return err
	}
	defer resp.Body.Close()

	// We can ignore the response body
	// But if you want to see it, uncomment this
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println("Error reading the response body:", err)
	// 	return err
	// }

	return nil
}
