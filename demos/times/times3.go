package main

import (
	"fmt"
	"time"
)


func main() {
	timer := time.NewTimer(3 * time.Second)

	go func() {
		<-timer.C
		fmt.Println("Timer has expired.")
	}()

	timer.Reset(0)
	time.Sleep(10 * time.Second)
}