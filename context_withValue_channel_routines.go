package main

import (
	"fmt"
	"golang.org/x/net/context"
	"time"
)

type key int

const (
	messageKey key = iota
)

func printMessage(ctx context.Context, ch chan<- string) { // channel chi gui message
	message, ok := ctx.Value(messageKey).(string)
	if !ok {
		ch <- "Khong tim thay thong diep trong context"
	}

	time.Sleep(1 * time.Second)
	ch <- message
}

func main() {
	// tao mot context cha
	ctx := context.Background()

	// tao mot context con voi thong diep
	ctx = context.WithValue(ctx, messageKey, "Xin chao!")

	// tao mot kenh cho ket qua
	ch := make(chan string)

	// thuc thi goroutine de in thong diep
	go printMessage(ctx, ch)

	// doc ket qua tu kenh
	select {
	case result := <-ch:
		fmt.Println(result)
	case <-time.After(10 * time.Nanosecond):
		fmt.Println("Qúa thời gian chờ")
	}

}
