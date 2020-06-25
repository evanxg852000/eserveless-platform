package main

import (
	"fmt"
	"net/http"
	"time"
)

// API is a http function handler
func API(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

// Ticker is a cron function handler
func Ticker(envs interface{}) {
	fmt.Printf("Ticker ticked at %s %v \n", time.Now().String(), envs)
}
