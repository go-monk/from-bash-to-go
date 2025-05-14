# From Bash to Go

> Bash is great until it isn't.

If you're in DevOps, SRE or cybersecurity, you’ve probably written a thousand Bash scripts to automate or glue things together. And Bash is a good tool for these tasks providing the programs are small and simple. But the moment they get just a bit more complex, things start falling apart. Slowly but surely they'll become harder and harder to understand and modify. Also the dependency on external tools (like `curl`, `awk`, `jq`) makes them difficult to deploy to diverse systems. Well written programs in Go alleviate all of the mentioned Bash shortcomings significantly and bring new advantages of which a cultural agenda of radical simplicity is not the least. I began using Go (after experimenting with Python) for tools, automation, security and integration work in 2018. I have never looked back. 

## 0) Quick healthcheck script 

Look at this simple healthcheck script:

```sh
#!/bin/bash

URL="http://localhost:8080/healthz"

STATUS=$(curl -s -o /dev/null -w "%{http_code}" -m 2 $URL)
if [ "$STATUS" -ne 200 ]; then
  echo "Service unhealthy!"
  exit 1
fi
```

Here's a Go equivalent:

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

If you want to test the programs above, you can run a dummy web server with /healthz endpoints:

```sh
go run ./healthz/main.go
```

No you can run the scripts:

```sh
./0/healthcheck.sh 

go run ./0/healthcheck.go
# or
go build ./0/healthcheck.go && ./healthcheck
```

(Remember the Unix philosophy, no news is good news :-). Also you can check the exit status of the commands with `echo $?` - zero means all good.)

Not much difference at first sight, except for the syntax of course. However, the Go code has no external dependencies and can be compiled to run on any operating system and CPU architecture. For example, if you are developing on a Mac but you want to deploy to a Linux server:

```sh
GOOS=linux GOARCH=arm64 go build ./0/healthcheck.go
scp healthcheck user@linuxbox.com:
```

Now, small scripts do get often extended.

## 1) Check multiple services 

A colleague of yours or your boss come to like the script and they want you (or somobody else) to add some functionality. It should healthcheck more than one service. Sure, no problem you think. But when extending the script you find out that one of the services replies with 301 status instead of 200 (yeah, things tend to get messy).

Ok, let's continue with Go since it's becoming more complex.

First, we define a custom type (a struct with two fields of type string and intenger) that will hold data about the healthcheck endpoints:

```go
type HealthCheck struct {
	URL               string
	HealthyStatusCode int
}
```

Next, create a function attached to this custom type - via the `(h HealthCheck)` part - that will do the check:

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

We use the standard library [http](https://pkg.go.dev/net/http) package in place of `curl`. You can use `go doc http.Get` to see details about the Get method. The function (or method) returns a bool that says wether the service is healthy (true) or not (false).

Finally, let's define the services we want to healthcheck as a slice of `HealtCheck`s. Then we loop over them and call the `Do` method on each:

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

We check the program compiles and runs:

```sh
go run ./1/main.go
```

Ok, nice. Let's grab a cup of coffee, we deserve it ...

## 2) Different timeout

You come back to your desk with a coffee and read a Slack message: "please add also healthz3 endpoint to your script ...". Sure, easy enough - we add `{URL: "http://localhost:8080/healthz3", HealthyStatusCode: http.StatusOK},` and run the script:

```sh
❯ go run ./1/main.go 
http://localhost:8080/healthz3 is unhealthy
```

Hmmm. After some investigation you find out that the endpoint takes 3 seconds to reply. You let the guy know over the Slack and he says: "yeah I know, that's fine". Luckily we just need to make couple of easy changes to accomodate for this slow service. Run `diff ./1/main.go ./2/main.go` to see those changes.

## 3) Read configuration from JSON file

As you can see that this thing is outgrowing the original "quick and dirty script" approach. You think it would be nicer to read the healtcheck endpoints configuration from a file. Probably the easiest way is to work with a JSON file.

We create a function for reading a file and returning a slice (list) of healthchecks:

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

And we replace the harcoded healthchecks in `main()` like this:

```go
	healthChecks, err := readConfig("healthchecks.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck: %v\n", err)
		os.Exit(1)
	}
```

As an exercise remove the hardcoded filename (`healthchecks.json`) and get the filename from the command arguments instead.