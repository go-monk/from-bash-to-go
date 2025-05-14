package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	URL := "http://localhost:8080/healthz"

	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(URL)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Service unhealthy!")
		os.Exit(1)
	}
}
