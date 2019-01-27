package main

import(
	"fmt"
	"runtime"
	"sync"
)

func main() {
	p := &sync.Pool{
		New: func() interface{} {
			return 0
		},
	}

	a := p.Get().(int)
	p.Put(1)
	runtime.GC()
	b := p.Get().(int)
	fmt.Println(a, b)
}