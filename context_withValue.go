package main

import (
	"fmt"
	"golang.org/x/net/context"
)

type key int

const (
	userKey key = iota
)

func greet(ctx context.Context) {
	// truy cap gia tri tu context
	user, ok := ctx.Value(userKey).(string)
	if !ok {
		fmt.Println("khong the tim thay thong tin nguoi dung trong context")
	}

	fmt.Printf("Xin chao, %s\n", user)
}

func main() {
	// taoj mot context cha
	ctx := context.Background()

	// tao mot context con voi thong tin nguoi dung
	ctx = context.WithValue(ctx, userKey, "Alice")

	// goi ham greet voi context da co gia tri thogn tin nguoi dung
	greet(ctx)
}
