## container 包
Go 语言的链表实现在标准库的 `container/list` 代码包中。

这个代码包中有两个公开的程序实体——`List` 和 `Element`，`List` 实现了一个双向链表（以下简称链表），
而 `Element` 则代表了链表中元素的结构。

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

- `InsertBefore` 和 `InsertAfter` 方法分别用于在指定的元素之前和之后插入新元素
- `PushFront` 和 `PushBack` 方法则分别用于在链表的最前端和最后端插入新元素。
