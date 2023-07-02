package main

import (
	"sync"
	"time"
)

func main() {
	c := make(chan string) //The channel is created with an empty list of receivers and senders.

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		c <- `foo` ////Our first goroutine sends the value foo to the channel, line 16.
	}()

	//The channel acquires a struct sudog from a pool that will represent the sender.
	//This structure will keep reference to the goroutine and the value foo.
	//This sender is now enqueued in the sendq attribute.
	//The goroutine moves into a waiting state with the reason “chan send”.
	go func() {
		defer wg.Done()

		time.Sleep(time.Second * 1)
		println(`Message: ` + <-c) //Our second goroutine will read a message from the channel, line 23.
	}()

	wg.Wait()
}
