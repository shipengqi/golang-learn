---
title: Go 汇编工具和 dlv 使用
weight: 12
---

## 生成汇编

生成汇编代码有两种方式：

```bash
# 即将源代码编译成 .o 目标文件，并输出汇编代码。
go tool compile -S main.go

# 或者
# 反汇编，即从可执行文件反编译成汇编代码，所以要先用 go build 命令编译出可执行文件。
go build main.go && go tool objdump ./main
```

使用时可以加上 `grep` 的来加一个过滤条件：

```bash
go tool compile -S main.go | grep "main.go:4"

# 或者
go build main.go && go tool objdump ./main | grep "main.go:4"
```

## dlv 使用

```go
package main

func main() {
	var m map[int]int
	m[1] = 1
}
```

使用 dlv 来调试向一个 `nil` 的 `map` 插入新元素，为什么会报 `panic`。

```bash
[root@shcCDFrh75vm7 ~]# go build main.go # 编译生成可执行文件
[root@shcCDFrh75vm7 ~]# dlv exec ./main # 使用 dlv 进入调试状态
Type 'help' for list of commands.
(dlv) 
```

使用 `b` 命令可以打断点，有三种方式：

- `b + 地址`
- `b + 代码行数`
- `b + 函数名`

现在在对 `map` 赋值的地方加个断点。先找到代码位置，`m[1] = 1` 这行在第 5 行：

```bash
(dlv) b main.go:5
Breakpoint 1 set at 0x45d9ee for main.main() ./main.go:5
(dlv) 
```

现在可以使用 `c` 命令来继续运行程序，当程序运行到断点处时会停下来：

```bash
Breakpoint 1 set at 0x45d9ee for main.main() ./main.go:5
(dlv) c
> main.main() ./main.go:5 (hits goroutine(1):1 total:1) (PC: 0x45d9ee)
Warning: debugging optimized function
     1: package main
     2:
     3: func main() {
     4:         var m map[int]int
=>   5:         m[1] = 1
     6: }
(dlv) 
```

使用 `disass` 命令，可以看到汇编指令：

```bash
(dlv) disass
TEXT main.main(SB) /root/main.go
        main.go:3       0x45d9e0        493b6610        cmp rsp, qword ptr [r14+0x10]
        main.go:3       0x45d9e4        762c            jbe 0x45da12
        main.go:3       0x45d9e6        55              push rbp
        main.go:3       0x45d9e7        4889e5          mov rbp, rsp
        main.go:3       0x45d9ea        4883ec18        sub rsp, 0x18
=>      main.go:5       0x45d9ee*       488d050b730000  lea rax, ptr [rip+0x730b]
        main.go:5       0x45d9f5        31db            xor ebx, ebx
        main.go:5       0x45d9f7        b901000000      mov ecx, 0x1
        main.go:5       0x45d9fc        0f1f4000        nop dword ptr [rax], eax
        main.go:5       0x45da00        e85b0ffbff      call $runtime.mapassign_fast64
        main.go:5       0x45da05        48c70001000000  mov qword ptr [rax], 0x1
        main.go:6       0x45da0c        4883c418        add rsp, 0x18
        main.go:6       0x45da10        5d              pop rbp
        main.go:6       0x45da11        c3              ret
        main.go:3       0x45da12        e869cdffff      call $runtime.morestack_noctxt
        main.go:3       0x45da17        ebc7            jmp $main.main
(dlv) 
```

这时使用 `si` 命令（`si` 全称是 Step Instruction，用于在汇编级别单步执行一条机器指令），执行单条指令，多次执行 `si`，就会执行到 `map` 赋值函数 `mapassign_fast64`：

```bash
# ...
(dlv) si
> main.main() ./main.go:5 (PC: 0x45da00)
Warning: debugging optimized function
        main.go:3       0x45d9ea        4883ec18        sub rsp, 0x18
        main.go:5       0x45d9ee*       488d050b730000  lea rax, ptr [rip+0x730b]
        main.go:5       0x45d9f5        31db            xor ebx, ebx
        main.go:5       0x45d9f7        b901000000      mov ecx, 0x1
        main.go:5       0x45d9fc        0f1f4000        nop dword ptr [rax], eax
=>      main.go:5       0x45da00        e85b0ffbff      call $runtime.mapassign_fast64
        main.go:5       0x45da05        48c70001000000  mov qword ptr [rax], 0x1
        main.go:6       0x45da0c        4883c418        add rsp, 0x18
        main.go:6       0x45da10        5d              pop rbp
        main.go:6       0x45da11        c3              ret
        main.go:3       0x45da12        e869cdffff      call $runtime.morestack_noctxt
(dlv) si
> runtime.mapassign_fast64() /usr/local/go/src/runtime/map_fast64.go:93 (PC: 0x40e960)
Warning: debugging optimized function
TEXT runtime.mapassign_fast64(SB) /usr/local/go/src/runtime/map_fast64.go
=>      map_fast64.go:93        0x40e960        493b6610        cmp rsp, qword ptr [r14+0x10]
        map_fast64.go:93        0x40e964        0f86e0020000    jbe 0x40ec4a
        map_fast64.go:93        0x40e96a        55              push rbp
        map_fast64.go:93        0x40e96b        4889e5          mov rbp, rsp
        map_fast64.go:93        0x40e96e        4883ec30        sub rsp, 0x30
        map_fast64.go:93        0x40e972        48894c2450      mov qword ptr [rsp+0x50], rcx
(dlv) 
```

