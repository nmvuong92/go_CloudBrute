package main

import (
	"fmt"
	"github.com/gookit/goutil/dump"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"time"
)

func main() {
	req, err := http.NewRequest(http.MethodGet, "https://httpbin.org/get", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*80) // quá timeout
	defer cancel()

	req = req.WithContext(ctx)

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil { // throw ra lỗi khi quá timeout
		dump.P(err)
		return
	}
	defer res.Body.Close()
	out, err := io.ReadAll(res.Body)
	if err != nil {
		dump.P(err)
		return
	}

	dump.P(string(out))
}
