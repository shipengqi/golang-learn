# 函数

### 声明函数
`func`关键字声明函数：
```go
func 函数名(形式参数列表) (返回值列表) {
    函数体
}
```
如果函数返回一个无名变量或者没有返回值，返回值列表的括号可以省略。如果一个函数声明没有返回值列表，那么这个
函数不会返回任何值。

```go
// 两个int 类型参数 返回一个 int 类型的值
func max(num1, num2 int) int {
   /* 定义局部变量 */
   var result int

   if (num1 > num2) {
      result = num1
   } else {
      result = num2
   }
   return result 
}

// 返回多个类型的值
func swap(x int, y string) (string, int) {
   return y, x
}

// 有名返回值
func Size(rect image.Rectangle) (width, height int, err error)
```

在函数体中，函数的形参作为局部变量，被初始化为调用者提供的值（函数调用必须按照声明顺序为所有参数提供实参）。函数的形参和有名返回值（也就是对返回值命名）作为函数最外层的局部变量，
被存储在相同的词法块中。

**Go 语言使用的是值传递，当我们传一个参数值到被调用函数里面时，实际上是传了这个值的一份copy，当在被调用函数中修改参数值的时候，调用函数中相应实参不会发生任何变化，
因为数值变化只作用在copy上。但是如果是引用传递，在调用函数时将实际参数的地址传递到函数中，那么在函数中对参数所进行的修改，将影响到实际参数。**

注意，如果实参是`slice`、`map`、`function`、`channel`等类型，实参可能会由于函数的间接引用被修改。

没有函数体的函数声明，这表示该函数不是以 Go 实现的。这样的声明定义了函数标识符。

### 函数作为值
Go 函数被看作第一类值：函数像其他值一样，拥有类型，可以被赋值给其他变量，传递给函数，从函数返回。

```go
func main(){
	/* 声明函数变量 */
	getSquareRoot := func(x float64) float64 {
		return math.Sqrt(x)
	}

	/* 使用函数 */
	fmt.Println(getSquareRoot(9)) // 3
}
```

### 函数作为参数
声明一个名叫operate的函数类型，它有两个参数和一个结果，都是int类型的。
```go
type operate func(x, y int) int
```

编写calculate函数的签名部分。这个函数除了需要两个int类型的参数之外，还应该有一个operate类型的参数。
```go
func calculate(x int, y int, op operate) (int, error) {
    if op == nil {
        return 0, errors.New("invalid operation")
    }
    return op(x, y), nil
}
```

### 闭包
Go 语言支持匿名函数，可作为闭包。
```go
// 返回一个函数
func getSequence() func() int { // func() 是没有参数也没有返回值的函数类型
	 i:=0
	 // 闭包
   return func() int {
      i+=1
     return i  
   }
}
```

### 递归
递归，就是在运行的过程中调用自己。

### 错误
Go 中，对于大部分函数而言，永远无法确保能否成功运行（有一部分函数总是能成功的运行。比如`strings.Contains`和`strconv.FormatBool`）。
通常 Go 函数的最后一个返回值用来传递错误信息。如果导致失败的原因只有一个，返回值可以是一个布尔值，通常被命名为`ok`。否则应该返回一个`error`类型。

#### error 类型
`error`类型是内置的接口类型。`error`类型可能是`nil`或者`non-nil`，`nil`表示成功。

#### 错误处理
当函数调用返回错误时，最常用的处理方式是传播错误，如。
```go
resp, err := http.Get(url)
if err != nil{ // 将这个HTTP错误返回给调用者
    return nil, err
}


doc, err := html.Parse(resp.Body)
resp.Body.Close()
if err != nil {
	// fmt.Errorf函数使用fmt.Sprintf格式化错误信息并返回
	// 使用该函数前缀添加额外的上下文信息到原始错误信息。
  return nil, fmt.Errorf("parsing %s as HTML: %v", url,err)
}
```
由于错误信息经常是以链式组合在一起的，所以错误信息中应避免大写和换行符。

编写错误信息时，我们要确保错误信息对问题细节的描述是详尽的。尤其是要注意错误信息表达的一致性，即相同的函数或同包内
的同一组函数返回的错误在构成和处理方式上是相似的。

根据不同的场景，我们可能要对错误做些特殊处理，比如错误重试机制，或者打印错误日志，或者直接忽略错误。

#### 文件结尾错误
`io`包在任何由文件结束引起的读取失败都返回同一个错误`io.EOF`：
```go
in := bufio.NewReader(os.Stdin)
for {
    r, _, err := in.ReadRune()
    if err == io.EOF {
        break // finished reading
    }
    if err != nil {
        return fmt.Errorf("read failed:%v", err)
    }
    // ...use r…
}
```

### 可变参数函数
可变参数函数值的是参数数量可变的函数。比如`fmt.Printf`，`Printf`接收一个的必备参数，之后接收任意个数的后续参数。

