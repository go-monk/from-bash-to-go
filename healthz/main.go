package main

import (
	"log"
	"net/http"
	"time"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("healthy"))
}

func Healthz2Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMovedPermanently)
	w.Write([]byte("healthy"))
}

func Healthz3Handler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("healthy"))
}

func main() {
	http.HandleFunc("/healthz", HealthzHandler)
	http.HandleFunc("/healthz2", Healthz2Handler)
	http.HandleFunc("/healthz3", Healthz3Handler)

	port := "8080"
	log.Printf("Starting server on port %s ...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
