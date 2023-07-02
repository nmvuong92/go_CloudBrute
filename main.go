package main

// an early version of cloud brute plugin for HunterSuite.io
import (
	engine "github.com/0xsha/cloudbrute/internal"
)

func main() {
	var threads = 80
	urls := []string{
		"yyy.s3.amazonaws.com",
		"xxx.s3.amazonaws.com",
	}
	var details engine.RequestDetails
	engine.AsyncHTTPHead(urls, threads, 10, details, "out.txt")
}
