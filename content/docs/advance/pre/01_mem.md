---
title: 计算机内存和寄存器
weight: 1
---

## 计算机为什么需要内存

计算机是运行程序的载体，进程由可执行代码被执行后产生。那么计算机在运行程序的过程中为什么需要内存呢？

### 代码的本质

简单来看代码主要包含两部分：

- 指令部分：中央处理器 CPU 可执行的指令
- 数据部分：常量等

代码包含了指令，代码被转化为可执行二进制文件，被执行后加载到内存中，中央处理器 CPU 通过内存获取指令：

![code-in-memory](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/code-in-memory.png)

### 程序的运行过程

可执行代码文件被执行之后，代码中的待执行指令被加载到了内存当中。

CPU 执行指令可以简单的分为三步：

1. 取指：CPU 控制单元从内存中获取指令
2. 译指：CPU 控制单元解析从内存中获取指令
3. 执行：CPU 运算单元负责执行具体的指令操作

### 内存的作用

- 暂存二进制可执行代码文件中的指令、预置数据(常量)等
- 暂存指令执行过程中的中间数据
- 等等

### 为什么需要栈内存

进程在运行过程中会产生很多临时数据，要关注两个问题：

- 内存的分配
- 内存的回收

最简单、高效地分配和回收方式**线性分配**。线性分配分配的内存是一段连续的区域。栈内存就是使用线性分配的方式进行内存管理的。

**栈内存**的简易管理过程：

1. 栈内存分配逻辑：current - alloc

![mem-alloc](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/mem-alloc.png)

2. 栈内存释放逻辑：current + alloc

![mem-release](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/mem-release.png)

通过利用**栈内存**，**CPU 在执行指令过程中可以高效的存储临时变量**。其次：

- 栈内存的分配过程：类似**栈**的**入栈**过程。
- 栈内存的释放过程：类似**栈**的**出栈**过程。

![mem-stack](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/mem-stack.png)

**栈内存的分配和释放类似一个栈结构，所以叫做栈内存**。

### 为什么需要堆内存

如果函数 A 内的变量 `too` 是个指针且被函数外的代码依赖，如果 `too` 变量指向的内存被回收了，那么这个指针就成了野指针不安全。

**什么是野指针？**

**野指针就是指向一个已被释放的内存地址的指针**。野指针指向的内存地址**可能是其他变量的内存地址，也可能是无效的内存地址**。野指针指向的内存地址可能会被其他程序占用，也可能会被操作系统回收。如果程序继续访问野指针指向的内存地址，就会导致程序崩溃或者数据错误。

怎么解决这个问题？

这就是**堆内存**存在的意义，Go 语言会在代码编译期间通过**逃逸分析**把**分配在栈上的变量分配到堆上去**。**堆内存再通过垃圾回收器回收**。

![mem-escape](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/mem-escape.png)


### 虚拟内存

程序实际操作的都是虚拟内存，最终由 CPU 通过内存管理单元 MMU (Memory Manage Unit) 把虚拟内存的地址转化为实际的物理内存地址。

虚拟内存是一种内存管理技术，它：

- 为**每个进程提供独立的、连续的地址空间**
- 防止了进程直接对物理内存的操作 (如果进程可以直接操作物理内存，那么存在某个进程篡改其他进程数据的可能)
- 提升物理内存的利用率，当进程真正要使用物理内存时再分配
- 通过分页机制将虚拟地址映射到物理地址
- 虚拟内存和物理内存是通过 MMU (管理单元内存 Memory Management Unit) 映射的

对于 **Go，不管是栈内存还是堆内存，都是对虚拟内存的操作**。

### 内存管理组成部分

内存管理一般包含三个不同的组件，分别是用户程序（Mutator）、分配器（Allocator）和收集器（Collector）。

- **用户程序**：用户程序是指正在运行的程序，它可以是一个操作系统、一个数据库系统、一个 Web 浏览器等等。用户程序需要访问内存来存储数据和执行指令。
- **分配器**：分配器会负责从堆中初始化相应的内存区域（栈内存是由编译器自动分配回收的）。
- **收集器**：负责回收堆中不再使用的对象和内存空间。

## 寄存器

寄存器是 CPU 内部的存储单元，用于存放从**内存读取而来的数据（包括指令）和 CPU 运算的中间结果**。

