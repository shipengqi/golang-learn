---
title: Go 性能优化
weight: 5
---

# Go 性能优化

## JSON 优化

Go 官方的 `encoding/json` 是通过反射来实现的。性能相对有些慢。 可以使用第三方库来替代标准库：

- [json-iterator/go](https://github.com/json-iterator/go)，完全兼容标准库，性能有很大提升。
- [go-json](https://github.com/goccy/go-json)，完全兼容标准库，性能强于 `json-iterator/go`。
- [sonic](https://github.com/bytedance/sonic)，字节开发的的 JSON 序列化/反序列化库，速度快，但是对硬件有一些要求。

实际开发中可以根据编译标签来选择 `JSON` 库，参考 [component-base/json](https://github.com/shipengqi/component-base/tree/main/json)。

## 内存对齐

## 字符串与字节转换优化，减少内存分配

## 利用 sync.Pool 减少堆分配

## 垃圾回收优化

## 设置 GOMAXPROCS