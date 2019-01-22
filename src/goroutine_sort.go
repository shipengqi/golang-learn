package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	var count uint32 = 0
	trigger := func(i uint32, fn func()) { // func()代表的是既无参数声明也无结果声明的函数类型
		for {
			if n := atomic.LoadUint32(&count); n == i {
				fn()
				atomic.AddUint32(&count, 1)
				break
			}
			time.Sleep(time.Nanosecond)
		}
	}
	for i := uint32(0); i < 10; i++ {
		go func(i uint32) {
			fn := func() {
				fmt.Println(i)
			}
			trigger(i, fn)
		}(i)
	}
	trigger(10, func(){})
}