在参数列表的最后一个参数类型之前加上省略符号`...`，表示该函数会接收任意数量的该类型参数。
```go
func sum(vals ...int) int {
	total := 0
	for _, val := range vals {
			total += val
	}
	return total
}

// 调用
fmt.Println(sum())           // "0"
fmt.Println(sum(3))          // "3"
fmt.Println(sum(1, 2, 3, 4)) // "10"

// 还可以使用类似 ES6 的解构赋值的语法
values := []int{1, 2, 3, 4}
fmt.Println(sum(values...)) // "10"
```

### 关键字 defer
在普通函数或方法前加关键字`defer`，会使函数或方法延迟执行，直到包含该`defer`语句的函数执行完毕时（**无论函数是否出错**），
`defer`后的函数才会被执行。

`defer`语句一般被用于处理成对的操作，如打开、关闭、连接、断开连接、加锁、释放锁。因为`defer`可以保证让你更任何情况下，
资源都会被释放。
```go
package ioutil
func ReadFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
			return nil, err
	}
	defer f.Close()
	return ReadAll(f)
}

// 互斥锁
var mu sync.Mutex
var m = make(map[string]int)
func lookup(key string) int {
	mu.Lock()
	defer mu.Unlock()
	return m[key]
}

// 记录何时进入和退出函数
func bigSlowOperation() {
	defer trace("bigSlowOperation")() // 运行trace函数，记录了进入函数的时间，并返回一个函数值，这个函数值会延迟执行
	extra parentheses
	// ...lots of work…
	time.Sleep(10 * time.Second) // simulate slow
	operation by sleeping
}
func trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s", msg)
	return func() { 
		log.Printf("exit %s (%s)", msg,time.Since(start)) 
	}
}

// 观察函数的返回值
func double(x int) (result int) { // 有名返回值
  // 由于 defer 在 return 之后执行，所以这里的result 就是函数最终的返回值
	defer func() { fmt.Printf("double(%d) = %d\n", x,result) }()

	return x + x
}

_ = double(4) // 输出 "double(4) = 8"
```
上面的例子中我们知道`defer`函数可以观察函数返回值，`defer`函数还可以修改函数的返回值：
```go
func triple(x int) (result int) {
	defer func() { result += x }()
	return double(x)
}
fmt.Println(triple(4)) // "12"
```

#### 如果一个函数中有多条defer语句，那么那几个defer函数调用的执行顺序是怎样的
在同一个函数中，defer函数调用的执行顺序与它们分别所属的defer语句的出现顺序（更严谨地说，是执行顺序）完全相反。

在defer语句每次执行的时候，Go 语言会把它携带的defer函数及其参数值另行存储到一个队列中。

这个队列与该defer语句所属的函数是对应的，并且，它是先进后出（FILO）的，相当于一个栈。

在需要执行某个函数中的defer函数调用的时候，Go 语言会先拿到对应的队列，然后从该队列中一个一个地取出defer函数及其参数值，并逐个执行调用。

### Panic 异常
Go 运行时错误会引起`painc`异常。
一般而言，当`panic`异常发生时，程序会中断运行，并立即执行在该`goroutine`中被延迟的函数（defer 机制）。随后，程序崩溃并输出日志信息。

由于`panic`会引起程序的崩溃，因此`panic`一般用于严重错误，如程序内部的逻辑不一致。但是对于大部分漏洞，我们应该使用 Go 提供的错误机制，
而不是`panic`，尽量避免程序的崩溃。

#### panic 函数
`panic`函数接受任何值作为参数。当某些不应该发生的场景发生时，我们就应该调用`panic`。

#### panic 详情中都有什么
```bash
panic: runtime error: index out of range

goroutine 1 [running]:
main.main()
/Users/haolin/GeekTime/Golang_Puzzlers/src/puzzlers/article19/q0/demo47.go:5 +0x3d
exit status 2
```
第一行是`panic: runtime error: index out of range`。其中的`runtime error`的含义是，这是一个`runtime`代码包中抛出的`panic`。

`goroutine 1 [running]`，它表示有一个 ID 为1的 goroutine 在此 panic 被引发的时候正在运行。这里的 ID 其实并不重要。

`main.main()`表明了这个 goroutine 包装的go函数就是命令源码文件中的那个`main`函数，也就是说这里的 goroutine 正是**主 goroutine**。

再下面的一行，指出的就是这个 goroutine 中的哪一行代码在此 panic 被引发时正在执行。含了此行代码在其所属的源码文件中的行数，以及这个源码文件的绝对路径。

`+0x3d`代表的是：此行代码相对于其所属函数的入口程序计数偏移量。用处并不大。

`exit status 2`表明我的这个程序是以退出状态码2结束运行的。**在大多数操作系统中，只要退出状态码不是0，都意味着程序运行的非正常结束。**在 Go 语言中，
**因 panic 导致程序结束运行的退出状态码一般都会是2**。


