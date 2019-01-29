# WEB 编程
## 访问网络服务
socket，常被翻译为套接字，它应该算是网络编程世界中最为核心的知识之一了。

所谓 socket，是一种 IPC 方法。IPC 是 Inter-Process Communication 的缩写，可以被翻译为进程间通信。顾名思义，IPC 这个概念（或者说规范）主要定义的是多个进程之间，
相互通信的方法。

这些方法主要包括：系统信号（signal）、管道（pipe）、套接字 （socket）、文件锁（file lock）、消息队列（message queue）、信号灯（semaphore，有的地方也称之为信号量）等。
现存的主流操作系统大都对 IPC 提供了强有力的支持，尤其是 socket。

`os/signal`代码包中就有针对系统信号的 API。

又比如，`os.Pipe`函数可以创建命名管道，而`os/exec`代码包则对另一类管道（匿名管道）提供了支持。

对于 **socket，Go 语言与之相应的程序实体都在其标准库的`net`代码包中**。

所谓的**系统调用**，你可以理解为特殊的 C 语言函数。它们是**连接应用程序和操作系统内核的桥梁，也是应用程序使用操作系统功能的唯一渠道**。

syscall包中的`Socket`函数本身是平台不相关的。在其底层，Go 语言为它支持的每个操作系统都做了适配，这才使得这个函数无论在哪个平台上，总是有效的。
Go 语言的net代码包中的很多程序实体，都会直接或间接地使用到`syscall.Socket`函数。

### net.Dial函数的第一个参数`network`
`net.Dial`函数会接受两个参数，分别名为`network`和`address`，都是`string`类型的。

`network`常用的可选值一共有 9 个：
- "tcp"：代表 TCP 协议，其基于的 IP 协议的版本根据参数address的值自适应。
- "tcp4"：代表基于 IP 协议第四版的 TCP 协议。
- "tcp6"：代表基于 IP 协议第六版的 TCP 协议。
- "udp"：代表 UDP 协议，其基于的 IP 协议的版本根据参数address的值自适应。
- "udp4"：代表基于 IP 协议第四版的 UDP 协议。
- "udp6"：代表基于 IP 协议第六版的 UDP 协议。
- "unix"：代表 Unix 通信域下的一种内部 socket 协议，以 SOCK_STREAM 为 socket 类型。
- "unixgram"：代表 Unix 通信域下的一种内部 socket 协议，以 SOCK_DGRAM 为 socket 类型。
- "unixpacket"：代表 Unix 通信域下的一种内部 socket 协议，以 SOCK_SEQPACKET 为 socket 类型。

### net.DialTimeout函数
调用`net.DialTimeout`函数时给定超时时间，这里的超时时间，代表着函数为网络连接建立完成而等待的最长时间。这是一个相对的时间。它会由这个函数的参数`timeout`的值表示。
### http.Client类型中的Transport字段
`http.Client`类型中的`Transport`字段代表着：**向网络服务发送 HTTP 请求，并从网络服务接收 HTTP 响应的操作过程**。

这个字段是`http.RoundTripper`接口类型的，它有一个由`http.DefaultTransport`变量代表的缺省值（以下简称DefaultTransport）。当我们在初始化一个`http.Client`
类型的值（以下简称Client值）的时候，如果**没有显式地为该字段赋值，那么这个`Client`值就会直接使用`DefaultTransport`**。

`http.Client`类型的`Timeout`字段，代表的正是前面所说的单次 HTTP 事务的**超时时间**，它是`time.Duration`类型的。它的**零值是可用的，用于表示没有设置超时时间**。

### DefaultTransport
`DefaultTransport`的实际类型是`*http.Transport`，后者即为`http.RoundTripper`接口的默认实现。**这个类型是可以被复用**的，也推荐被复用，同时，它也是**并发安全**的。

`http.Transport`类型，会在内部使用一个`net.Dialer`类型的值（以下简称`Dialer`值），并且，它会把该值的`Timeout`字段的值，设定为30秒。也就是说，
这个`Dialer`值如果在 30 秒内还没有建立好网络连接，那么就会被判定为操作超时。

