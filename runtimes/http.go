package main

import "net/http"

func main() {
	http.HandleFunc("/", API)
	http.ListenAndServe(":8000", nil)
}
