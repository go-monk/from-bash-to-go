Source: https://github.com/go-monk/from-bash-to-go

# From Bash to Go

> Bash is great until it isn't.

If you're in DevOps, SRE, or cybersecurity, you’ve probably written countless Bash scripts to automate or glue things together. Bash is a good tool for these tasks, provided the programs are small and simple. However, as they grow more complex, they become harder to understand and modify. Additionally, the dependency on external tools (like `curl`, `awk`, `jq`) makes them difficult to deploy across diverse systems. Well-written programs in Go alleviate these Bash shortcomings significantly and bring new advantages, including a cultural agenda of radical simplicity that brings more joy :-).

## 0) Quick Health Check Script 

Consider this simple health check script:

```sh
#!/bin/bash

URL="http://localhost:8080/healthz"

STATUS=$(curl -s -o /dev/null -w "%{http_code}" -m 2 $URL)
if [ "$STATUS" -ne 200 ]; then
	echo "Service unhealthy!"
	exit 1
fi
```

Here's the Go equivalent:

```go
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
```

To test the programs above, you can run a dummy web server with `/healthz` endpoints:

```sh
go run ./healthz/main.go
```

Now you can run the scripts:

```sh
./0/healthcheck.sh 

go run ./0/healthcheck.go
# or
go build ./0/healthcheck.go && ./healthcheck
```

(Remember the Unix philosophy: no news is good news :-). You can check the exit status of the commands with `echo $?`—zero means all good.)

At first glance, there isn't much difference except for the syntax. However, the Go code has no external dependencies and can be compiled to run on any operating system and CPU architecture. For example, if you are developing on a Mac but want to deploy to a Linux server:

```sh
GOOS=linux GOARCH=arm64 go build ./0/healthcheck.go
scp healthcheck user@linuxbox.com:
```

Small scripts often grow over time.

## 1) Check Multiple Services 

A colleague or your boss likes the script and asks you (or someone else) to add functionality to health check more than one service. Sure, no problem, you think. But when extending the script, you discover that one of the services replies with a 301 status instead of 200 (well, things tend to get messy).

Let's continue with Go since the task is becoming more complex.

First, define a custom type (a struct with two fields: a string and an integer) to hold data about the health check endpoints:

```go
type HealthCheck struct {
	URL               string
	HealthyStatusCode int
}
```

Next, create a function attached to this custom type—via the `(h HealthCheck)` part—that performs the check:

```go
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
```

We use the standard library [http](https://pkg.go.dev/net/http) package instead of `curl`. You can use `go doc http.Get` to see details about the `Get` method from the package. Our `Do` method returns a boolean indicating whether the service is healthy (`true`) or not (`false`).

Finally, define the services to health check as a slice of `HealthCheck` structs. Then loop over them and call the `Do` method on each:

```go
	healthChecks := []HealthCheck{
		{URL: "http://localhost:8080/healthz", HealthyStatusCode: http.StatusOK},
		{URL: "http://localhost:8080/healthz2", HealthyStatusCode: http.StatusMovedPermanently},
	}
	for _, h := range healthChecks {
		if ok := h.Do(); !ok {
			fmt.Printf("%s is unhealthy\n", h.URL)
		}
	}
```

Check that the program compiles and runs:

```sh
go run ./1/main.go
```

Nice! Time for a coffee break—you deserve it.

## 2) Different Timeouts

You return to your desk with a coffee and see a Slack message like "please add the `healthz3` endpoint to your script". Sure, easy enough—you add `{URL: "http://localhost:8080/healthz3", HealthyStatusCode: http.StatusOK},` and run the script:

```sh
❯ go run ./1/main.go 
http://localhost:8080/healthz3 is unhealthy
```

Hmm. After some investigation, you discover that the endpoint takes 3 seconds to reply. You inform the requester, and they reply, "yeah, i know, that's fine". Ok then. Luckily, you just need to make a couple of easy changes to accommodate this slow service. Run `diff ./1/main.go ./2/main.go` to see the changes.

## 3) Read Configuration from a JSON File

At this point, it's clear that the script is outgrowing the original "quick and dirty" approach. It would be better to read the health check endpoints configuration from a file. A JSON file is a simple choice to start with.

Create a function to read a file and return a slice of health checks:

```go
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
```

Replace the hardcoded health checks in `main()` like this:

```go
	healthChecks, err := readConfig("healthchecks.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck: %v\n", err)
		os.Exit(1)
	}
```

As an exercise, remove the hardcoded filename (`healthchecks.json`) and get the filename from the command-line arguments instead.
