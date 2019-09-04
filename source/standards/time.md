---
title: time
---

# time

`time` 提供了一个数据类型 `time.Time`（作为值使用）以及显示和测量时间和日期的功能函数，比如：
- `time.Now()` 获取当前时间。
- `t.Day()`、`t.Minute()` 获取时间的一部分。
- `time.After`、`time.Ticker` 在经过一定时间或周期执行某项任务（事件处理的特例）。
- `time.Sleep（Duration d）` 暂停某个进程（ goroutine），暂停时长为 `d`。
- `Duration` 代表两个时间点之间经过的时间，以**纳秒**为单位，类型为 `int64`。
- `Location` 类型映射某个时区的时间，UTC 表示通用协调世界时间。

## 时区
Go 语言使用 `Location` 来表示地区相关的时区，一个 `Location` 可能表示多个时区。

`time` 包提供了 `Location` 的两个实例：
- `Local` 代表当前系统本地时区；
- `UTC` 代表通用协调时间，也就是零时区。`time` 包默认（为显示提供时区）使用 `UTC` 时区。

### Local 是如何做到表示本地时区的？

在初始化 `Local` 时，通过读取 `/etc/localtime` （这是一个符号链接，指向 `/usr/share/zoneinfo` 中某一个时区）可以获取到系统本地时区。

如果设置了**环境变量 `TZ`**，则会优先使用它。

```go
tz, ok := syscall.Getenv("TZ")
switch {
case !ok:
	z, err := loadZoneFile("", "/etc/localtime")
	if err == nil {
		localLoc = *z
		localLoc.name = "Local"
		return
	}
case tz != "" && tz != "UTC":
	if z, err := loadLocation(tz); err == nil {
		localLoc = *z
		return
	}
}
```
### 获得特定时区的实例

函数 `LoadLocation` 可以根据名称获取特定时区的实例：
```go
func LoadLocation(name string) (*Location, error)
```

如果 `name` 是 "" 或 "UTC"，返回 UTC；
如果 `name` 是 "Local"，返回 `Local`；
否则 `name` 应该是 IANA 时区数据库里有记录的地点名（该数据库记录了地点和对应的时区），如 "America/New_York"。

## Time
`Time` 代表一个纳秒精度的时间点。程序中应使用 `time.Time` 类型值来保存和传递时间，而不是指针。
程序中应使用 Time 类型值来保存和传递时间，而不是指针。
一个 Time 类型值可以被多个 go 协程同时使用。时间点可以使用 Before、After 和 Equal 方法进行比较。Sub 方法让两个时间点相减，
生成一个 Duration 类型值（代表时间段）。Add 方法给一个时间点加上一个时间段，生成一个新的 Time 类型时间点。

`Time` 零值代表时间点 `January 1, year 1, 00:00:00.000000000 UTC`。因为本时间点一般不会出现在使用中，`IsZero` 方法提供
了检验时间是否是显式初始化的一个简单途径。

每一个 `Time` 都具有一个地点信息（即对应地点的时区信息），当计算时间的表示格式时，如 Format、Hour 和 Year 等方法，
都会考虑该信息。Local、UTC 和 In 方法返回一个指定时区（但指向同一时间点）的 Time。修改地点 / 时区信息只是会改变其表示；不会修改被表
示的时间点，因此也不会影响其计算。

通过 `==` 比较 Time 时，Location 信息也会参与比较，因此 Time 不应该作为 map 的 key。

```go
type Time struct {
	// sec gives the number of seconds elapsed since
	// January 1, year 1 00:00:00 UTC.
	sec int64

	// nsec specifies a non-negative nanosecond
	// offset within the second named by Seconds.
	// It must be in the range [0, 999999999].
	nsec int32

	// loc specifies the Location that should be used to
	// determine the minute, hour, month, day, and year
	// that correspond to this Time.
	// Only the zero Time has a nil Location.
	// In that case it is interpreted to mean UTC.
	loc *Location
}
```

先看 `time.Now()` 函数。

```go
// Now returns the current local time.
func Now() Time {
	sec, nsec := now()
	return Time{sec + unixToInternal, nsec, Local}
}
```
`now()` 的具体实现在 `runtime` 包中。从 `Time{sec + unixToInternal, nsec, Local}` 可以看出，`Time` 结构的 `sec` 并非 Unix 
时间戳，实际上，加上的 `unixToInternal` 是 `1-1-1` 到 `1970-1-1` 经历的秒数。也就是 `Time` 中的 `sec` 是从 `1-1-1` 算起的秒数，
而不是 Unix 时间戳。

