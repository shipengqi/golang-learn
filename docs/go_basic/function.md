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

在函数体中，函数的形参作为局部变量，被初始化为调用者提供的值（函数调用必须按照声明顺序为所有参数提供实参）。函数的形参和有名返回值（也就是对返回值命名）作为函数最外层的局部变量，被存储在相同的词法块中。

**Go 语言使用的是值传递，当我们传一个参数值到被调用函数里面时，实际上是传了这个值的一份copy，当在被调用函数中修改参数值的时候，调用函数中相应实参不会发生任何变化，因为数值变化只作用在copy上。但是如果是引用传递，在调用函数时将实际参数的地址传递到函数中，那么在函数中对参数所进行的修改，将影响到实际参数。**

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

### 闭包
Go 语言支持匿名函数，可作为闭包。
```go
// 返回一个函数
func getSequence() func() int {
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
func sum(vals...int) int {
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

### Panic 异常
Go 运行时错误会引起`painc`异常。
一般而言，当`panic`异常发生时，程序会中断运行，并立即执行在该`goroutine`中被延迟的函数（defer 机制）。随后，程序崩溃并输出日志信息。

由于`panic`会引起程序的崩溃，因此`panic`一般用于严重错误，如程序内部的逻辑不一致。但是对于大部分漏洞，我们应该使用 Go 提供的错误机制，
而不是`panic`，尽量避免程序的崩溃。

#### panic 函数
`panic`函数接受任何值作为参数。当某些不应该发生的场景发生时，我们就应该调用`panic`。

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