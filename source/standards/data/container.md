---
title: container
---

# container

`container` 该包实现了三个复杂的数据结构：链表，环，堆。也就是说使用这三个数据结构的时候不需要再从头开始写算法了。

##  链表

链表就是一个有 `prev` 和 `next` 指针的数组了。
`container` 包中有两个公开的结构—— `List` 和 `Element`，`List` 实现了一个双向链表（简称链表），
而 `Element` 则代表了链表中元素的结构。

```go
type Element struct {
	next, prev *Element  // 上一个元素和下一个元素
	list *List  // 元素所在链表
	Value interface{}  // 元素
}

type List struct {
	root Element  // 链表的根元素
	len  int      // 链表的长度
}
```
List的四种方法:
- `MoveBefore` 方法和 `MoveAfter` 方法，它们分别用于把给定的元素移动到另一个元素的前面和后面。
- `MoveToFront` 方法和 `MoveToBack` 方法，分别用于把给定的元素移动到链表的最前端和最后端。


```go
// moves element "e" to its new position before "mark".
func (l *List) MoveBefore(e, mark *Element)
// moves element "e" to its new position after "mark".
func (l *List) MoveAfter(e, mark *Element)

// moves element "e" to the front of list "l".
func (l *List) MoveToFront(e *Element)
// moves element "e" to the back of list "l".
func (l *List) MoveToBack(e *Element)
```

“给定的元素”都是 `*Element` 类型。

如果我们自己生成这样的值，然后把它作为“给定的元素”传给链表的方法，那么会发生什么？链表会接受它吗？

不会接受，这些方法将不会对链表做出任何改动。因为我们自己生成的 `Element` 值并不在链表中，所以也就谈不上“在链表中移动元素”。

- `InsertBefore` 和 `InsertAfter` 方法分别用于在指定的元素之前和之后插入新元素。
- `PushFront` 和 `PushBack` 方法则分别用于在链表的最前端和最后端插入新元素。

示例：

```go
package main

import (
	"container/list"
	"fmt"
)

func main() {
    list := list.New()
    list.PushBack(1)
    list.PushBack(2)

    fmt.Printf("len: %v\n", list.Len())
    fmt.Printf("first: %#v\n", list.Front())
    fmt.Printf("second: %#v\n", list.Front().Next())
}

output:
len: 2
first: &list.Element{next:(*list.Element)(0x2081be1b0), prev:(*list.Element)(0x2081be150), list:(*list.List)(0x2081be150), Value:1}
second: &list.Element{next:(*list.Element)(0x2081be150), prev:(*list.Element)(0x2081be180), list:(*list.List)(0x2081be150), Value:2}
```

List 的其他方法：
```go
type Element
    func (e *Element) Next() *Element
    func (e *Element) Prev() *Element
type List
    func New() *List
    func (l *List) Back() *Element   // 最后一个元素
    func (l *List) Front() *Element  // 第一个元素
    func (l *List) Init() *List  // 链表初始化
    func (l *List) InsertAfter(v interface{}, mark *Element) *Element // 在某个元素后插入
    func (l *List) InsertBefore(v interface{}, mark *Element) *Element  // 在某个元素前插入
    func (l *List) Len() int // 在链表长度
    func (l *List) PushBackList(other *List)  // 在队列最后插入接上新队列
    func (l *List) PushFrontList(other *List) // 在队列头部插入接上新队列
    func (l *List) Remove(e *Element) interface{} // 删除某个元素
```

## 环

环的结构有点特殊，环的尾部就是头部，指向环形链表任一元素的指针都可以作为整个环形链表看待。
它不需要像 List 一样保持 List 和 Element 两个结构，只需要保持一个结构就行。

```go
type Ring struct {
	next, prev *Ring
	Value      interface{}
}
```

初始化环的时候，需要定义好环的大小，然后对环的每个元素进行赋值。环还提供一个 `Do` 方法，能遍历一遍环，对每个元素执行
一个 `function`。

示例：

```go
package main

import (
	"container/ring"
	"fmt"
)

func main() {
    ring := ring.New(3)

    for i := 1; i <= 3; i++ {
        ring.Value = i
        ring = ring.Next()
    }

    // 计算 1+2+3
    s := 0
    ring.Do(func(p interface{}){
        s += p.(int)
    })
    fmt.Println("sum is", s)
}

output:
sum is 6
```

ring 提供的方法有

```go
type Ring
    func New(n int) *Ring // 创建一个长度为 n 的环形链表
    func (r *Ring) Do(f func(interface{})) // 遍历环形链表中的每一个元素 x 进行 f(x) 操作
    func (r *Ring) Len() int // 获取环形链表长度
    
    // 如果 r 和 s 在同一环形链表中，则删除 r 和 s 之间的元素，
    // 被删除的元素组成一个新的环形链表，返回值为该环形链表的指针（即删除前，r->Next() 表示的元素）
    // 如果 r 和 s 不在同一个环形链表中，则将 s 插入到 r 后面，返回值为
    // 插入 s 后，s 最后一个元素的下一个元素（即插入前，r->Next() 表示的元素）
    func (r *Ring) Link(s *Ring) *Ring

    func (r *Ring) Move(n int) *Ring // 移动 n % r.Len() 个位置，n 正负均可
    func (r *Ring) Next() *Ring // 返回下一个元素
    func (r *Ring) Prev() *Ring // 返回前一个元素
    func (r *Ring) Unlink(n int) *Ring // 删除 r 后面的 n % r.Len() 个元素
```

## 堆