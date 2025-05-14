package main

import (
	"fmt"
	"net/http"
	"time"
)

type HealthCheck struct {
	URL               string
	HealthyStatusCode int
}

func (h HealthCheck) Do() bool {
	client := http.Client{Timeout: 2 * time.Second}
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
		{URL: "http://localhost:8080/healthz", HealthyStatusCode: http.StatusOK},
		{URL: "http://localhost:8080/healthz2", HealthyStatusCode: http.StatusMovedPermanently},
		// {URL: "http://localhost:8080/healthz3", HealthyStatusCode: http.StatusOK},
	}
	for _, h := range healthChecks {
		if ok := h.Do(); !ok {
			fmt.Printf("%s is unhealthy\n", h.URL)
		}
	}
}
