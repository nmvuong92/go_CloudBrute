package mainx

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

	results := make(chan URLResult)
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			resp, err := http.Get(u)

			code := 0
			if err == nil {
				code = resp.StatusCode
				resp.Body.Close()
			}
			results <- URLResult{URL: u, Code: code, Error: err}
		}(url)
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
