package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type HealthCheck struct {
	URL               string
	ResponseTimeout   time.Duration // defaults to zero
	HealthyStatusCode int
}

func (h HealthCheck) Do() (bool, error) {
	client := http.Client{Timeout: h.ResponseTimeout} // zero means no timeout
	resp, err := client.Get(h.URL)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != h.HealthyStatusCode {
		return false, err
	}
	return true, nil
}

func readConfig(filepath string) ([]HealthCheck, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var hs []HealthCheck
	if err := json.Unmarshal(data, &hs); err != nil {
		return nil, err
	}
	return hs, nil
}

func main() {
	healthChecks, err := readConfig("healthchecks.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "x: %v\n", err)
		os.Exit(1)
	}
	for _, h := range healthChecks {
		ok, err := h.Do()
		if !ok {
			fmt.Printf("%s is unhealthy (%v)\n", h.URL, err)
		}
	}
}
