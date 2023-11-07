---
title: SingleFlight
weight: 11
---

# SingleFlight

Go 的扩展库 `golang.org/x/sync` 提供了 `singleflight` 包，它的作用在处理多个 goroutine 同时调用同一个函数的时候，只让一个 goroutine 去调用这个函数，等到这个 goroutine 返回结果时，再把结
果返回给这几个 goroutine，这样可以减少并发调用的数量。

一个常见的使用场景：在使用 Redis 对数据库中的数据进行缓存，如果发生缓存击穿，大量的流量都会打到后端数据库上，导致后端服务响应延时等问题。
`singleflight` 可以将对同一个 key 的多个请求合并为一个，减轻后端服务的压力。

## 使用

```go
package main

import (
	"fmt"
	"time"
	
	"golang.org/x/sync/singleflight"
)

func GetValueFromRedis(key string) string {
	fmt.Println("query ...")
	time.Sleep(10 * time.Second) // 模拟一个比较耗时的操作
	return "singleflight demo"
}

func main() {
	requestGroup := new(singleflight.Group)

	cachekey := "demokey"
	go func() {
		v1, _, shared := requestGroup.Do(cachekey, func() (interface{}, error) {
			ret := GetValueFromRedis(cachekey)
			return ret, nil
		})
		fmt.Printf("1st call: v1: %v, shared: %v\n", v1, shared)
	}()

	time.Sleep(2 * time.Second)

	// 重复查询 key，第一次查询还未结束
	v2, _, shared := requestGroup.Do(cachekey, func() (interface{}, error) {
		ret := GetValueFromRedis(cachekey)
		return ret, nil
	})
	fmt.Printf("2nd call: v2:%v, shared:%v\n", v2, shared)
}
```

输出：

```
query ...
1st call: v1: singleflight demo, shared:true
2nd call: v2: singleflight demo, shared:true
```

`query ...` 只打印了一次，请求被合并了。


## 原理

`singleflight.Group` 的结构体：

```go
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// 代表一个正在处理的请求，或者已经处理完的请求
type call struct {
	wg sync.WaitGroup

	// 这个字段代表处理完的值，在 waitgroup 完成之前只会写一次
	// waitgroup 完成之后就读取这个值
	val interface{}
	err error

	// These fields are read and written with the singleflight
	// mutex held before the WaitGroup is done, and are read but
	// not written after the WaitGroup is done.
	dups  int
	chans []chan<- Result
}
```