这时再用单步命令 `s`（`s` 全称是 Step（或 Step Into），用于在源码级别单步执行程序，并进入被调用的函数内部），就会进入判断 `h` 的值为 `nil` 的分支，然后执行 `panic` 函数：

```bash
(dlv) si
> runtime.mapassign_fast64() /usr/local/go/src/runtime/map_fast64.go:93 (PC: 0x40e960)
Warning: debugging optimized function
TEXT runtime.mapassign_fast64(SB) /usr/local/go/src/runtime/map_fast64.go
=>      map_fast64.go:93        0x40e960        493b6610        cmp rsp, qword ptr [r14+0x10]
        map_fast64.go:93        0x40e964        0f86e0020000    jbe 0x40ec4a
        map_fast64.go:93        0x40e96a        55              push rbp
        map_fast64.go:93        0x40e96b        4889e5          mov rbp, rsp
        map_fast64.go:93        0x40e96e        4883ec30        sub rsp, 0x30
        map_fast64.go:93        0x40e972        48894c2450      mov qword ptr [rsp+0x50], rcx
(dlv) s
> runtime.mapassign_fast64() /usr/local/go/src/runtime/map_fast64.go:94 (PC: 0x40e980)
Warning: debugging optimized function
    89:         }
    90:         return unsafe.Pointer(&zeroVal[0]), false
    91: }
    92:
    93: func mapassign_fast64(t *maptype, h *hmap, key uint64) unsafe.Pointer {
=>  94:         if h == nil {
    95:                 panic(plainError("assignment to entry in nil map"))
    96:         }
    97:         if raceenabled {
    98:                 callerpc := getcallerpc()
    99:                 racewritepc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(mapassign_fast64))
(dlv) s
> runtime.mapassign_fast64() /usr/local/go/src/runtime/map_fast64.go:95 (PC: 0x40ec36)
Warning: debugging optimized function
    90:         return unsafe.Pointer(&zeroVal[0]), false
    91: }
    92:
    93: func mapassign_fast64(t *maptype, h *hmap, key uint64) unsafe.Pointer {
    94:         if h == nil {
=>  95:                 panic(plainError("assignment to entry in nil map"))
    96:         }
    97:         if raceenabled {
    98:                 callerpc := getcallerpc()
    99:                 racewritepc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(mapassign_fast64))
   100:         }
(dlv) 
```

至此，向 `nil` 的 `map` 赋值时，产生 panic 的代码就被找到了。

可以使用 `bt` 命令看到调用栈：

```bash
(dlv) bt
0  0x000000000042fa04 in runtime.fatalpanic
   at /usr/local/go/src/runtime/panic.go:1217
1  0x000000000042e9f8 in runtime.gopanic
   at /usr/local/go/src/runtime/panic.go:779
2  0x000000000040ec49 in runtime.mapassign_fast64
   at /usr/local/go/src/runtime/map_fast64.go:95
3  0x000000000045da05 in main.main
   at ./main.go:5
4  0x0000000000431f7d in runtime.main
   at /usr/local/go/src/runtime/proc.go:271
5  0x000000000045aa41 in runtime.goexit
   at /usr/local/go/src/runtime/asm_amd64.s:1695
(dlv)
```

之后可以使用 `frame {num}` 命令可以跳转到相应位置。例如 `frame 1` 就是跳转到 `panic.go:779` 的位置。

### 更多命令

| 命令 | 全称 | 作用 | 使用场景 |
| --- | --- | --- | --- |
| `n` | Next | 执行当前行，但不进入函数 | 快速跳过已知函数 |
| `so` | Step Out | 跳出当前函数 | 快速返回调用处 |
| `ni` | Next Instruction | 类似 `si`，但不会进入函数调用 |  |
| `r` | restart | 重新启动调试会话 |  |
| `clear` | Clear | 删除断点 | 如 `clear 1` 删除 `ID=1` 的断点 |
| `clearall` | Clear All | 删除所有断点 |  |
| `p` | Print | 打印变量的值 | 调试代码时，需要查看变量的值，如 `p x` 或 `p x + y` |
| `set` | Set | 修改变量值 | 如 `set x = 10` |
| `disass` | `disassemble` | 反汇编当前函数 | 查看汇编代码 |
| `gr` | goroutine | 显示当前 goroutine 信息 |  |
| `grs` |  goroutines | 列出所有 goroutines | |
| `switch` | `switch` | 切换 goroutine | 如 `switch 2` |