之所以要使用寄存器来临时存放数据而不是直接操作内存，一是因为 CPU 的工作原理决定了**有些操作运算只能在 CPU 内部进行**，二是因为 **CPU 读写寄存器的速度比读写内存的速度快得多**。

{{< callout type="info" >}}
寄存器，是 CPU 内置的容量小、但速度极快的内存。而程序计数器，则是用来存储 CPU 正在执行的指令位置、或者即将执行的下一条指令位置。它们都是 CPU 在运行任何任务前，必须的依赖环境，因此也被叫做 **CPU 上下文**。
{{< /callout >}}

CPU 厂商为每个寄存器都取了一个名字，比如 AMD64 CPU 中的 rax, rbx, rcx, rdx 等等，这样就可以很方便的在汇编代码中使用寄存器的名字来进行编程。

示例，Go 代码：

```go
c = a + b
```

在 AMD64 Linux 平台下，使用 go 编译器编译它可得到如下 `AT&T` 格式的汇编代码：

```asm
mov (%rsp),%rdx     // 把变量 a 的值从内存中读取到寄存器 rdx 中
mov 0x8(%rsp),%rax  // 把变量 b 的值从内存中读取到寄存器 rax 中
add %rdx,%rax       // 把寄存器 rdx 和 rax 中的值相加，并把结果放回 rax 寄存器中
mov %rax,0x10(%rsp) // 把寄存器 rax 中的值写回变量 c 所在的内存
```

上面的一行 go 语言代码被编译成了 4 条汇编指令，指令中出现的 rax，rdx 和 rsp 都是寄存器的名字（`AT&T`格式的汇编代码中所有寄存器名字前面都有一个 `%`符号）。

**汇编代码**其实比较简单，它所做的工作不外乎就是**把数据在内存和寄存器中搬来搬去或做一些基础的数学和逻辑运算**。

不同体系结构的CPU，其内部寄存器的数量、种类以及名称可能大不相同。

以 AMD64 架构的 CPU 为例，常用的三类寄存器：

1. **通用寄存器**：用来存放一般性的数据，用途没有做特殊规定，程序员和编译器可以自定义其用途。

  16 个通用寄存器，分别是：

  - rax, rbx, rcx, rdx
  - rsp（栈顶寄存器）, rbp（栈基址寄存器）
  - rsi, rdi
  - r8, r9, r10, r11, r12, r13, r14, r15

2. **程序计数寄存器（rip）**：也叫做 PC 寄存器或者 IP 寄存器。**存放下一条即将执行的指令的地址**。
3. **段寄存器**：fs 和 gs 寄存器。**一般用它来实现线程本地存储（TLS）**

除了 fs 和 gs 段寄存器是 16 位的，其它都是 64 位的，也就是 8 个字节，其中的 16 个通用寄存器还可以作为 `32/16/8` 位寄存器使用。只是使用时需要换一个名字，比如可以用 eax 这个名字来表示一个 32 位的寄存器，它使用的是 rax 寄存器的低 32 位。

### 程序计数寄存器（rip）

rip 寄存器里面存放的是 CPU 即将执行的下一条指令在内存中的地址。例如下面的汇编代码：

```asm
0x0000000000400770:    add %rdx,%rax
0x0000000000400773:    mov $0x0,%ecx
```

假设当前 CPU 正在执行第一条指令，这条指令在内存中的地址是 `0x0000000000400770`，紧接它后面的下一条指令的地址是 `0x0000000000400773`，所以此时 rip 寄存器里面存放的值是 `0x0000000000400773`。

**rip 寄存器的值是 CPU 自动控制的，CPU 也提供了几条可以间接修改 rip 寄存器的指令**。

### 栈顶寄存器（rsp）和栈基址寄存器（rbp）

rsp 寄存器一般用来存放函数调用栈的栈顶地址，而 rbp 寄存器通常用来存放函数的栈帧起始地址，编译器一般**使用这两个寄存器加一定偏移的方式来访问函数局部变量或函数参数**。

#### 函数栈帧

1. 每个未运行完的函数都有对应的栈帧。
2. 栈帧保存了函数的返回地址和局部变量。

#### 栈帧创建于销毁过程

假设代码：