`http.Transport`类型还包含了很多其他的字段，其中有一些字段是关于操作超时的。
- `IdleConnTimeout`：含义是空闲的连接在多久之后就应该被关闭。`DefaultTransport`会把该字段的值设定为90秒。**如果该值为0，那么就表示不关闭空闲的连接。注意，
这样很可能会造成资源的泄露**。
- `ResponseHeaderTimeout`：含义是，从客户端把请求完全递交给操作系统到从操作系统那里接收到响应报文头的最大时长。`DefaultTransport`并没有设定该字段的值。
- `ExpectContinueTimeout`：含义是，在客户端递交了请求报文头之后，等待接收第一个响应报文头的最长时间。在客户端想要使用 HTTP 的“POST”方法把一个很大的报文体发送给服务端的时候，
它可以先通过发送一个包含了“Expect: 100-continue”的请求报文头，来询问服务端是否愿意接收这个大报文体。这个字段就是用于设定在这种情况下的超时时间的。注意，如果该字段的值不大于0，
那么无论多大的请求报文体都将会被立即发送出去。这样可能会造成网络资源的浪费。`DefaultTransport`把该字段的值设定为了1秒。
- `TLSHandshakeTimeout`：TLS 是 Transport Layer Security 的缩写，可以被翻译为传输层安全。这个字段代表了基于 TLS 协议的连接在被建立时的握手阶段的超时时间。若该值为0，
则表示对这个时间不设限。`DefaultTransport`把该字段的值设定为了10秒。
- `MaxIdleConns`字段对空闲连接的总数做出限定，`DefaultTransport`把`MaxIdleConns`字段的值设定为了100。
- `MaxIdleConnsPerHost`字段限定的则是，该Transport值访问的每一个网络服务的最大空闲连接数。
- `MaxConnsPerHost`限制的是，针对某一个Transport值访问的每一个网络服务的最大连接数，不论这些连接是否是空闲的。并且，该字段没有相应的缺省值，它的零值表示不对此设限。

### 空闲连接
HTTP 协议有一个请求报文头叫做“Connection”。在 HTTP 协议的 1.1 版本中，这个报文头的值默认是“keep-alive”。

在这种情况下的网络连接都是持久连接，它们会在当前的 HTTP 事务完成后仍然保持着连通性，因此是可以被复用的。

既然连接可以被复用，那么就会有两种可能。一种可能是，针对于同一个网络服务，有新的 HTTP 请求被递交，该连接被再次使用。另一种可能是，不再有对该网络服务的 HTTP 请求，该连接被闲置。

后一种可能就产生了空闲的连接。

如果我们想彻底地杜绝空闲连接的产生，那么可以在初始化Transport值的时候把它的**`DisableKeepAlives`字段的值设定为`true`**。这时，HTTP 请求的“Connection”报文头的值就会被
设置为“close”。这会告诉网络服务，这个网络连接不必保持，当前的 HTTP 事务完成后就可以断开它了。但是，每当一个 HTTP 请求被递交时，就都会产生一个新的网络连接。这样做会明显地加重
网络服务以及客户端的负载，并会让每个 HTTP 事务都耗费更多的时间。所以，在一般情况下，我们都不要去设置这个`DisableKeepAlives`字段。

## 搭建一个web
```go
package main

import (
  "fmt"
  "net/http"
  "strings"
  "log"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()  //解析参数，默认是不会解析的
  fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
  fmt.Println("path", r.URL.Path)
  fmt.Println("scheme", r.URL.Scheme)
  fmt.Println(r.Form["url_long"])
  for k, v := range r.Form {
    fmt.Println("key:", k)
    fmt.Println("val:", strings.Join(v, ""))
  }
  fmt.Fprintf(w, "Hello World!") //这个写入到w的是输出到客户端的
}

func main() {
  http.HandleFunc("/", sayhelloName) //设置访问的路由
  err := http.ListenAndServe(":9090", nil) //设置监听的端口
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
```

使用`http://localhost:9090/?url_long=111&url_long=222`访问，观察输出结果。

Go 是通过函数`ListenAndServe`来监听端口，这个函数的底层实现就是初始化一个`server`对象，然后调用了`net.Listen("tcp", addr)`，
也就是底层用TCP协议搭建了一个服务，然后监控我们设置的端口。
`http`包源码：
```go
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("http: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		tempDelay = 0
		c, err := srv.newConn(rw)
		if err != nil {
				continue
		}
		go c.serve()
	}
}
```

