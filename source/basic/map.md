---
title: map
---

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

对于所有的基本类型、指针类型，以及数组类型、结构体类型和接口类型，Go 语言都有一套算法与之对应。这套算法中就包含了哈希和判等。以求哈希的操作为例，宽度越小的类型速度通常越快。对于布尔类型、整数类型、浮点数类型、复数类型和指针类型来说都是如此。对于字符串类型，由于它的宽度是不定的，所以要看它的值的具体长度，长度越短求哈希越快。

类型的宽度是指它的单个值需要占用的字节数。比如，`bool`、`int8` 和 `uint8` 类型的一个值需要占用的字节数都是 1，因此这些类型的宽度就都是 1。

#### 在值为 nil 的字典上执行读写操作会成功吗

当我们仅声明而不初始化一个字典类型的变量的时候，它的值会是 `nil`。如果你尝试使用一个 `nil` 的 `map`，你会得到一个 `nil` 指针异常，这将导致程序终止运行。所以不应该初始化一个空的 map 变量，比如 `var m map[string]string`。

**除了添加键 - 元素对，我们在一个值为 `nil` 的字典上做任何操作都不会引起错误**。当我们试图在一个值为 `nil` 的字典中添加键 - 元素对的时候，Go 语言的运行时系统就会立即抛出一个 panic。

可以先使用 `make` 函数初始化，或者 `dictionary = map[string]string{}`。这两种方法都可以创建一个空的 `hash map` 并指向 `dictionary`。这确保永远不会获得 `nil 指针异常`。

## hash 表

要实现一个性能优异的哈希表，需要注意两个关键点 —— **哈希函数和冲突解决方法**。

哈希函数的选择在很大程度上能够决定哈希表的读写性能。在理想情况下，哈希函数应该能够将不同键映射到不同的索引上，这要求哈希函数的输出范围大于输入范围，但是由于键的数量会远远大于映射的范围，所以在实际使用时，这个理想的效果是不可能实现的。

比较实际的方式是让哈希函数的结果能够尽可能的均匀分布，然后通过工程上的手段解决哈希碰撞的问题。哈希函数映射的结果一定要尽可能均匀，结果不均匀的哈希函数会带来更多的哈希冲突以及更差的读写性能。

如果使用结果分布较为均匀的哈希函数，那么哈希的增删改查的时间复杂度为 `O(1)`；但是如果哈希函数的结果分布不均匀，那么所有操作的时间复杂度可能会达到 `O(n)` （为什么是  `O(n)` ，如果使用拉链法解决哈希冲突，极端情况下，hash 函数的结构都在一个索引的链表上，复杂度就是 `O(n)`），由此看来，使用好的哈希函数是至关重要的。

常见解决哈希冲突方法的就是开放寻址法和拉链法。

开放寻址法2是一种在哈希表中解决哈希碰撞的方法，这种方法的核心思想是**依次探测和比较数组中的元素以判断目标键值对是否存在于哈希表中**。

开放寻址法中对性能影响最大的是**装载因子**，它是数组中元素的数量与数组大小的比值。随着装载因子的增加，线性探测的平均用时就会逐渐增加，这会影响哈希表的读写性能。当装载率超过 70% 之后，哈希表的性能就会急剧下降，而一旦装载率达到 100%，整个哈希表就会完全失效，这时查找和插入任意元素的时间复杂度都是 `O(n)` 的，这时需要遍历数组中的全部元素，所以在实现哈希表时一定要关注装载因子的变化。

拉链法是哈希表最常见的实现方法。一般会使用**数组加上链表**，不过一些编程语言会在拉链法的哈希中引入红黑树以优化性能，拉链法会使用链表数组作为哈希底层的数据结构。

![](separate-chaing-and-set.png)

上图所示，当我们需要将一个键值对 (Key6, Value6) 写入哈希表时，键值对中的键 Key6 都会先经过一个哈希函数，哈希函数返回的哈希会帮助我们选择一个桶，和开放地址法一样，选择桶的方式是直接对哈希返回的结果取模：

```go
index := hash("Key6") % array.len
```

选择了 2 号桶后就可以遍历当前桶中的链表了，在遍历链表的过程中会遇到以下两种情况：

1. 找到键相同的键值对 — 更新键对应的值；
2. 没有找到键相同的键值对 — 在链表的末尾追加新的键值对；

## 数据结构

`runtime.hmap` 是最核心的结构体：

