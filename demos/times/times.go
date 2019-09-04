package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)

	go func() {
		// time.Sleep(1 * time.Second)
		time.Sleep(3 * time.Second)
		<-c
	}()

	select {
	case c <- 1:
		fmt.Println("channel...")
	case <-time.After(2 * time.Second):
		close(c)
		fmt.Println("timeout...")
	}
}
