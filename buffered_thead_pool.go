package main

import (
	"fmt"
	"net/http"
	"sync"
)

type URLResult struct {
	URL   string
	Code  int
	Error error
}

func main() {
	urls := []string{
		"https://example.com",
		"https://google.com",
		"https://facebook.com",
		"https://twitter.com",
		"https://instagram.com",
	}

	results := make(chan URLResult, len(urls))
	var wg sync.WaitGroup
	poolSize := 4 // Kích thước thread pool

	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urls {
				resp, err := http.Get(url)

				code := 0
				if err == nil {
					code = resp.StatusCode
					resp.Body.Close()
				}
				results <- URLResult{URL: url, Code: code, Error: err}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res.Error != nil {
			fmt.Printf("Error crawling %s: %s\n", res.URL, res.Error.Error())
		} else {
			fmt.Printf("URL: %s, Status Code: %d\n", res.URL, res.Code)
		}
	}
}