`srv.Serve(net.Listener)`函数，这个函数就是处理接收客户端的请求信息。这个函数里面起了一个`for{}`，首先通过`Listener`接收请求，其次创建一个`Conn`，最后单独开了一个`goroutin`e，把这个请求的数据当做参数扔给这个`conn`去服务：`go c.serve()`。这个就是高并发体现了，用户的每一次请求都是在一个新的`goroutine`去服务，相互不影响。

### http包
Go的`http`有两个核心功能：`Conn`、`ServeMux`。

Go为了实现高并发和高性能, 使用了`goroutines`来处理`Conn`的读写事件, 这样每个请求都能保持独立，相互不会阻塞，可以高效的响应网络事件。

http包的路由器是怎么实现的？
```go
type ServeMux struct {
	mu sync.RWMutex   //锁，由于请求涉及到并发处理，因此这里需要一个锁机制
	m  map[string]muxEntry  // 路由规则，一个string对应一个mux实体，这里的string就是注册的路由表达式
	hosts bool // 是否在任意的规则中带有host信息
}

type muxEntry struct {
	explicit bool   // 是否精确匹配
	h        Handler // 这个路由表达式对应哪个handler
	pattern  string  //匹配字符串
}

type Handler interface {
  ServeHTTP(ResponseWriter, *Request)  // 路由实现器
}
```
`Handler`是一个接口，但是前一小节中的`sayhelloName`函数并没有实现`ServeHTTP`这个接口。但是在`http`包里面还定义了一个类型`HandlerFunc`,
我们定义的函数`sayhelloName`就是这个`HandlerFunc`调用之后的结果，这个类型默认就实现了`ServeHTTP`这个接口，即我们调用了`HandlerFunc(f)`,强制类型转换`f`成为`HandlerFunc`类型，这样`f`就拥有了`ServeHTTP`方法。
```go
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
  f(w, r)
}
```

路由分发：
```go
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	if r.RequestURI == "*" {
		w.Header().Set("Connection", "close")
		w.WriteHeader(StatusBadRequest)
		return
	}
	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}
```
如果是`*`那么关闭链接，不然调用`mux.Handler(r)`返回对应设置路由的处理`Handler`，然后执行`h.ServeHTTP(w, r)`。
`mux.Handler(r)`：
```go
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) {
    if r.Method != "CONNECT" {
        if p := cleanPath(r.URL.Path); p != r.URL.Path {
            _, pattern = mux.handler(r.Host, p)
            return RedirectHandler(p, StatusMovedPermanently), pattern
        }
    }    
    return mux.handler(r.Host, r.URL.Path)
}

func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
    mux.mu.RLock()
    defer mux.mu.RUnlock()

    // Host-specific pattern takes precedence over generic ones
    if mux.hosts {
        h, pattern = mux.match(host + path)
    }
    if h == nil {
        h, pattern = mux.match(path)
    }
    if h == nil {
        h, pattern = NotFoundHandler(), ""
    }
    return
}
```
根据用户请求的URL和路由器里面存储的`map`去匹配的，当匹配到之后返回存储的`handler`，调用这个`handler`的`ServeHTT`P接口就可以执
行到相应的函数了。

#### http.Server类型的ListenAndServe方法

`http.Server`类型的`ListenAndServe`方法的功能是：监听一个基于 TCP 协议的网络地址，并对接收到的 HTTP 请求进行处理。这个方法会默认开启针对网络连接的存活探测机制，以保证连接是
持久的。同时，该方法会一直执行，直到有严重的错误发生或者被外界关掉。当被外界关掉时，它会返回一个由`http.ErrServerClosed`变量代表的错误值。

ListenAndServe方法主要会做下面这几件事情：
1. 检查当前的`http.Server`类型的值（以下简称当前值）的`Addr`字段。该字段的值代表了当前的网络服务需要使用的网络地址，即：IP 地址和端口号. 如果这个字段的值为空字符串，
那么就用":http"代替。也就是说，使用任何可以代表本机的域名和 IP 地址，并且端口号为80。
2. 通过调用`net.Listen`函数在已确定的网络地址上启动基于 TCP 协议的监听。
3. 检查`net.Listen`函数返回的错误值。如果该错误值不为`nil`，那么就直接返回该值。否则，通过调用当前值的`Serve`方法准备接受和处理将要到来的 HTTP 请求。