```cpp
#include<stdio.h>
 
int add(int a, int b)
{
	int c = 0;
	c = a + b;
	return c;
}
 
int main()
{
	int a = 1;
	int b = 1;
	int sum;
	sum = add(a, b);
	return 0;
}
```

调用 `add(a, b)` 之前，栈的情况如下：

<img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/rbp-rsp.png" alt="rbp-rsp" width="280px">

1. 函数调用涉及到传参，因此**在调用函数之前，需要先将传入的参数保存**，以方便函数的调用，因此需要将 `add` 函数的 `a=1`，`b=2` 入栈保存。

<img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/add-before.png" alt="add-before" width="280px">

2. 函数调用前，要创建新的栈帧，rsp 和 rbp 都要改变，为了**函数调用结束后，栈顶恢复到调用前的位置，因此需要先将 rbp 寄存器的值保存到栈中**。

<img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/add-save-rbp.png" alt="add-save-rbp" width="280px">

3. 创建创建所需调用函数的栈帧，使 rbp 指向当前 rsp 的位置，根据 `add` 函数的参数个数，创建合适的栈帧大小。

<img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/add-change-rbp.png" alt="add-change-rbp" width="280px">

<img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/add-create-stack.png" alt="add-create-stack" width="280px">

4. 保存局部变量。将 `add` 函数中创建的变量 `int c = 0` 放入刚刚开辟的栈帧空间中。

<img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/add-save-var.png" alt="add-save-var" width="280px">

5. 参数运算。根据形参与局部变量，进行对应的运算，这里执行 `c = a + b`, 得到 `c = 2`,放入刚才 `c` 对应的位置。

<img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/add-calc.png" alt="add-calc" width="280px">

{{< callout type="info" >}}
假设上一栈帧的 rbp 和参数都是 4 字节，那么参数寻址就是 `rbp + 偏移量`：

- `rbp+4`：上一栈帧的 rbp
- `rbp+8`：参数 a
- `rbp+12`：参数 b
{{< /callout >}}


6. 函数返回。`add` 函数执行完成，需要将 `add` 创建的函数栈销毁，以返回到 `main` 函数中继续执行。在销毁 `add` 函数栈之前，要先将 `add` 函数的返回值 `c = 2` 保存起来，存储到 rax 寄存器中。

<img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/add-save-result.png" alt="add-save-result" width="380px">

7. 销毁 `add` 函数栈。

<img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/add-destroy-stack.png" alt="add-destroy-stack" width="380px">

8. rbp 寄存器拿到之前存储的上一栈帧栈底的值，回到相应的位置，于此同时，栈空间内存储的 rbp 的值没有用了，也将被销毁。

<img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/add-destroy-rbp.png" alt="add-destroy-rbp" width="380px">

9. 销毁形参。

<img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/add-destroy-arg.png" alt="add-destroy-arg" width="380px">

10. `main` 函数拿到返回值。`main` 函数是一个函数，它有自己的栈帧。因此所谓的前一栈帧实际上就是调用 `add` 函数的 `main` 函数的栈帧。因此要让 `main` 函数拿到返回值，只需要把 rax 寄存器中的值放入 `main` 栈帧中 `sum` 对应的位置就行。

<img src="https://raw.gitcode.com/shipengqi/illustrations/files/main/go/main-stack.png" alt="main-stack" width="380px">

绿色部分就是 `main` 函数的栈帧（这里的 `a=1,b=1` 是 `main` 栈帧的局部变量）。至此栈帧的创建与销毁结束，函数调用完成。


### Go 汇编寄存器

Go 汇编格式跟前面讨论过的 AT&T 汇编基本上差不多。

Go 汇编语言中使用的寄存器的名字与 AMD64 汇编中的寄存器的名字不一样，它们之间的对应关系如下：

| go 寄存器 | amd64 寄存器 |
| --------- | ------------ |
| AX        | rax         |
| BX        | rbx         |
| CX        | rcx         |
| DX        | rdx         |
| SI        | rsi         |
| DI        | rdi         |
| SP        | rsp         |
| BP        | rbp         |
| PC        | rip         |
| `R8 ~ R15`        | `r8 ~ r15`   |

Go 汇编还引入了几个没有任何硬件寄存器与之对应的**虚拟寄存器**。这些寄存器一般用来存放内存地址。