```go
type hmap struct {
 count     int        // 哈希表中的元素数量
 flags     uint8      // 状态标识，主要是 goroutine 写入和扩容机制的相关状态控制。并发读写的判断条件之一就是该值
 B         uint8      // 哈希表持有的 buckets 数量，但是因为哈希表中桶的数量都 2 的倍数，所以该字段会存储对数，也就是 len(buckets) == 2^B
 noverflow uint16     // 溢出桶的数量
 hash0     uint32     // 哈希的种子，它能为哈希函数的结果引入随机性，这个值在创建哈希表时确定，并在调用哈希函数时作为参数传入

 buckets    unsafe.Pointer  // 当前桶
 oldbuckets unsafe.Pointer  // 哈希在扩容时用于保存之前 buckets 的字段，它的大小是当前 buckets 的一半
 nevacuate  uintptr         // 迁移进度

 extra *mapextra
}

type mapextra struct {
 overflow    *[]*bmap   为 hmap.buckets （当前）溢出桶的指针地址
 oldoverflow *[]*bmap   为 hmap.oldbuckets （旧）溢出桶的指针地址
 nextOverflow *bmap     为空闲溢出桶的指针地址
}
```

![](hmap-and-buckets.png)

`runtime.hmap` 的桶是 `runtime.bmap`。每一个 **`runtime.bmap` 都能存储 8 个键值对**，当哈希表中存储的数据过多，单个桶无法装满时就会使用 `extra.nextOverflow` 中桶存储溢出的数据。

述两种不同的桶在内存中是连续存储的，我们在这里将它们分别称为**正常桶**和**溢出桶**。黄色的就是正常桶，绿色的是溢出桶。**溢出桶能够减少扩容的频率**。

```go
type bmap struct {
 tophash [bucketCnt]uint8
}
```

`tophash` 存储了**键的哈希的高 8 位，通过比较不同键的哈希的高 8 位可以减少访问键值对次数以提高性能**。

```go
type bmap struct {
    topbits  [8]uint8
    keys     [8]keytype
    values   [8]valuetype
    pad      uintptr
    overflow uintptr
}

```

存储 k 和 v 的载体并不是用 `k/v/k/v/k/v/k/v` 的模式，而是 `k/k/k/k/v/v/v/v` 的形式去存储。这是为什么呢？

例如一个 map `map[int64]int8`，如果按照 `k/v` 的形式存放 int64 的 key 占用 8 个字节，最然值 int8 只占用一个字节，但是却需要 7 个填充字节来做内存对齐，就会浪费大量内存空间。

随着哈希表存储的数据逐渐增多，我们会扩容哈希表或者使用额外的桶存储溢出的数据，不会让单个桶中的数据超过 8 个，不过溢出桶只是临时的解决方案，创建过多的溢出桶最终也会导致哈希的扩容。

## 访问

`hash[key]` 以及类似的操作都会被转换成哈希的 OINDEXMAP 操作，中间代码生成阶段会在 `cmd/compile/internal/gc.walkexpr` 函数中将这些 OINDEXMAP 操作转换成如下的代码：

```go
v     := hash[key] // => v     := *mapaccess1(maptype, hash, &key)
v, ok := hash[key] // => v, ok := mapaccess2(maptype, hash, &key)
```

赋值语句左侧接受参数的个数会决定使用的运行时方法：

- 当接受一个参数时，会使用 `runtime.mapaccess1`，该函数仅会返回一个指向目标值的指针；
- 当接受两个参数时，会使用 `runtime.mapaccess2`，除了返回目标值之外，它还会返回一个用于表示当前键对应的值是否存在的 bool 值：

`runtime.mapaccess1` 会先通过哈希表设置的哈希函数、种子获取当前键对应的哈希，再通过 `runtime.bucketMask` 和 `runtime.add` 拿到该键值对所在的桶序号和哈希高位的 8 位数字。

```go
func mapaccess1(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer {
 alg := t.key.alg
 hash := alg.hash(key, uintptr(h.hash0))
 m := bucketMask(h.B)
 b := (*bmap)(add(h.buckets, (hash&m)*uintptr(t.bucketsize)))
 top := tophash(hash)
bucketloop:
 for ; b != nil; b = b.overflow(t) {
  for i := uintptr(0); i < bucketCnt; i++ {
   if b.tophash[i] != top {
    if b.tophash[i] == emptyRest {
     break bucketloop
    }
    continue
   }
   k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))
   if alg.equal(key, k) {
    v := add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.keysize)+i*uintptr(t.valuesize))
    return v
   }
  }
 }
 return unsafe.Pointer(&zeroVal[0])
}
```

