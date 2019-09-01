---
title: time
---

# time

`time` 提供了一个数据类型 `time.Time`（作为值使用）以及显示和测量时间和日期的功能函数，比如：
- `time.Now()` 获取当前时间。
- `t.Day()`、`t.Minute()` 获取时间的一部分。
- `time.After`、`time.Ticker` 在经过一定时间或周期执行某项任务（事件处理的特例）。
- `time.Sleep（Duration d）` 暂停某个进程（ goroutine），暂停时长为 `d`。
- `Duration` 类型表示两个连续时刻所相差的纳秒数，类型为 `int64`。
- `Location` 类型映射某个时区的时间，UTC 表示通用协调世界时间。

## Format
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