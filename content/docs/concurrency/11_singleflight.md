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

`singleflight.Group` 提供了三个方法：

- `Do`：接受两个参数，第一个参数是一个 key，第二个参数是一个函数。同一个 key 对应的函数，在同一时间只会有一个在执行，其他的并发执行的请求会等待。当第一个执行的函数返回结果
其他的并发请求会使用这个结果。
- `DoChan`：和 `Do` 方法差不多，只不过是返回一个 channel，当执行的函数返回结果时，就可以从这个 channel 中接收这个结果。
- `Forget`：在 `Group` 的映射表中删除某个 key。接下来这个 key 的请求就不会等待前一个未完成的函数的返回结果了。

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

	//  val 和 err 只会在执行传入的函数时赋值一次并在 WaitGroup.Wait 返回时被读取
	val interface{}
	err error

	// 抑制的请求数量
	dups  int
	// 用于同步结果
	chans []chan<- Result
}
```

`Do` 的实现：

```go
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok { // 存在相同的 key
		c.dups++
		g.mu.Unlock()
		c.wg.Wait() // 等待这个 key 的第一个请求完成
		return c.val, c.err, true // 使用 key 的请求结果
	}
    // 第一个请求，创建一个 call
	c := new(call) 
	c.wg.Add(1)
    // 将 key 放到 map
	g.m[key] = c
	g.mu.Unlock()

	// 执行函数
	g.doCall(c, key, fn)
	return c.val, c.err, c.dups > 0
}

func (g *Group) doCall(c *call, key string, fn func() (interface{}, error)) {
	// 执行函数
	// 将函数的返回值赋值给 c.val 和 c.err
	c.val, c.err = fn()
	// 当前函数已经执行完成，通知所有等待结果的 goroutine 可以从 call 结构体中取出返回值并返回了
	c.wg.Done()

	g.mu.Lock()
	// 从 map 中删除已经执行一次的 key
	delete(g.m, key)
	// 将结果通过 channel 同步给使用 DoChan 的 goroutine
	for _, ch := range c.chans {
		ch <- Result{c.val, c.err, c.dups > 0}
	}
	g.mu.Unlock()
}
```