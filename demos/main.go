package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mailbox uint8
	var lock sync.RWMutex
	sendCond := sync.NewCond(&lock)
	recvCond := sync.NewCond(lock.RLocker())



	go func() {
		lock.Lock()
		for mailbox == 1 {
			sendCond.Wait()
		}
		mailbox = 1
		fmt.Println("An email is sent")
		lock.Unlock()
		recvCond.Signal()
	}()

	go func() {
		time.Sleep(time.Second * 2)
		lock.RLock()
		for mailbox == 0 {
			recvCond.Wait()
		}
		mailbox = 0
		fmt.Println("An email is received")
		lock.RUnlock()
		sendCond.Signal()
	}()

	time.Sleep(time.Second * 5)
}