`Time` 的最后一个字段表示地点时区信息。本章后面会专门介绍。

### 常用方法

- `Time.IsZero()` 函数用于判断 Time 表示的时间是否是 0 值。表示 1 年 1 月 1 日。
- `time.Unix(sec, nsec int64`) 通过 Unix 时间戳生成 `time.Time` 实例；
- `time.Time.Unix()` 得到 Unix 时间戳；
- `time.Time.UnixNano()` 得到 Unix 时间戳的纳秒表示；
- `time.Parse` 和 `time.ParseInLocation`
- `time.Time.Format`

```go
t, _ := time.Parse("2006-01-02 15:04:05", "2016-06-13 09:14:00")
fmt.Println(time.Now().Sub(t).Hours())
```
这段代码的结果跟预期的不一样。原因是 `time.Now()` 的时区是 `time.Local`，而 `time.Parse` 解析出来的时区却是 `time.UTC`
（可以通过 `Time.Location()` 函数知道是哪个时区）。在中国，它们相差 8 小时。

所以，一般的，我们应该总是使用 `time.ParseInLocation` 来解析时间，并给第三个参数传递 `time.Local`。

#### 为什么是 2006-01-02 15:04:05

可能你已经注意到：`2006-01-02 15:04:05` 这个字符串了。没错，这是固定写法，类似于其他语言中 `Y-m-d H:i:s` 等。为什么采用这种形式？
又为什么是这个时间点而不是其他时间点？

- 官方说，使用具体的时间，比使用 `Y-m-d H:i:s` 更容易理解和记忆；
- 而选择这个时间点，也是出于好记的考虑，官方的例子：`Mon Jan 2 15:04:05 MST 2006`，很好记：2006 年 1 月 2 日 3 点 4 分 5 秒


#### Round 和 Truncate 方法

比如，有这么个需求：获取当前时间整点的 `Time` 实例。例如，当前时间是 `15:54:23`，需要的是 `15:00:00`。我们可以这么做：

```
t, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:00:00"), time.Local)
fmt.Println(t)
```
实际上，`time` 包给我们提供了专门的方法，功能更强大，性能也更好，这就是 `Round` 和 `Trunate`，它们区别，一个是取最接近的，一个是向下取整。

```go
t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2016-06-13 15:34:39", time.Local)
// 整点（向下取整）
fmt.Println(t.Truncate(1 * time.Hour))
// 整点（最接近）
fmt.Println(t.Round(1 * time.Hour))

// 整分（向下取整）
fmt.Println(t.Truncate(1 * time.Minute))
// 整分（最接近）
fmt.Println(t.Round(1 * time.Minute))

t2, _ := time.ParseInLocation("2006-01-02 15:04:05", t.Format("2006-01-02 15:00:00"), time.Local)
fmt.Println(t2)
```

#### Format
自定义时间格式化字符串，例如： `fmt.Printf("%02d.%02d.%4d\n", t.Day(), t.Month(), t.Year())` 将会输出 `21.07.2011`。

包中的一个预定义函数 `func (t Time) Format(layout string) string` 可以根据一个格式化字符串来将一个时间 t 转换为相应
格式的字符串，你可以使用一些预定义的格式，如：`time.ANSIC` 或 `time.RFC822`。

一般的格式化设计是通过对于一个标准时间的格式化描述来展现的，示例：

```go
fmt.Println(t.Format("02 Jan 2006 15:04")) // 21 Jul 2011 10:31
```

示例：

```go
package main
import (
	"fmt"
	"time"
)

var week time.Duration
func main() {
	t := time.Now()
	fmt.Println(t) // e.g. Wed Dec 21 09:52:14 +0100 RST 2011
	fmt.Printf("%02d.%02d.%4d\n", t.Day(), t.Month(), t.Year())
	// 21.12.2011
	t = time.Now().UTC()
	fmt.Println(t) // Wed Dec 21 08:52:14 +0000 UTC 2011
	fmt.Println(time.Now()) // Wed Dec 21 09:52:14 +0100 RST 2011
	// calculating times:
	week = 60 * 60 * 24 * 7 * 1e9 // must be in nanosec
	week_from_now := t.Add(time.Duration(week))
	fmt.Println(week_from_now) // Wed Dec 28 08:52:14 +0000 UTC 2011
	// formatting times:
	fmt.Println(t.Format(time.RFC822)) // 21 Dec 11 0852 UTC
	fmt.Println(t.Format(time.ANSIC)) // Wed Dec 21 08:56:34 2011
	fmt.Println(t.Format("02 Jan 2006 15:04")) // 21 Dec 2011 08:52
	s := t.Format("20060102")
	fmt.Println(t, "=>", s)
	// Wed Dec 21 08:52:14 +0000 UTC 2011 => 20111221
}
```

