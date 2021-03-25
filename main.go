package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
)

func testSSRF(payloads, match string, appendMode bool) {
	file, err := os.Open(payloads)
	if err != nil {
		log.Fatalf("File could not be read: %v", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(os.Stdin)
	payloadScanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for payloadScanner.Scan() {
			link := scanner.Text()
			payload := payloadScanner.Text()

			u, err := url.Parse(link)
			if err != nil {
				log.Fatalf("URL format err: %v", err)
			}
			qs := url.Values{}
			for param, value := range u.Query() {
				if appendMode {
					qs.Set(param, value[0]+payload)
				} else {
					qs.Set(param, payload)
				}
			}
			u.RawQuery = qs.Encode()
			fmt.Printf("[+] Testing URL with payload: %s", u.RawQuery)
			// TO DO: Send request with payload and capture response
		}
	}
}

func main() {
	var concurrency int
	var payloads string
	var match string
	var appendMode bool
	flag.IntVar(&concurrency, "c", 20, "Set the threads to use")
	flag.StringVar(&payloads, "p", "", "Payload list")
	flag.StringVar(&match, "m", "", "Match the response with a pattern (e.g): Success")
	flag.BoolVar(&appendMode, "a", true, "Append the payload to parameter")
	flag.Parse()

	if payloads != "" {
		// Create goroutines
		var wg sync.WaitGroup
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				testSSRF(payloads, match, appendMode)
				wg.Done()
			}()
			wg.Wait()
		}
	}
}
