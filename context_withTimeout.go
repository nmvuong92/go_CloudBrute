package main

import (
	"fmt"
	"golang.org/x/net/context"
	"time"
)

func doSomething(ctx context.Context) {
	if ctx.Err() != nil {
		fmt.Printf("Context da bi huy bo")
		return
	}

	// Semulate (tinh toan) mot tac vu lau dai
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Tac vu da hoan thanh")
	case <-ctx.Done():
		fmt.Println("Tac vu da bi huy bo")
	}

}

func main() {
	// tao mot context cha
	ctx := context.Background()

	// tao mot context voi thoi gian toi da 3s
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// thuc thi mot tac vu voi context
	doSomething(ctx)

	// kiem tra xem tac vu co hoan thanh hay khong
	if ctx.Err() != nil {
		fmt.Println("Tac vu da bi huy bo do vuot qua thoi gian toi da")
	}
}