bucketloop 循环中，哈希会依次遍历正常桶和溢出桶中的数据，它先会比较哈希的高 8 位和桶中存储的 tophash，后比较传入的和桶中的值以加速数据的读写。用于选择桶序号的是哈希的最低几位，而用于加速访问的是哈希的高 8 位，这种设计能够减少同一个桶中有大量相等 tophash 的概率影响性能。

![](hashmap-mapaccess.png)

每一个桶都是一整片的内存空间，当发现桶中的 tophash 与传入键的 tophash 匹配之后，我们会通过指针和偏移量获取哈希中存储的键 `keys[0]` 并与 key 比较，如果两者相同就会获取目标值的指针 `values[0]` 并返回。

判断是否正在发生扩容（h.oldbuckets 是否为 nil），若正在扩容，则到老的 buckets 中查找（因为 buckets 中可能还没有值，搬迁未完成），若该 bucket 已经搬迁完毕。则到 buckets 中继续查找

## 写入

当形如 `hash[k]` 的表达式出现在赋值符号左侧时，该表达式也会在编译期间转换成 `runtime.mapassign` 函数的调用

```go
func mapassign(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer {
 alg := t.key.alg
 hash := alg.hash(key, uintptr(h.hash0))

 h.flags ^= hashWriting

again:
 bucket := hash & bucketMask(h.B)
 b := (*bmap)(unsafe.Pointer(uintptr(h.buckets) + bucket*uintptr(t.bucketsize)))
 top := tophash(hash)
```

通过遍历比较桶中存储的 tophash 和键的哈希，如果找到了相同结果就会返回目标位置的地址。其中 inserti 表示目标元素的在桶中的索引，insertk 和 val 分别表示键值对的地址，获得目标地址之后会通过算术计算寻址获得键值对 k 和 val：

```go
 var inserti *uint8
 var insertk unsafe.Pointer
 var val unsafe.Pointer
bucketloop:
 for {
  for i := uintptr(0); i < bucketCnt; i++ {
   if b.tophash[i] != top {
    if isEmpty(b.tophash[i]) && inserti == nil {
     inserti = &b.tophash[i]
     insertk = add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))
     val = add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.keysize)+i*uintptr(t.valuesize))
    }
    if b.tophash[i] == emptyRest {
     break bucketloop
    }
    continue
   }
   k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))
   if !alg.equal(key, k) {
    continue
   }
   val = add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.keysize)+i*uintptr(t.valuesize))
   goto done
  }
  ovf := b.overflow(t)
  if ovf == nil {
   break
  }
  b = ovf
 }
```

 for 循环会依次遍历正常桶和溢出桶中存储的数据，整个过程会分别判断 tophash 是否相等、key 是否相等，遍历结束后会从循环中跳出。

 ![](hashmap-overflow-bucket.png)

如果当前桶已经满了，哈希会调用 `runtime.hmap.newoverflow` 创建新桶或者使用 `runtime.hmap` 预先在 `noverflow` 中创建好的桶来保存数据，新创建的桶不仅会被追加到已有桶的末尾，还会增加哈希表的 `noverflow` 计数器。

```go
 if inserti == nil {
  newb := h.newoverflow(t, b)
  inserti = &newb.tophash[0]
  insertk = add(unsafe.Pointer(newb), dataOffset)
  val = add(insertk, bucketCnt*uintptr(t.keysize))
 }

 typedmemmove(t.key, insertk, key)
 *inserti = top
 h.count++

done:
 return val
}
```

## 扩容

```bash
装载因子 := 元素数量 ÷ 桶数量
```

`runtime.mapassign` 函数会在以下两种情况发生时触发哈希的扩容：

- 装载因子已经超过 6.5；
- 哈希使用了太多溢出桶；

哈希的扩容不是一个原子的过程，所以 `runtime.mapassign` 还需要**判断当前哈希是否已经处于扩容状态，避免二次扩容造成混乱**。

根据触发的条件不同扩容的方式分成两种，如果这次扩容是溢出的桶太多导致的，那么这次扩容就是**等量扩容 `sameSizeGrow`**，sameSizeGrow 是一种特殊情况下发生的扩容，当我们持续向哈希中插入数据并将它们全部删除时，如果哈希表中的数据量没有超过阈值，就会不断积累溢出桶造成缓慢的内存泄漏。runtime: limit the number of map overflow buckets 引入了 **sameSizeGrow 通过复用已有的哈希扩容机制解决该问题，一旦哈希中出现了过多的溢出桶，它会创建新桶保存数据，垃圾回收会清理老的溢出桶并释放内存**。

