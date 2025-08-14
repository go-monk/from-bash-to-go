package main

import (
	"fmt"
	"net/http"
	"time"
)

type HealthCheck struct {
	URL               string
	ResponseTimeout   time.Duration // defaults to zero
	HealthyStatusCode int
}

func (h HealthCheck) Do() bool {
	client := http.Client{Timeout: h.ResponseTimeout} // zero means no timeout
	resp, err := client.Get(h.URL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != h.HealthyStatusCode {
		return false
	}
	return true
}

func main() {
	healthChecks := []HealthCheck{
		{URL: "http://localhost:8080/healthz", ResponseTimeout: 2 * time.Second, HealthyStatusCode: http.StatusOK},
		{URL: "http://localhost:8080/healthz2", ResponseTimeout: 2 * time.Second, HealthyStatusCode: http.StatusMovedPermanently},
		{URL: "http://localhost:8080/healthz3", ResponseTimeout: 10 * time.Second, HealthyStatusCode: http.StatusOK},
	}
	for _, h := range healthChecks {
		if ok := h.Do(); !ok {
			fmt.Printf("%s is unhealthy\n", h.URL)
		}
	}
}