示例，计算函数执行时间：
```go
start := time.Now() // 起始时间
longCalculation()
end := time.Now() // 结束时间
delta := end.Sub(start) // 消耗时间
fmt.Printf("longCalculation took this amount of time: %s\n", delta)
```

## 定时器
`time` 包有两种定时器：
- `Timer`（到达指定时间触发且只触发一次）
- `Ticker`（间隔特定时间触发）。

### Timer
```go
type Timer struct {
	C <-chan Time	 // The channel on which the time is delivered.
	r runtimeTimer
}
```

`Timer` 的实例必须通过 `NewTimer` 或 `AfterFunc` 获得。
**当 `Timer` 到期时，当时的时间会被发送给 `C` (channel)，除非 `Timer` 是被 `AfterFunc` 函数创建的**。

`runtimeTimer` 定义在 `sleep.go` 文件中，必须和 `runtime` 包中 `time.go` 文件中的 `timer` 保持一致：
```go
type timer struct {
	i int // heap index

	// Timer wakes up at when, and then at when+period, ... (period > 0 only)
	// each time calling f(now, arg) in the timer goroutine, so f must be
	// a well-behaved function and not block.
	when   int64
	period int64
	f      func(interface{}, uintptr)
	arg    interface{}
	seq    uintptr
}
```

通过 `NewTimer()` 来看这些字段：

```go
// NewTimer creates a new Timer that will send
// the current time on its channel after at least duration d.
func NewTimer(d Duration) *Timer {
	c := make(chan Time, 1)
	t := &Timer{
		C: c,
		r: runtimeTimer{
			when: when(d),
			f:    sendTime,
			arg:  c,
		},
	}
	startTimer(&t.r)
	return t
}
```

在 `when` 表示的时间到时，会往 `Timer.C` 中发送当前时间。`when` 表示的时间是纳秒时间，正常通过 `runtimeNano() + int64(d)` 赋值。

`f` 参数的值是 `sendTime`，定时器时间到时，会调用 `f`，并将 `arg` 和 `seq` 传给 `f`。

因为 `Timer` 是一次性的，所以 `period` 保留默认值 0。


### 常用方法
#### time.After
`time.After` 模拟超时：
```go
c := make(chan int)

go func() {
	// time.Sleep(1 * time.Second)
	time.Sleep(3 * time.Second)
	<-c
}()

select {
case c <- 1:
	fmt.Println("channel...")
case <-time.After(2 * time.Second):
	close(c)
	fmt.Println("timeout...")
}
```

#### time.Stop 和 time.Reset
```go
start := time.Now()
timer := time.AfterFunc(2*time.Second, func() {
	fmt.Println("after func callback, elaspe:", time.Now().Sub(start))
})

time.Sleep(1 * time.Second)
// time.Sleep(3*time.Second)
// Reset 在 Timer 还未触发时返回 true；触发了或 Stop 了，返回 false
if timer.Reset(3 * time.Second) {
	fmt.Println("timer has not trigger!")
} else {
	fmt.Println("timer had expired or stop!")
}

time.Sleep(10 * time.Second)

// output:
// timer has not trigger!
// after func callback, elaspe: 4.00026461s
```

`timer.Stop()` 不会关闭 `Timer.C` 这个 channel，可以使用 `timer.Reset(0)` 代替 `timer.Stop()` 来停止定时器。

#### Sleep

`Sleep` 的是通过 `Timer` 实现的（查看 `runtime/time.go` 文件）。用于暂停当前 `goroutine`。

## Ticker

`Ticker` 和 `Timer` 类似，区别是：`Ticker` 中的 `runtimeTimer` 字段的 `period` 字段会赋值为 `NewTicker(d Duration)` 中的 `d`，
表示每间隔 `d` 纳秒，定时器就会触发一次。

除非程序终止前定时器一直需要触发，否则，不需要时应该调用 `Ticker.Stop` 来释放相关资源。

如果程序终止前需要定时器一直触发，可以使用更简单方便的 `time.Tick` 函数，因为 `Ticker` 实例隐藏起来了，因此，该函数启动的定时器无法停止。