扩容的入口是 `runtime.hashGrow`：

```go
func hashGrow(t *maptype, h *hmap) {
 bigger := uint8(1)
 if !overLoadFactor(h.count+1, h.B) {
  bigger = 0
  h.flags |= sameSizeGrow
 }
 oldbuckets := h.buckets
 newbuckets, nextOverflow := makeBucketArray(t, h.B+bigger, nil)

 h.B += bigger
 h.flags = flags
 h.oldbuckets = oldbuckets
 h.buckets = newbuckets
 h.nevacuate = 0
 h.noverflow = 0

 h.extra.oldoverflow = h.extra.overflow
 h.extra.overflow = nil
 h.extra.nextOverflow = nextOverflow
}
```

哈希在扩容的过程中会通过 runtime.makeBucketArray 创建一组新桶和预创建的溢出桶，随后将原有的桶数组设置到 oldbuckets 上并将新的空桶设置到 buckets 上，溢出桶也使用了相同的逻辑更新，下图展示了触发扩容后的哈希：

![](hashmap-hashgrow.png)

为什么是增量扩容？

“渐进式”地方式，原有的 key 并不会一次性搬迁完毕，每次最多只会搬迁 2 个 bucket。

如果是全量扩容的话，那问题就来了。假设当前 hmap 的容量比较大，直接全量扩容的话，就会导致扩容要花费大量的时间和内存，导致系统卡顿，最直观的表现就是慢。

```go
type evacDst struct {
 b *bmap  // 当前目标桶
 i int    // 当前目标桶存储的键值对数量
 k unsafe.Pointer  // 指向当前 key 的内存地址
 v unsafe.Pointer  // 指向当前 value 的内存地址
}
```

evacDst 是迁移中的基础数据结构

```go
func evacuate(t *maptype, h *hmap, oldbucket uintptr) {
 b := (*bmap)(add(h.oldbuckets, oldbucket*uintptr(t.bucketsize)))
 newbit := h.noldbuckets()
 if !evacuated(b) {
  var xy [2]evacDst
  x := &xy[0]
  x.b = (*bmap)(add(h.buckets, oldbucket*uintptr(t.bucketsize)))
  x.k = add(unsafe.Pointer(x.b), dataOffset)
  x.v = add(x.k, bucketCnt*uintptr(t.keysize))

  if !h.sameSizeGrow() {
   y := &xy[1]
   y.b = (*bmap)(add(h.buckets, (oldbucket+newbit)*uintptr(t.bucketsize)))
   y.k = add(unsafe.Pointer(y.b), dataOffset)
   y.v = add(y.k, bucketCnt*uintptr(t.keysize))
  }

  for ; b != nil; b = b.overflow(t) {
            ...
  }

  if h.flags&oldIterator == 0 && t.bucket.kind&kindNoPointers == 0 {
   b := add(h.oldbuckets, oldbucket*uintptr(t.bucketsize))
   ptr := add(b, dataOffset)
   n := uintptr(t.bucketsize) - dataOffset
   memclrHasPointers(ptr, n)
  }
 }

 if oldbucket == h.nevacuate {
  advanceEvacuationMark(h, t, newbit)
 }
}
```

计算并得到 oldbucket 的 bmap 指针地址
计算 hmap 在增长之前的桶数量
判断当前的迁移（搬迁）状态，以便流转后续的操作。若没有正在进行迁移 !evacuated(b) ，则根据扩容的规则的不同，当规则为等量扩容 sameSizeGrow 时，只使用一个 evacDst 桶用于分流。而为双倍扩容时，就会使用两个 evacDst 进行分流操作
当分流完毕后，需要迁移的数据都会通过 typedmemmove 函数迁移到指定的目标桶上
若当前不存在 flags 使用标志、使用 oldbucket 迭代器、bucket 不为指针类型。则取消链接溢出桶、清除键值
在最后 advanceEvacuationMark 函数中会对迁移进度 hmap.nevacuate 进行累积计数，并调用 bucketEvacuated 对旧桶 oldbuckets 进行不断的迁移。直至全部迁移完毕。那么也就表示扩容完毕了，会对 hmap.oldbuckets 和 h.extra.oldoverflow 进行清空
