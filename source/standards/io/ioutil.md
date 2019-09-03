---
title: ioutil
---

# ioutil

ioutil 提供了一些常用、方便的 IO 操作函数。

## NopCloser 函数

有时候我们需要传递一个 `io.ReadCloser` 的实例，而我们现在有一个 `io.Reader` 的实例，比如：`strings.Reader` ，这个时候 `NopCloser` 
就派上用场了。它包装一个 `io.Reader`，返回一个 `io.ReadCloser` ，而相应的 `Close` 方法啥也不做，只是返回 `nil`。

比如，在标准库 `net/http` 包中的 `NewRequest`，接收一个 `io.Reader` 的 `body`，而实际上，`Request` 的 `Body `的类
型是 `io.ReadCloser`，因此，代码内部进行了判断，如果传递的 `io.Reader` 也实现了 `io.ReadCloser` 接口，则转换，否则通
过 `ioutil.NopCloser` 包装转换一下。相关代码如下：
```go
rc, ok := body.(io.ReadCloser)
if !ok && body != nil {
    rc = ioutil.NopCloser(body)
}
```

## ReadAll 函数

很多时候，我们需要一次性读取 `io.Reader` 中的数据，考虑到读取所有数据的需求比较多，Go 提供了 `ReadAll` 这个函数，用来从 `io.Reader` 中
一次读取所有数据。
```go
func ReadAll(r io.Reader) ([]byte, error)
```
该函数成功调用后会返回 `err == nil` 而不是 `err == EOF`。

## ReadDir 函数



在 Go 中如何输出目录下的所有文件呢？

在 `ioutil` 中提供了一个方便的函数：`ReadDir`，它读取目录并返回排好序的文件和子目录名（ `[]os.FileInfo` ）。

```go
func main() {
	dir := os.Args[1]
	listAll(dir,0)
}

func listAll(path string, curHier int){
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil{fmt.Println(err); return}

	for _, info := range fileInfos{
		if info.IsDir(){
			for tmpHier := curHier; tmpHier > 0; tmpHier--{
				fmt.Printf("|\t")
			}
			fmt.Println(info.Name(),"\\")
			listAll(path + "/" + info.Name(),curHier + 1)
		}else{
			for tmpHier := curHier; tmpHier > 0; tmpHier--{
				fmt.Printf("|\t")
			}
			fmt.Println(info.Name())
		}
	}
}
```

## ReadFile 和 WriteFile 函数

`ReadFile` 读取整个文件的内容，`ReadFile` 会先判断文件的大小，给 `bytes.Buffer` 一个预定义容量，避免额外分配内存。

```go
func ReadFile(filename string) ([]byte, error)
```

> `ReadFile` 从 `filename` 指定的文件中读取数据并返回文件的内容。成功的调用返回的 `err` 为 `nil` 而非 `EOF`。因为本函数定义为
读取整个文件，它不会将读取返回的 `EOF` 视为应报告的错误。

**WriteFile**：
```go
func WriteFile(filename string, data []byte, perm os.FileMode) error
```


> `WriteFile` 将 `data` 写入 `filename` 文件中，当文件不存在时会根据 `perm` 指定的权限进行创建一个,文件存在时会先清空文件内容。
对于 `perm` 参数，我们一般可以指定为：`0666`。

## TempDir 和 TempFile 函数

操作系统中一般都会提供临时目录，比如 linux 下的 `/tmp` 目录（通过 `os.TempDir()` 可以获取到)。有时候，我们自己需要创建临时目录，
比如 Go 工具链源码中（`src/cmd/go/build.go`），通过 `TempDir` 创建一个临时目录，用于存放编译过程的临时文件：
```go
b.work, err = ioutil.TempDir("", "go-build")
```
第一个参数如果为空，表明在系统默认的临时目录（ `os.TempDir` ）中创建临时目录；第二个参数指定临时目录名的前缀，该函数返回临时目录的路径。

相应的，`TempFile` 用于创建临时文件。如 `gofmt` 命令的源码中创建临时文件：
```go
f1, err := ioutil.TempFile("", "gofmt")
```
参数和 `ioutil.TempDir` 参数含义类似。

这里需要**注意**：创建者创建的临时文件和临时目录要负责删除这些临时目录和文件。如删除临时文件：
```go
defer func() {
    f.Close()
    os.Remove(f.Name())
}()
```