#### 从 panic 被引发到程序终止运行的大致过程是什么

此行代码所属函数的执行随即终止。紧接着，控制权并不会在此有片刻停留，它又会立即转移至再上一级的调用代码处。控制权如此一级一级地沿着调用栈的反方向传播至顶端，
也就是我们编写的最外层函数那里。

这里的最外层函数指的是go函数，对于主 goroutine 来说就是main函数。但是控制权也不会停留在那里，而是被 Go 语言运行时系统收回。

随后，程序崩溃并终止运行，承载程序这次运行的进程也会随之死亡并消失。与此同时，在这个控制权传播的过程中，panic 详情会被逐渐地积累和完善，并会在程序终止之前被打印出来。

#### 怎样让 panic 包含一个值，以及应该让它包含什么样的值
其实很简单，在调用panic函数时，把某个值作为参数传给该函数就可以了。`panic`函数的唯一一个参数是空接口（也就是`interface{}`）类型的，所以从语法上讲，它可以接受任何类型的值。

但是，我们**最好传入`error`类型的错误值，或者其他的可以被有效序列化的值。这里的“有效序列化”指的是，可以更易读地去表示形式转换**。

### Recover 捕获异常
一般情况下，我们不能因为某个处理函数引发的`panic`异常，杀掉整个进程，可以使用`recover`函数恢复`panic`异常。

`panic`时会调用`recover`，但是`recover`不能滥用，可能会引起资源泄漏或者其他问题。我们可以将`panic value`设置成特殊类型，
来标识某个`panic`是否应该被恢复。
```go
func soleTitle(doc *html.Node) (title string, err error) {
	type bailout struct{}
	defer func() {
		switch p := recover(); p {
            case nil:       // no panic
            case bailout{}: // "expected" panic
                err = fmt.Errorf("multiple title elements")
            default:
                panic(p) // unexpected panic; carry on panicking
		}
	}()
  ...
}
```

上面的代码，`deferred`函数调用`recover`，并检查`panic value`。当`panic value`是`bailout{}`类型时，`deferred`函数生成一个`error`返回给调用者。
当`panic value`是其他`non-nil`值时，表示发生了未知的`pani`c异常。

#### 正确调用 recover 函数
```go
package main

import (
    "fmt"
    "errors"
)

func main() {
    fmt.Println("Enter function main.")
    // 引发 panic。
    panic(errors.New("something wrong"))
    p := recover()
    fmt.Printf("panic: %s\n", p)
    fmt.Println("Exit function main.")
}
```
上面的代码，`recover`函数调用并不会起到任何作用，甚至都没有机会执行。因为panic 一旦发生，控制权就会讯速地沿着调用栈的反方向传播。所以，在panic函数调用之后的代码，
根本就没有执行的机会。

先调用recover函数，再调用panic函数会怎么样呢？
如果在我们调用recover函数时未发生 panic，那么该函数就不会做任何事情，并且只会返回一个`nil`。

`defer`语句调用recover函数才是正确的打开方式。

无论函数结束执行的原因是什么，其中的defer函数调用都会在它即将结束执行的那一刻执行。即使导致它执行结束的原因是一个 panic 也会是这样。

要注意，我们要**尽量把defer语句写在函数体的开始处，因为在引发 panic 的语句之后的所有语句，都不会有任何执行机会**。

### 传入函数的那些参数值后来怎么样了
```go
package main

import "fmt"

func main() {
	array1 := [3]string{"a", "b", "c"}
	fmt.Printf("The array: %v\n", array1)
	array2 := modifyArray(array1)
	fmt.Printf("The modified array: %v\n", array2)
	fmt.Printf("The original array: %v\n", array1)
}

func modifyArray(a [3]string) [3]string {
	a[1] = "x"
	return a
}
```
在`main`函数中声明了一个数组`array1`，然后把它传给了函数`modify`，`modify`对参数值稍作修改后将其作为结果值返回。`main`函数中的代码拿到这个结果之后打印了它（即`array2`），以及原来的数组`array1`。关键问题是，原数组会因`modify`函数对参数值的修改而改变吗？

答案是：原数组不会改变。为什么呢？原因是，**所有传给函数的参数值都会被复制，函数在其内部使用的并不是参数值的原值，而是它的副本**。

由于数组是值类型，所以每一次复制都会拷贝它，以及它的所有元素值。

注意，**对于引用类型，比如：切片、字典、通道，像上面那样复制它们的值，只会拷贝它们本身而已，并不会拷贝它们引用的底层数据。也就是说，这时只是浅表复制，而不是深层复制**。

以切片值为例，如此复制的时候，只是拷贝了它指向底层数组中某一个元素的指针，以及它的长度值和容量值，而它的底层数组并不会被拷贝。