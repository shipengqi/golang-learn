---
title: map
---
# map
`map` 是一个无序的 `key/value` 对的集合。**`map` 是引用类型**。这意味着它拥有对底层数据结构的引用，
就像指针一样。它底层的数据结构是 `hash table` 或 `hash map`。

`map` 作为引用类型是非常好的，因为无论 `map` 有多大，都只会有一个副本。

定义 `map`，使用 `map` 关键字：
```go
/* 声明变量，默认 map 是 nil */
var 变量名 map[键类型]值类型

/* 使用 make 函数 */
变量名 := make(map[键类型]值类型)

/* 字面值的语法创建 */
变量名 := map[键类型]值类型{
  key1: value1,
  key2: value2,
  ...
}
```

一个 `map` 在未初始化之前默认为 `nil`。
通过索引下标 `key` 来访问 `map` 中对应的 `value`
```go
age, ok := ages["bob"]
if !ok { /* "bob" is not a key in this map; age == 0. */ }
```
`ok` 表示操作结果，是一个布尔值。**这叫做 `ok-idiom` 模式，就是在多返回值中返回一个 `ok` 布尔值，表示是否操作
成功**。

使用 `map` 过程中需要注意的几点：
- **`map` 是无序的，每次打印出来的 `map` 都会不一样**，它不能通过 `index` 获取，而必须通过 `key` 获取
- `map` 的长度是不固定的，也就是和 `slice` 一样，也是一种引用类型
- 内置的 `len` 函数同样适用于 `map`，返回 `map` 拥有的 `key` 的数量
- `map` 的值可以很方便的修改，通过 `numbers["one"]=11` 可以很容易的把 `key` 为 `one` 的字典值改为 11
- **`map` 和其他基本型别不同，它不是 `thread-safe` 的**，在多个 `go-routine` 存取时，必须使用 `mutex lock` 机制

#### delete()
`delete` 函数删除 `map` 元素。
```go
delete(mapName, key)
```

#### 遍历
可以使用 `for range` 遍历 `map`：
```go
for key, value := range mapName {
	fmt.Println(mapName[key])
}
```
**`Map` 的迭代顺序是不确定的。可以先使用 `sort` 包排序**。

#### map 为什么是无序的
编译器对于 slice 和 map 的循环迭代有不同的实现方式，`for` 遍历 map，调用了两个方法：
- `runtime.mapiterinit`
- `runtime.mapiternext`

```go
func mapiterinit(t *maptype, h *hmap, it *hiter) {
	...
	it.t = t
	it.h = h
	it.B = h.B
	it.buckets = h.buckets
	if t.bucket.kind&kindNoPointers != 0 {
		h.createOverflow()
		it.overflow = h.extra.overflow
		it.oldoverflow = h.extra.oldoverflow
	}

	r := uintptr(fastrand())
	if h.B > 31-bucketCntBits {
		r += uintptr(fastrand()) << 31
	}
	it.startBucket = r & bucketMask(h.B)
	it.offset = uint8(r >> h.B & (bucketCnt - 1))
	it.bucket = it.startBucket
    ...

	mapiternext(it)
}
```

`fastrand` 部分，它是一个生成随机数的方法，它生成了随机数。用于决定从哪里开始循环迭代。
因此**每次 `for range map` 的结果都是不一样的。那是因为它的起始位置根本就不固定**。

#### map 的键类型不能是哪些类型
`map` 的键和元素的最大不同在于，前者的类型是受限的，而后者却可以是任意类型的。

**`map` 的键类型不可以是函数类型、字典类型和切片类型**。

为什么？

Go 语言规范规定，在**键类型的值之间必须可以施加操作符 `==` 和 `!=`**。换句话说，键类型的值必须要支持判等操作。由于
函数类型、字典类型和切片类型的值并不支持判等操作，所以字典的键类型不能是这些类型。

另外，如果键的类型是接口类型的，那么键值的实际类型也不能是上述三种类型，否则在程序运行过程中会引发 panic（即运行时恐慌）。
```go
var badMap2 = map[interface{}]int{
"1":   1,
[]int{2}: 2, // 这里会引发 panic。
3:    3,
}
```

#### 优先考虑哪些类型作为字典的键类型
求哈希和判等操作的速度越快，对应的类型就越适合作为键类型。

对于所有的基本类型、指针类型，以及数组类型、结构体类型和接口类型，Go 语言都有一套算法与之对应。这套算法中就包含了哈希和判等。
以求哈希的操作为例，宽度越小的类型速度通常越快。对于布尔类型、整数类型、浮点数类型、复数类型和指针类型来说都是如此。对于字
符串类型，由于它的宽度是不定的，所以要看它的值的具体长度，长度越短求哈希越快。

类型的宽度是指它的单个值需要占用的字节数。比如，`bool`、`int8` 和 `uint8` 类型的一个值需要占用的字节数都是 1，因此这
些类型的宽度就都是 1。


#### 在值为 nil 的字典上执行读写操作会成功吗
当我们仅声明而不初始化一个字典类型的变量的时候，它的值会是 `nil`。如果你尝试使用一个 `nil` 的 `map`，你会
得到一个 `nil` 指针异常，这将导致程序终止运行。所以不应该初始化一个空的 map 变量，比如 `var m map[string]string`。

**除了添加键 - 元素对，我们在一个值为 `nil` 的字典上做任何操作都不会引起错误**。当我们试图在一个值为 `nil` 的字典中
添加键 - 元素对的时候，Go 语言的运行时系统就会立即抛出一个 panic。

可以先使用 `make` 函数初始化，或者 `dictionary = map[string]string{}`。这两种方法都可以创建一个空的 `hash map`
 并指向 `dictionary`。这确保永远不会获得 `nil 指针异常`。