#### http.Server类型的Serve方法是怎样接受和处理 HTTP 请求的
在一个`for`循环中，网络监听器的`Accept`方法会被不断地调用，该方法会返回两个结果值；第一个结果值是`net.Conn`类型的，它会代表包含了新到来的 HTTP 请求的网络连接；
第二个结果值是代表了可能发生的错误的`error`类型值。

如果这个错误值不为`nil`，除非它代表了一个暂时性的错误，否则循环都会被终止。如果是暂时性的错误，那么循环的下一次迭代将会在一段时间之后开始执行。

如果这里的`Accept`方法没有返回非`nil`的错误值，那么这里的程序将会先把它的第一个结果值包装成一个`*http.conn`类型的值（以下简称conn值），然后通过在新的 goroutine 中调用
这个`conn`值的`serve`方法，来对当前的 HTTP 请求进行处理。

### 自己实现一个简易的路由器
Go其实支持外部实现的路由器，`ListenAndServe`的第二个参数就是用以配置外部路由器的，它是一个`Handler`接口，即外部路由器只要实
现了`Handler`接口就可以,我们可以在自己实现的路由器的`ServeHTTP`里面实现自定义路由功能。
```go
package main

import (
    "fmt"
    "net/http"
)

type MyMux struct {
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/" {
        sayhelloName(w, r)
        return
    }
    http.NotFound(w, r)
    return
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello myroute!")
}

func main() {
    mux := &MyMux{}
    http.ListenAndServe(":9090", mux)
}
```

## 表单
```go
func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析url传递的参数，对于POST则解析响应包的主体（request body）
	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello world!") //这个写入到w的是输出到客户端的
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}
```
上面的代码可以看出通过`r.Method`来获取请求方法。
`r.Form`里面包含了所有请求的参数，比如URL中`query-string`、`POST`的数据、`PUT`的数据，所有当你在URL的`query-string`字段和`POST`冲突时，
会保存成一个`slice`，里面存储了多个值，Go官方文档中说在接下来的版本里面将会把`POST`、`GET`这些数据分离开来。

### 预防跨站脚本
对`XSS`最佳的防护应该结合以下两种方法：
- 验证所有输入数据，有效检测攻击
- 对所有输出数据进行适当的处理，以防止任何已成功注入的脚本在浏览器端运行。

Go的`html/template`里面带有下面几个函数可以帮你转义：
```go
func HTMLEscape(w io.Writer, b []byte) //把b进行转义之后写到w
func HTMLEscapeString(s string) string //转义s之后返回结果字符串
func HTMLEscaper(args ...interface{}) string //支持多个参数一起转义，返回结果字符串


fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) //输出到服务器端
fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
template.HTMLEscape(w, []byte(r.Form.Get("username"))) //输出到客户端
```

### 文件上传
```go
http.HandleFunc("/upload", upload)

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method) //获取请求的方法
    if r.Method == "GET" {
        crutime := time.Now().Unix()
        h := md5.New()
        io.WriteString(h, strconv.FormatInt(crutime, 10))
        token := fmt.Sprintf("%x", h.Sum(nil))

        t, _ := template.ParseFiles("upload.gtpl")
        t.Execute(w, token)
    } else {
        r.ParseMultipartForm(32 << 20)
        file, handler, err := r.FormFile("uploadfile")
        if err != nil {
            fmt.Println(err)
            return
        }
        defer file.Close()
        fmt.Fprintf(w, "%v", handler.Header)
        f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer f.Close()
        io.Copy(f, file)
    }
}
```

处理文件上传需要调用`r.ParseMultipartForm`，里面的参数表示`maxMemory`，调用`ParseMultipartForm`之后，上传的文件存储在`maxMemory`大小的
内存里面，如果文件大小超过了`maxMemory`，那么剩下的部分将存储在系统的临时文件中。我们可以通过`r.FormFile`获取上面的文件句柄，然后实例中使
用了`io.Copy`来存储文件。

## 访问数据库
Go定义了`database/sql`接口。