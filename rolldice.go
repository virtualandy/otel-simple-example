package main

import (
	"fmt"
	"io"
	"io/ioutil"
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
	time.Sleep(time.Millisecond * time.Duration(roll*100))

	ctx, childSpan := tracer.Start(ctx, "child")
	// defer childSpan.End()
	span.AddEvent("calling external api")
	callAPI(500 * roll)
	childSpan.End()

	rollValueAttr := attribute.Int("roll.value", roll)
	span.SetAttributes(rollValueAttr)
	rollCnt.Add(ctx, 1, metric.WithAttributes(rollValueAttr))

	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}

// A function for simulating slowness with an external http call
func callAPI(sleepInMilliseconds int) {
	// Construct the URL with the sleep parameter
	baseURL := "https://fakeresponder.com/"
	sleepParam := "?sleep=" + strconv.Itoa(sleepInMilliseconds)

	// Create the full URL
	fullURL := baseURL + sleepParam

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Println("Error making the GET request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}
