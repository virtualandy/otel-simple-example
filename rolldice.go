package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	tracer  = otel.Tracer("rolldice")
	meter   = otel.Meter("rolldice")
	rollCnt metric.Int64Counter
)

func init() {
	var err error
	rollCnt, err = meter.Int64Counter("dice.rolls",
		metric.WithDescription("The number of rolls by roll value"),
		metric.WithUnit("{roll}"))
	if err != nil {
		panic(err)
	}
}

func rolldice(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "rolling dice")
	defer span.End()

	roll := 1 + rand.Intn(6)

	// Sleep for a few seconds, so we can see longer traces
	// And breakout the slow api call versus other work
	// that might be done in this function, i.e. a DB query
	time.Sleep(time.Millisecond * time.Duration(roll*100))

	span.AddEvent("simulating a call to slow api")
	err := simulateSlowAPI(500*roll, ctx)
	if err != nil {
		log.Printf("Failed to make call to simulate slow api")
	}

	rollValueAttr := attribute.Int("roll.value", roll)
	span.SetAttributes(rollValueAttr)
	rollCnt.Add(ctx, 1, metric.WithAttributes(rollValueAttr))

	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}

// A function for simulating slowness with an external http call
func simulateSlowAPI(sleepInMilliseconds int, ctx context.Context) error {
	ctx, childSpan := tracer.Start(ctx, "simulatedSlowApiSpan")
	defer childSpan.End()
	// Construct the URL with the sleep parameter
	baseURL := "https://fakeresponder.com/"
	sleepParam := "?sleep=" + strconv.Itoa(sleepInMilliseconds)

	// Create the full URL
	fullURL := baseURL + sleepParam

	// Make the GET request
	childSpan.AddEvent("GET fakeresponder.com api")
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
