---
title: 哈希表
weight: 4
---

`map` 是一个无序的 `key/value` 对的集合，同一个 key 只会出现一次。

## 哈希表的设计原理


哈希表其实是数组的扩展。哈希表是利用数组可以根据下标随机访问（时间复杂度是 `O(1)`）这一特性来实现快速查找的。

### 哈希函数

哈希表是通过**哈希函数**将 key 转化为数组的下标，然后将数据存储在数组下标对应的位置。查询时，也是同样的使用哈希函数计算出数组下标，从下标对应的位置取出数据。

![map-seek-addr](https://raw.gitcode.com/shipengqi/illustrations/blobs/20553472dd8c6502c7cd1a510d4dca24f5198397/map-seek-addr.png)

哈希函数的基本要求：

1. 哈希函数计算出来的值是一个非负整数。
2. 如果 `key1 == key2` 那么 `hash(key1) == hash(key2)`
3. 如果 `key1 != key2` 那么 `hash(key1) != hash(key2)`


第三点，想要实现一个不同的 key 对应的哈希值绝对不一样的哈希函数，几乎是不可能的，也就说无法避免**哈希冲突**。

常用的处理哈希冲突的方法有两种：**开放寻址法**和**链表法**。

### 开放寻址法

开放寻址法核心思想是，如果出现了哈希冲突，就重新探测一个空闲位置，将其插入。

![map-seek-addr](https://raw.gitcode.com/shipengqi/illustrations/blobs/20553472dd8c6502c7cd1a510d4dca24f5198397/map-seek-addr.png)

上图蓝色表示已经插入的元素，`key9` 哈希后得到的数组下标为 6，但是已经有数据了，产生了冲突。那么就按顺序向后查找直到找到一个空闲的位置，如果到数组的尾部都没有找到空闲的位置，就从头开始继续找。上图最终找到位置 1 并插入元素。

查找的逻辑和插入类似，从哈希函数计算出来的下标位置开始查找，比较数组中下标位置的元素和要查找的元素。如果相等，则说明就是要找的元素；否则就顺序往后依次查找。直到找到数组中的空闲位置，还没有找到，就说明要查找的元素并没有在哈希表中。

可以看出当数组中空闲位置不多的时候，哈希冲突的概率就会大大提高。**装载因子**（load factor）就是用来表示空位的多少。

```
装载因子=已插入的元素个数/哈希表的长度
```

装载因子越大，说明空闲位置越少，冲突越多，哈希表的性能会下降。

### 链表法

链表法是最常见的哈希冲突的解决办法。在哈希表中，每个桶（bucket）会对应一条链表，所有哈希值相同的元素都放到相同桶对应的链表中。

![map-link](https://raw.gitcode.com/shipengqi/illustrations/blobs/f242cc1652d75a584b78212fe4b077d3e0f22b72/map-link.png)

插入时，哈希函数计算后得出存放在几号桶，然后遍历桶中的链表了：

- 找到键相同的键值对，则更新键对应的值；
- 没有找到键相同的键值对，则在链表的末尾追加新的键值对

链表法实现的哈希表的装载因子：

```
装载因子=已插入的元素个数/桶数量
```

## Go map 原理

表示 map 的结构体是 `hmap`：

```go
// src/runtime/map.go
type hmap struct {
    // 哈希表中的元素数量 
	count     int        
	// 状态标识，主要是 goroutine 写入和扩容机制的相关状态控制。并发读写的判断条件之一就是该值
	flags     uint8
	// 哈希表持有的 buckets 数量，但是因为哈希表中桶的数量都 2 的倍数，
	// 所以该字段会存储对数，也就是 len(buckets) == 2^B
	B         uint8
	// 溢出桶的数量
	noverflow uint16     
    // 哈希种子，它能为哈希函数的结果引入随机性，这个值在创建哈希表时确定，并在调用哈希函数时作为参数传入
	hash0     uint32
    // 指向 buckets 数组，长度为 2^B
	buckets   unsafe.Pointer  
    // 哈希在扩容时用于保存之前 buckets 的字段
	// 等量扩容的时候，buckets 长度和 oldbuckets 相等
	// 双倍扩容的时候，buckets 长度是 oldbuckets 的两倍
	oldbuckets unsafe.Pointer
	// 迁移进度，小于此地址的 buckets 是已迁移完成的
	nevacuate  uintptr         
	extra *mapextra
}

type mapextra struct {
    // hmap.buckets （当前）溢出桶的指针地址
	overflow    *[]*bmap
	// 为 hmap.oldbuckets （旧）溢出桶的指针地址
	oldoverflow *[]*bmap
    // 为空闲溢出桶的指针地址
	nextOverflow *bmap     
}
```

**`hmap.buckets` 就是指向一个 `bmap` 数组**。`bmap` 的结构体：

```go
type bmap struct {
	tophash [bucketCnt]uint8
}

// 编译时，编译器会推导键值对占用内存空间的大小，然后修改 bmap 的结构
type bmap struct {
	topbits  [8]uint8
	keys     [8]keytype
	values   [8]valuetype
	pad      uintptr
	overflow uintptr
}
```

`bmap` 就是桶，**一个桶里面会最多存储 8 个键值对**。

![bmap-struct](https://raw.gitcode.com/shipengqi/illustrations/blobs/1ed637b480a103c526ff108da4e3658c1fcbf9c0/bmap-struct.png)

1. 在桶内，会**根据 key 计算出来的 hash 值的高 8 位来决定 key 存储在桶中的位置**（桶内的键值对，根据类型的大小就可以计算出偏移量）。
2. **key 和 value 是分别放在一块连续的内存，这样做的目的是为了节省内存**。例如一个 `map[int64]int8` 类型的 map，如果按照 `key1/value1/key2/value2 ...` 这样的形式来存储，那么**内存对齐**每个 `key/value` 都需要 padding 7 个字节。分开连续存储的话，就只需要在最后 padding 一次。
3. 每个桶只能存储 8 个 `key/value`，如果有更多的 key 放入当前桶，就需要一个溢出桶，通过 `overflow` 指针连接起来（**链表法**）。

![hmap](https://raw.gitcode.com/shipengqi/illustrations/blobs/7299ec739895315b49422660e585719b3bb38f6d/hmap.png)

### 初始化

初始化 `map`：
```go
hash := map[string]int{
	"1": 2,
	"3": 4,
	"5": 6,
}
hash2 := make(map[string]int, 3)
```

不管是使用字面量还是 `make` 初始化 map，最后都是调用 `makemap` 函数：

```go
func makemap(t *maptype, hint int, h *hmap) *hmap {
	// ...
	// initialize Hmap
	if h == nil {
		h = new(hmap)
	}
	// 获取一个随机的哈希种子
	h.hash0 = fastrand()
	
	// 根据传入的 hint 计算出需要的最小需要的桶的数量
	B := uint8(0)
	for overLoadFactor(hint, B) {
		B++
	}
	h.B = B

	// 初始化 hash table
	// 如果 B 等于 0，那么 buckets 就会在赋值的时候再分配
	// 如果 hint 长度比较大，分配内存会花费长一点
	if h.B != 0 {
		var nextOverflow *bmap
		// makeBucketArray 根据传入的 B 计算出的需要创建的桶数量
		// 并在内存中分配一片连续的空间用于存储数据
		h.buckets, nextOverflow = makeBucketArray(t, h.B, nil)
		if nextOverflow != nil {
			h.extra = new(mapextra)
			h.extra.nextOverflow = nextOverflow
		}
	}

	return h
}
```

预分配的溢出桶和正常桶是在一块连续的内存中。

### 查询

查询 map 中的值：

```go
v := hash[key]
v, ok := hash[key]
```

这两种查询方式会被转换成 [`mapaccess1`](https://github.com/golang/go/blob/8da6405e0db80fa0a4136fb816c7ca2db716c2b2/src/runtime/map.go#L396) 和 `mapaccess2` 函数，两个函数基本一样，不过 `mapaccess2` 函数的返回值多了一个 `bool` 类型。

查询过程：

![map-get](https://raw.gitcode.com/shipengqi/illustrations/blobs/a0add709c3e56df3ee2aeaa95c66ad4015aa2721/map-get.png)

#### 1. 计算哈希值

通过哈希函数和种子获取当前 key 的 64 位的哈希值（64 位机）。以上图哈希值：`11010111 | 110000110110110010001111001010100010010110010101001 │ 00011` 为例。

#### 2. 计算这个 key 要放在哪个桶

**根据哈希值的 `B` （`hmap.B`）个 bit 位来计算，也就是 `00011`，十进制的值是 `3`，那么就是 `3` 号桶。** 

#### 3. 计算这个 key 在桶内的位置

根据哈希值的高 8 位，也就是 `10010111`，十进制的值是 `151`，**先用 `151` 和桶内存储的 `tophash` 比较，`tophash` 一致的话，再比较桶内的存储的 key 和传入的 key，这种方式可以优化桶内的读写速度（`tophash` 不一致就不需要比较了）**。 


```go
// src/runtime/map.go#L434 mapaccess1
for i := uintptr(0); i < bucketCnt; i++ {
	// 先比较 tophash，如果不相等，就直接进入下次循环
    if b.tophash[i] != top {
        if b.tophash[i] == emptyRest {
            break bucketloop
        }
        continue
    }
	// ...
	// 再比较桶内的 key 和传入的 key，如果相等，再获取目标值的指针
    if t.Key.Equal(key, k) {
		// ...
    }
}
```

{{< callout type="info" >}}
计算在几号桶用的是后 `B` 位，`tophash` 使用的是高 8 位，分别用前后的不同位数，避免了冲突。**这种方式可以避免一个桶内出现大量相同的 `tophash`**，影响读写的性能。
{{< /callout >}}

如果当前桶中没有找到 key，而且**存在溢出桶，那么会接着遍历所有的溢出桶中的数据**。

### 写入

写入 map 和查询 map 的实现原理类似，计算哈希值和存放在哪个桶，然后遍历当前桶和溢出桶的数据：

- 如果当前 key 不存在，则通过偏移量存储到桶中。
- 如果已经存在，则返回 value 的内存地址，然后修改 value。
- 如果桶已满，则会创建新桶或者使用空闲的溢出桶，新桶添加到已有桶的末尾，`noverflow` 计数加 1。将键值对添加到桶中。

{{< callout type="info" >}}
前面说的找到 key 的位置，进行赋值操作，实际上并不准确。看 `mapassign` 函数的原型：

`func mapassign(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer`

`mapassign` 函数返回的指针就是指向的 key 所对应的 value 值位置，有了地址，就很好操作赋值了。在汇编中赋值。
{{< /callout >}}

### 删除

删除 map 中的 `key/value`：

```go
delete(hashmap, key)
```

`delete` 关键字的唯一作用就是将某一个 `key/value` 从哈希表中删除。会被编译器被转换成 `mapdelete` 方法。删除操作先是找到 key 的位置，清空 `key/value`，然后将 `hmap.count - 1`，并且对应的 `tophash` 设置为 `Empty`。

底层执行函数 `mapdelete`：

```go
func mapdelete(t *maptype, h *hmap, key unsafe.Pointer)
```

和上面的定位 key 的逻辑一样，找到对应位置后，对 key 或者 value 进行“清零”操作：

```go
// 对 key 清零
if t.indirectkey {
	*(*unsafe.Pointer)(k) = nil
} else {
	typedmemclr(t.key, k)
}

// 对 value 清零
if t.indirectvalue {
	*(*unsafe.Pointer)(v) = nil
} else {
	typedmemclr(t.elem, v)
}
```

最后，将 `count` 值减 1，将对应位置的 `tophash` 值置成 `Empty`。

### 扩容

随着 map 中写入的 `key/value` 增多，装载因子会越来越大，哈希冲突的概率越来越大，性能会跟着下降。如果大量的 key 都落入到同一个桶中，哈希表会退化成链表，查询的时间复杂度会从 `O(1)` 退化到 `O(n)`。

所以当装载因子大到一定程度之后，哈希表就不得不进行扩容。

#### Go map 在什么时候会触发扩容？

```go
func mapassign(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer {
// src/runtime/map.go mapassign
// If we hit the max load factor or we have too many overflow buckets,
// and we're not already in the middle of growing, start growing.
    if !h.growing() && (overLoadFactor(h.count+1, h.B) || tooManyOverflowBuckets(h.noverflow, h.B)) {
		hashGrow(t, h)
		goto again // Growing the table invalidates everything, so try again
    }
}
```

触发扩容的条件：

1. 装载因子超过阈值 6.5。
2. 溢出桶的数量过多：
   - 当 `B < 15` 时，如果溢出桶的数量超过 `2^B` 则触发扩容。
   - 当 `B >= 15` 时，如果溢出桶的数量超过 `2^15` 则触发扩容，也就说**最多 `2^15` 个溢出桶，避免无限制的增长**。

触发的条件不同扩容的方式也分为两种：

1. 如果这次**扩容是溢出的桶太多导致的，那么这次扩容就是“等量扩容” `sameSizeGrow`**。
2. 另一种就是**装载因子超过阈值导致翻倍扩容**了。

#### 为什么溢出桶过多需要进行扩容？

什么情况下会出现装载因子很小不超过阈值，但是溢出桶过多的情况？

先插入很多元素，导致创建了很多桶，但是未达到阈值，并没有触发扩容。之后再删除元素，降低元素的总量。反复执行前面的步骤，但是又不会触发扩容，就会导致创建了很多溢出桶，但是 map 中的 key 分布的很分散。导致查询和插入的效率很低。还可能导致内存泄漏。

##### 对于条件 2 溢出桶的数量过多

申请的新的 buckets 数量和原有的 buckets 数量是**相等的**，进行的是**等量扩容**。由于 buckets 数量不变，所以原有的数据在几号桶，迁移之后仍然在几号桶。比如原来在 0 号 bucket，到新的地方后，仍然放在 0 号 bucket。

**扩容完成后，溢出桶没有了，key 都集中到了一个 bucket，更为紧凑了，提高了查找的效率**。

##### 对于条件 1 当装载因子超过阈值后

申请的新的 buckets 数量和原有的 buckets 数量的 **2 倍**，也就是 `B+1`。桶的数量改变了，那么 key 的哈希值要重新计算，才能决定它到底落在哪个 bucket。

例如，原来 `B=5`，根据出 key 的哈希值的后 5 位，就能决定它落在哪个 bucket。扩容后的 buckets 数量翻倍，B 变成了 6，因此变成哈希值的后 6 位才能决定 key 落在哪个 bucket。这叫做 `rehash`。

![map-evacuate-bucket-num](https://raw.gitcode.com/shipengqi/illustrations/blobs/13c71d43f205e8f86976dcfe11d26ce99521c83c/map-evacuate-bucket-num.png)

因此，某个 key 在迁移前后 bucket 序号可能会改变，取决于 `rehash` 之后的哈希值倒数第 6 位是 0 还是 1。

**扩容完成后，老 buckets 中的 key 分裂到了 2 个新的 bucket**。

#### 渐进式扩容

扩容需要把原有的 buckets 中的数据迁移到新的 buckets 中。如果一个哈希表当前大小为 1GB，扩容为原来的两倍大小，那就需要对 1GB 的数据**重新计算哈希值（rehash）**，并且从原来的内存空间搬移到新的内存空间，这是非常耗时的操作。

所以 map 的扩容采用的是一种**渐进式**的方式，将迁移的操作穿插在插入操作的过程中，分批完成。

##### 迁移实现

Go map 扩容的实现在 `hashGrow` 函数中，`hashGrow` 只申请新的 buckets，并未参与真正的数据迁移：

```go
func hashGrow(t *maptype, h *hmap) {
	bigger := uint8(1)
	// 溢出桶过多触发的扩容是等量扩容，bigger 设置为 0
	if !overLoadFactor(h.count+1, h.B) {
		bigger = 0
		h.flags |= sameSizeGrow
	}
    // 将原有的 buckets 挂到 oldbuckets 上
	oldbuckets := h.buckets
	// 申请新的 buckets
	newbuckets, nextOverflow := makeBucketArray(t, h.B+bigger, nil)

	flags := h.flags &^ (iterator | oldIterator)
	if h.flags&iterator != 0 {
		flags |= oldIterator
	}
    // 如果是等量扩容，bigger 为 0，B 不变
	h.B += bigger
	h.flags = flags
	// 原有的 buckets 挂到 map 的 oldbuckets 上
	h.oldbuckets = oldbuckets
	// 新申请的 buckets 挂到 buckets 上
	h.buckets = newbuckets
	// 设置迁移进度为 0
	h.nevacuate = 0
	// 溢出桶数量为 0
	h.noverflow = 0
	// ...
}
```

迁移是在插入数据和删除数据时，也就是 `mapassign` 和 `mapdelete` 中进行的：

```go
func mapassign(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer {
    // ... 
again:
	bucket := hash & bucketMask(h.B)
    if h.growing() {
		// 真正的迁移在 growWork 中
        growWork(t, h, bucket)
	}	
	// ...
}

func mapdelete(t *maptype, h *hmap, key unsafe.Pointer) {
    // ... 
	bucket := hash & bucketMask(h.B)
	if h.growing() {
		// 真正的迁移在 growWork 中
		growWork(t, h, bucket)
	}
    // ... 
}

func (h *hmap) growing() bool {
	// oldbuckets 不为空，说明还没有迁移完成
	return h.oldbuckets != nil
}
```

也就是说**数据的迁移过程一般发生在插入或修改、删除 key 的时候**。在扩容完毕后 (预分配内存)，不会马上就进行迁移。而是采取**写时复制**的方式，**当有访问到具体 bucket 时，才会逐渐的将 `oldbucket` 迁移到新 bucket 中**。

`growWork` 函数：

```go
func growWork(t *maptype, h *hmap, bucket uintptr) {
	// 迁移
	evacuate(t, h, bucket&h.oldbucketmask())

	// 还没有迁移完成，额外再迁移一个 bucket，加快迁移进度
	if h.growing() {
		evacuate(t, h, h.nevacuate)
	}
}
```

`evacuate` 函数大致迁移过程如下：

1. 先判断当前 bucket 是不是已经迁移，没有迁移就做迁移操作：

```go
b := (*bmap)(add(h.oldbuckets, oldbucket*uintptr(t.bucketsize)))
newbit := h.noldbuckets()
// 判断旧桶是否已经被迁移了
if !evacuated(b) {
	// 旧桶没有被搬迁
	// do...  // 做转移操作
	// 遍历所有的 bucket，包括 overflow buckets
	// b 是老的 bucket 地址
	for ; b != nil; b = b.overflow(t) {
		// ...
	}
	// ...
}
```

真正的迁移在 `evacuate` 函数中，它会对传入桶中的数据进行再分配。`evacuate` 函数每次只完成一个 bucket 的迁移工作（包括这个 bucket 链接的溢出桶），它会遍历 bucket （包括溢出桶）中得到所有 `key/value` 并迁移。**已迁移的 `key/value` 对应的 `tophash` 会被设置为 `evacuatedEmpty`，表示已经迁移**。

### map 为什么是无序的

map 在扩容后，`key/value` 会进行迁移，在同一个桶中的 key，有些会迁移到别的桶中，有些 key 原地不动，导致遍历 map 就无法保证顺序。

Go 底层的实现简单粗暴，并不是固定地从 0 号 bucket 开始遍历，而是直接生成一个随机数，这个随机数决定从哪里开始遍历，因此**每次 `for range map` 的结果都是不一样的。那是因为它的起始位置根本就不固定**。


### 对 map 元素取地址

Go 无法对 map 的 key 或 value 进行取址。

```go
package main

import "fmt"

func main() {
	m := make(map[string]int)

	fmt.Println(&m["foo"])
}
```

上面的代码不能通过编译：

```
./main.go:8:14: cannot take the address of m["foo"]
```

#### 使用 `unsafe.Pointer` 获取 key 或 value 的地址

可以使用 `unsafe.Pointer` 等获取到了 key 或 value 的地址，也不能长期持有，因为**一旦发生扩容，key 和 value 的位置就会改变，之前保存的地址也就失效了**。

### 比较 map

**比较只能是遍历 `map` 的每个元素**。`map1 == map2` 这种编译是不通过的。

### map 不是并发安全的

在查找、赋值、遍历、删除的过程中都会检测写标志，一旦发现写标志位为 `1`，则直接 panic。

赋值和删除函数在检测完写标志是复位之后，先将写标志位设置位 `1`，才会进行之后的操作。

```go
// 检测写标志
if h.flags&hashWriting == 0 {
	throw("concurrent map writes")
}

// 设置写标志
h.flags |= hashWriting
```





