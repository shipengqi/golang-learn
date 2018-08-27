<h1 align="center">
  Go 语言学习笔记
</h1>

<p align="center">
  <a href="https://github.com/golang/go">
    <img alt="Go Logo" src="https://upload.wikimedia.org/wikipedia/commons/2/23/Go_Logo_Aqua.svg" width="20%" height="">
  </a>
</p>



## Go语言的主要特性
- 编译型语言，静态类型语言，但是语法非常简洁
- 自动回收垃圾
- Go 有自己原生的并发编程模型
- 函数式编程
- 丰富的标准库

## 语法
- 每个文件头部使用`package`声明，表示该文件属于哪个包，e.g. `package main`，紧跟着使用`import`导入包
，注意如果导入了不需要的包，编译会失败。
- 入口函数`mian`没有参数，必须在`main`包中，定义了一个独立可执行的程序。
- 可以直接运行（`go run *.go`，其实也是`run`命令进行编译运行），也可以编译后运行。
- 使用`var`声明变量，函数内部可以省略`var`。
- 函数可以返回多个值，函数是第一类型，可以作为参数或返回值。
- 控制语句只有三种`if`，`for`，`switch`。
- `os.Args`获取命令行参数。

### 类型
