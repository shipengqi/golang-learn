package main

import (
	"container/ring"
	"fmt"
)

func main() {

	const rLen = 3

	// 创建新的 Ring
	r := ring.New(rLen)

	for i := 0; i < rLen; i++ {
		r.Value = i
		r = r.Next()
	}

	fmt.Printf("Length of ring: %d\n", r.Len()) // Length of ring: 3

	// 该匿名函数用来打印 Ring 中的数据
	printRing := func(v interface{}) {
		fmt.Print(v, " ")
	}

	r.Do(printRing) // 0 1 2
	fmt.Println()

	// 将 r 之后的第二个元素的值乘以 2
	r.Move(2).Value = r.Move(2).Value.(int) * 2

	r.Do(printRing) // 0 1 4
	fmt.Println()

	// 删除 r 与 r+2 之间的元素，即删除 r+1
	// 返回删除的元素组成的Ring的指针
	result := r.Link(r.Move(2))

	r.Do(printRing) // 0 4
	fmt.Println()

	result.Do(printRing) // 1
	fmt.Println()

	another := ring.New(rLen)
	another.Value = 7
	another.Next().Value = 8 // 给 another + 1 表示的元素赋值，即第二个元素
	another.Prev().Value = 9 // 给 another - 1 表示的元素赋值，即第三个元素

	another.Do(printRing) // 7 8 9
	fmt.Println()

	// 插入another到r后面，返回插入前r的下一个元素
	result = r.Link(another)

	r.Do(printRing) // 0 7 8 9 4
	fmt.Println()

	result.Do(printRing) // 4 0 7 8 9
	fmt.Println()

	// 删除r之后的三个元素，返回被删除元素组成的Ring的指针
	result = r.Unlink(3)

	r.Do(printRing) // 0 4
	fmt.Println()

	result.Do(printRing) // 7 8 9
	fmt.Println()
}