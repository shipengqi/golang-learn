---
title: http
---
# http


`net/http` 可以用来处理 HTTP 协议，包括 HTTP 服务器和 HTTP 客户端，主要组成：

- Request，HTTP 请求对象
- Response，HTTP 响应对象
- Client，HTTP 客户端
- Server，HTTP 服务端

创建一个 server ：

```go
package main

import (
	"fmt"
	"net/http"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "hello")
}

func main() {
	http.HandleFunc("/", MyHandler)
	_ = http.ListenAndServe(":8080", nil)
}
```

发送请求：
```go
resp, err := http.Get("http://example.com/") // GET 

resp, err := http.Post("http://example.com/") // POST

resp, err := http.PostForm("http://example.com/", url.Values{"foo": "bar"}) // 提交表单
```

## Request

`Request` 对象是对 http 请求报文的抽象。包括起始行, `Headers`, `Body` 等。

使用 `http.NewRequest` 函数创建一个 `Request` 请求对象：
```go
func NewRequest(method, url string, body io.Reader) (*Request, error)
```

`Request` 对象主要用于数据的存储，结构：
```go
type Request struct {
       
        Method string // HTTP方法
        URL *url.URL // URL
        
        Proto      string // "HTTP/1.0"
        ProtoMajor int    // 1
        ProtoMinor int    // 0

        Header Header // 报文头
        Body io.ReadCloser // 报文体
        GetBody func() (io.ReadCloser, error)
        ContentLength int64 // 报文长度
        TransferEncoding []string // 传输编码
        Close bool // 关闭连接
        Host string // 主机名
        
        Form url.Values // 
        PostForm url.Values // POST表单信息
        MultipartForm *multipart.Form // multipart，

        Trailer Header
        RemoteAddr string
        RequestURI string
        TLS *tls.ConnectionState
        Cancel <-chan struct{}
        Response *Response
}
```

```go
// 从 b 中读取和解析一个请求. 
func ReadRequest(b *bufio.Reader) (req *Request, err error)

// 给 request 添加 cookie, AddCookie 向请求中添加一个 cookie.按照RFC 6265 
// section 5.4的规则, AddCookie 不会添加超过一个 Cookie 头字段.
// 这表示所有的 cookie 都写在同一行, 用分号分隔（cookie内部用逗号分隔属性） 
func (r *Request) AddCookie(c *Cookie)

// 返回 request 中指定名 name 的 cookie，如果没有发现，返回 ErrNoCookie 
func (r *Request) Cookie(name string) (*Cookie, error)

// 返回该请求的所有cookies 
func (r *Request) Cookies() []*Cookie

// 利用提供的用户名和密码给http基本权限提供具有一定权限的header。
// 当使用http基本授权时，用户名和密码是不加密的 
func (r *Request) SetBasicAuth(username, password string)

// 如果在request中发送，该函数返回客户端的user-Agent
func (r *Request) UserAgent() string

// 对于指定格式的key，FormFile返回符合条件的第一个文件，如果有必要的话，
// 该函数会调用ParseMultipartForm和ParseForm。 
func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)

// 返回key获取的队列中第一个值。在查询过程中post和put中的主题参数优先级
// 高于url中的value。为了访问相同key的多个值，调用ParseForm然后直接
// 检查RequestForm。 
func (r *Request) FormValue(key string) string

// 如果这是一个有多部分组成的post请求，该函数将会返回一个MIME 多部分reader，
// 否则的话将会返回一个nil和error。使用本函数代替ParseMultipartForm
// 可以将请求body当做流stream来处理。 
func (r *Request) MultipartReader() (*multipart.Reader, error)

// 解析URL中的查询字符串，并将解析结果更新到r.Form字段。对于POST或PUT
// 请求，ParseForm还会将body当作表单解析，并将结果既更新到r.PostForm也
// 更新到r.Form。解析结果中，POST或PUT请求主体要优先于URL查询字符串
// （同名变量，主体的值在查询字符串的值前面）。如果请求的主体的大小没有被
// MaxBytesReader函数设定限制，其大小默认限制为开头10MB。
// ParseMultipartForm会自动调用ParseForm。重复调用本方法是无意义的。
func (r *Request) ParseForm() error 

// ParseMultipartForm将请求的主体作为multipart/form-data解析。
// 请求的整个主体都会被解析，得到的文件记录最多 maxMemery字节保存在内存，
// 其余部分保存在硬盘的temp文件里。如果必要，ParseMultipartForm会
// 自行调用 ParseForm。重复调用本方法是无意义的。
func (r *Request) ParseMultipartForm(maxMemory int64) error 

// 返回post或者put请求body指定元素的第一个值，其中url中的参数被忽略。
func (r *Request) PostFormValue(key string) string 

// 检测在request中使用的http协议是否至少是major.minor 
func (r *Request) ProtoAtLeast(major, minor int) bool

// 如果request中有refer，那么refer返回相应的url。Referer在request
// 中是拼错的，这个错误从http初期就已经存在了。该值也可以从Headermap中
// 利用Header["Referer"]获取；在使用过程中利用Referer这个方法而
// 不是map的形式的好处是在编译过程中可以检查方法的错误，而无法检查map中
// key的错误。
func (r *Request) Referer() string 

// Write方法以有线格式将HTTP/1.1请求写入w（用于将请求写入下层TCPConn等）
// 。本方法会考虑请求的如下字段：Host URL Method (defaults to "GET")
//  Header ContentLength TransferEncoding Body如果存在Body，
// ContentLength字段<= 0且TransferEncoding字段未显式设置为
// ["identity"]，Write方法会显式添加”Transfer-Encoding: chunked”
// 到请求的头域。Body字段会在发送完请求后关闭。
func (r *Request) Write(w io.Writer) error 

// 该函数与Write方法类似，但是该方法写的request是按照http代理的格式去写。
// 尤其是，按照RFC 2616 Section 5.1.2，WriteProxy会使用绝对URI
// （包括协议和主机名）来初始化请求的第1行（Request-URI行）。无论何种情况，
// WriteProxy都会使用r.Host或r.URL.Host设置Host头。
func (r *Request) WriteProxy(w io.Writer) error 
```

## Response

Response 也是一个数据对象，描述 HTTP 响应：
```go
type Response struct {
        Status     string // HTTP 状态码
        StatusCode int    // 状态码 200
        Proto      string // 版本号 "HTTP/1.0"
        ProtoMajor int    // 主版本号 
        ProtoMinor int    // 次版本号

        Header Header // 响应报文头
        Body io.ReadCloser // 响应报文体
        ContentLength int64 // 报文长度
        TransferEncoding []string // 报文编码
        Close bool 
        Trailer Header
        Request *Request // 请求对象
        TLS *tls.ConnectionState
}
```

```go
// ReadResponse从r读取并返回一个HTTP 回复。req参数是可选的，指定该回复
// 对应的请求（即是对该请求的回复）。如果是nil，将假设请 求是GET请求。
// 客户端必须在结束resp.Body的读取后关闭它。读取完毕并关闭后，客户端可以
// 检查resp.Trailer字段获取回复的 trailer的键值对。
func ReadResponse(r *bufio.Reader, req *Request) (*Response, error)

// 解析cookie并返回在header中利用set-Cookie设定的cookie值。
func (r *Response) Cookies() []*Cookie 

// 返回response中Location的header值的url。如果该值存在的话，则对于
// 请求问题可以解决相对重定向的问题，如果该值为nil，则返回ErrNOLocation。
func (r *Response) Location() (*url.URL, error) 

// 判定在response中使用的http协议是否至少是major.minor的形式。
func (r *Response) ProtoAtLeast(major, minor int) bool 

// 将response中信息按照线性格式写入w中。
func (r *Response) Write(w io.Writer) error 
```


## client

前面以 `http.Get("http://example.com/")` `Get` 或 `Post` 函数发送请求，就是通过绑定一个默认 `Client` 实现的。
使用 `Client` 要先初始化一个 `Client` 对象。`Client` 具有 `Do`，`Get`，`Head`，`Post` 以及 `PostForm` 等方法。

```go
package main

import "net/http"

func main() {
    client := http.Client()
    res, err := client.Get("http://www.google.com")
}
```
对于常用 HTTP 动词，`Client` 对象对应的函数，下面的这些方法与 `http.Get` 等方法一致：
```go
func (c *Client) Get(url string) (resp *Response, err error)

func (c *Client) Head(url string) (resp *Response, err error)

func (c *Client) Post(url string, contentType string, body io.Reader) (resp *Response, err error)

func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error)
```
但是在很多情况下，需要支持对 headers，cookies 等的设定，上面提供的方法就不能满足需求了。就需要使用 `Do` 方法，
```go
func (c *Client) Do(req *Request) (resp *Response, err error)
```

`http.NewRequest` 可以灵活的对 `Request` 进行配置，然后再使用 `http.Client` 的 `Do` 方法发送这个 `Request` 请求。

**模拟 HTTP Request**：

```go
// 简式声明一个 http.Client
client := &http.Client{}

// 构建 Request
request, err := http.NewRequest("GET", "http://www.baidu.com", nil)
if err != nil {
    fmt.Println(err)
}

// 使用 http.Cookie 结构体初始化一个 cookie 键值对
cookie := &http.Cookie{Name: "userId", Value: strconv.Itoa(12345)}

// AddCookie
request.AddCookie(cookie)

// 设置 Header
request.Header.Set("Accept", "text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8")
request.Header.Set("Accept-Charset", "GBK, utf-8;q=0.7, *;q=0.3")
request.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
request.Header.Set("Accept-Language", "zh-CN, zh;q=0.8")
request.Header.Set("Cache-Control", "max-age=0")
request.Header.Set("Connection", "keep-alive")

// 使用 Do 方法发送请求
response, err := client.Do(request)
if err != nil {
    fmt.Println(err)
    return
}

// 程序结束时关闭 response.Body 响应流
defer response.Body.Close()

// http Response 状态值
fmt.Println(response.StatusCode)
if response.StatusCode == 200 {

    // gzip.NewReader对压缩的返回信息解压（考虑网络传输量，http Server
    // 一般都会对响应压缩后再返回）
    body, err := gzip.NewReader(response.Body)
    if err != nil {
        fmt.Println(err)
    }

    defer body.Close()

    r, err := ioutil.ReadAll(body)
    if err != nil {
        fmt.Println(err)
    }
    // 打印出http Server返回的http Response信息
    fmt.Println(string(r))
}
```


**Ge请求**：

```go
// http.Get 实际上是 DefaultClient.Get(url)
response, err := http.Get("http://www.baidu.com")
if err != nil {
    fmt.Println(err)
}

// 程序在使用完回复后必须关闭回复的主体
defer response.Body.Close()

body, _ := ioutil.ReadAll(response.Body)
fmt.Println(string(body))
```

**Post 请求**：

```go
// application/x-www-form-urlencoded：为 POST 的 contentType
resp, err := http.Post("http://localhost:8080/login.do",
    "application/x-www-form-urlencoded", strings.NewReader("mobile=xxxxxxxxxx&isRemberPwd=1"))
if err != nil {
    fmt.Println(err)
    return
}
defer resp.Body.Close()
body, err := ioutil.ReadAll(resp.Body)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(string(body))
```

**http.PostForm 请求**：

```go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	postParam := url.Values{
		"mobile":      {"xxxxxx"},
		"isRemberPwd": {"1"},
	}
	// 数据的键值会经过 URL 编码后作为请求的 body 传递
	resp, err := http.PostForm("http://localhost：8080/login.do", postParam)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
```

## HTTP Server

`server.go` 文件中定义了一个非常重要的接口：`Handler`，另外还有一个结构体 `response`，这和 `http.Response` 结构体只有首字母大小
写不一致，这个 `response` 也是响应，只不过是专门用在服务端，和 `http.Response` 结构体是完全两回事。

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

type Server struct

// 监听 srv.Addr 然后调用 Serve 来处理接下来连接的请求
// 如果 srv.Addr 是空的话，则使用 ":http"
func (srv *Server) ListenAndServe() error 

// 监听 srv.Addr ，调用 Serve 来处理接下来连接的请求
// 必须提供证书文件和对应的私钥文件。如果证书是由
// 权威机构签发的，certFile 参数必须是顺序串联的服务端证书和 CA 证书。

func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error 

// 接受 Listener l 的连接，创建一个新的服务协程。该服务协程读取请求然后调用
// srv.Handler 来应答。实际上就是实现了对某个端口进行监听，然后创建相应的连接。 
func (srv *Server) Serve(l net.Listener) error

// 该函数控制是否 http 的 keep-alives 能够使用，默认情况下，keep-alives 总是可用的。
// 只有资源非常紧张的环境或者服务端在关闭进程中时，才应该关闭该功能。 
func (s *Server) SetKeepAlivesEnabled(v bool)

// 是一个 http 请求多路复用器，它将每一个请求的 URL 和
// 一个注册模式的列表进行匹配，然后调用和 URL 最匹配的模式的处理器进行后续操作。
type ServeMux

// 初始化一个新的 ServeMux 
func NewServeMux() *ServeMux

// 将 handler 注册为指定的模式，如果该模式已经有了 handler，则会出错 panic。
func (mux *ServeMux) Handle(pattern string, handler Handler) 

// 将 handler 注册为指定的模式 
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request))

// 根据指定的 r.Method, r.Host 以及 r.RUL.Path 返回一个用来处理给定请求的 handler。
// 该函数总是返回一个 非 nil 的 handler，如果 path 不是一个规范格式，则 handler 会
// 重定向到其规范 path。Handler 总是返回匹配该请求的的已注册模式；在内建重定向
// 处理器的情况下，pattern 会在重定向后进行匹配。如果没有已注册模式可以应用于该请求，
// 本方法将返回一个内建的 ”404 page not found” 处理器和一个空字符串模式。
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) 

// 该函数用于将最接近请求 url 模式的 handler 分配给指定的请求。 
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request)
```

`Handler` 接口是 `server.go` 中最关键的接口，如果我们仔细看这个文件的源代码，将会发现很多结构体实现了这个接口的ServeHTTP方法。

注意这个接口的注释：Handler 响应 HTTP 请求。没错，最终我们的 HTTP 服务是通过实现 `ServeHTTP(ResponseWriter, *Request)` 来达到服务端
接收客户端请求并响应。

```go
func main() {
	http.HandleFunc("/", MyHandler)
	_ = http.ListenAndServe(":8080", nil)
}
```

以上两行代码，就成功启动了一个 HTTP 服务器。我们通过 `net/http` 包源代码分析发现，调用 `Http.HandleFunc`，按顺序做了几件事：


1. `Http.HandleFunc`调用了 `DefaultServeMux` 的 `HandleFunc`

```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    DefaultServeMux.HandleFunc(pattern, handler)
}
```

2. `DefaultServeMux.HandleFunc` 调用了 `DefaultServeMux` 的 `Handle`，`DefaultServeMux` 是一个 `ServeMux` 指针变量。
而 `ServeMux` 是 Go 语言中的 `Multiplexer`（多路复用器），通过 `Handle` 匹配 `pattern` 和我们定义的 `handler`
（其实就是 `http.HandlerFunc` 函数类型变量）。

```go
var DefaultServeMux = &defaultServeMux
var defaultServeMux ServeMux

func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	mux.Handle(pattern, HandlerFunc(handler))
}
```
注意：
上面的方法命名Handle，HandleFunc和HandlerFunc，Handler（接口），他们很相似，容易混淆。记住Handle和HandleFunc和pattern 匹配有关，
也即往DefaultServeMux的map[string]muxEntry中增加对应的handler和路由规则。

接着我们看看myfunc的声明和定义：

```go
func myfunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hi")
}
```
而type HandlerFunc func(ResponseWriter, *Request) 是一个函数类型，而我们定义的myfunc的函数签名刚好符合这个函数类型。

所以http.HandleFunc("/", myfunc)，实际上是mux.Handle("/", HandlerFunc(myfunc))。

HandlerFunc(myfunc) 让myfunc成为了HandlerFunc类型，我们称myfunc为handler。而HandlerFunc类型是具有ServeHTTP方法的，而有了
ServeHTTP方法也就是实现了Handler接口。

```go
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r) // 这相当于自身的调用
}

```
现在ServeMux和Handler都和我们的myfunc联系上了，myfunc是一个Handler接口变量也是HandlerFunc类型变量，接下来和结构体server有关了。

从http.ListenAndServe的源码可以看出，它创建了一个server对象，并调用server对象的ListenAndServe方法：

```go
func ListenAndServe(addr string, handler Handler) error {
    server := &Server{Addr: addr, Handler: handler}
    return server.ListenAndServe()
}
```
而我们HTTP服务器中第二行代码：

```go
http.ListenAndServe(":8080", nil)
```
创建了一个server对象，并调用server对象的ListenAndServe方法，这里没有直接传递Handler，而是默认使用DefautServeMux作为multiplexer，
myfunc是存在于handler和路由规则中的。

Server的ListenAndServe方法中，会初始化监听地址Addr，同时调用Listen方法设置监听。

```go
for {
    rw, e := l.Accept()
    ...
    c := srv.newConn(rw)
c.setState(c.rwc, StateNew) 
go c.serve(ctx)
}
```
监听开启之后，一旦客户端请求过来，Go就开启一个协程go c.serve(ctx)处理请求，主要逻辑都在serve方法之中。

func (c *conn) serve(ctx context.Context)，这个方法很长，里面主要的一句：serverHandler{c.server}.ServeHTTP(w, w.req)。
其中w由w, err := c.readRequest(ctx)得到，因为有传递context。

还是来看源代码：

```go
type serverHandler struct {
srv *Server
}

func (sh serverHandler) ServeHTTP(rw ResponseWriter, req Request) {
handler := sh.srv.Handler
if handler == nil {
handler = DefaultServeMux
}
if req.RequestURI == "" && req.Method == "OPTIONS" {
handler = globalOptionsHandler{}
}
handler.ServeHTTP(rw, req)
}
```
从http.ListenAndServe(":8080", nil)开始，handler是nil，所以最后实际ServeHTTP方法是DefaultServeMux.ServeHTTP(rw, req)。

```go
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(StatusBadRequest)
		return
	}
	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}
```

通过func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string)，我们得到Handler h，然后执行h.ServeHTTP(w, r)方法，
也就是执行我们的myfunc函数（别忘了myfunc是HandlerFunc类型，而他的ServeHTTP(w, r)方法这里其实就是自己调用自己），把response写
到http.ResponseWriter对象返回给客户端，fmt.Fprintf(w, "hi")，我们在客户端会接收到hi 。至此整个HTTP服务执行完成。


总结下，HTTP服务整个过程大概是这样：
```go
Request -> ServeMux(Multiplexer) -> handler-> Response
```

我们再看下面代码：

```go
http.ListenAndServe(":8080", nil)
func ListenAndServe(addr string, handler Handler) error {
    server := &Server{Addr: addr, Handler: handler}
    return server.ListenAndServe()
}
```

上面代码实际上就是server.ListenAndServe()执行的实际效果，只不过简单声明了一个结构体Server{Addr: addr, Handler: handler}实例。
如果我们声明一个Server实例，完全可以达到深度自定义 http.Server的目的：


```go
package main

import (
	"fmt"
	"net/http"
)

func myfunc(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "hi")
}

func main() {
	// 更多http.Server的字段可以根据情况初始化
	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  0,
		WriteTimeout: 0,
	}
	http.HandleFunc("/", myfunc)
	_ = server.ListenAndServe()
}
```
这样服务也能跑起来，而且我们完全可以根据情况来自定义我们的Server！

还可以指定Servemux的用法:

`GOPATH\src\go42\chapter-15\15.3\7\main.go`
```go
package main

import (
	"fmt"
	"net/http"
)

func myfunc(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "hi")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", myfunc)
	_ = http.ListenAndServe(":8080", mux)
}
```

如果既指定Servemux又自定义 http.Server，因为Server中有字段Handler，所以我们可以直接把Servemux变量作为Server.Handler：


```go
package main

import (
	"fmt"
	"net/http"
)

func myfunc(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "hi")
}

func main() {
	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  0,
		WriteTimeout: 0,
	}
	mux := http.NewServeMux()
	server.Handler = mux

	mux.HandleFunc("/", myfunc)
	_ = server.ListenAndServe()
}
```

在前面pprof 包的内容中我们也用了本章开头这段代码，当我们访问http://localhost:8080/debug/pprof/ 时可以看到对应的性能分析报告。
因为我们这样导入 _"net/http/pprof" 包时，在文件 pprof.go 文件中init 函数已经定义好了handler：

```go
func init() {
	http.HandleFunc("/debug/pprof/", Index)
	http.HandleFunc("/debug/pprof/cmdline", Cmdline)
	http.HandleFunc("/debug/pprof/profile", Profile)
	http.HandleFunc("/debug/pprof/symbol", Symbol)
	http.HandleFunc("/debug/pprof/trace", Trace)
}
```
所以，我们就可以通过浏览器访问上面地址来看到报告。现在再来看这些代码，我们就明白怎么回事了！

## 36.5 自定义处理器（Custom Handlers）

自定义的Handler：

标准库http提供了Handler接口，用于开发者实现自己的handler。只要实现接口的ServeHTTP方法即可。

```go
package main

import (
	"log"
	"net/http"
	"time"
)

type timeHandler struct {
	format string
}

func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(th.format)
	_, _ = w.Write([]byte("The time is: " + tm))
}

func main() {
	mux := http.NewServeMux()

	th := &timeHandler{format: time.RFC1123}
	mux.Handle("/time", th)

	log.Println("Listening...")
	_ = http.ListenAndServe(":3000", mux)
}
```
我们知道，NewServeMux可以创建一个ServeMux实例，ServeMux同时也实现了ServeHTTP方法，因此代码中的mux也是一种handler。把它当成参数
传给http.ListenAndServe方法，后者会把mux传给Server实例。因为指定了handler，因此整个http服务就不再是DefaultServeMux，而是mux，
无论是在注册路由还是提供请求服务的时候。

任何有 func(http.ResponseWriter，*http.Request) 签名的函数都能转化为一个 HandlerFunc 类型。这很有用，因为 HandlerFunc 对象
内置了 ServeHTTP 方法，后者可以聪明又方便的调用我们最初提供的函数内容。

## 将函数作为处理器

```go
package main

import (
	"log"
	"net/http"
	"time"
)

func timeHandler(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	_, _ = w.Write([]byte("The time is: " + tm))
}

func main() {
	mux := http.NewServeMux()

	// Convert the timeHandler function to a HandlerFunc type
	th := http.HandlerFunc(timeHandler)
	// And add it to the ServeMux
	mux.Handle("/time", th)

	log.Println("Listening...")
	_ = http.ListenAndServe(":3000", mux)
}
```
创建新的server：

```go
func index(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")

    html := `<doctype html>
        <html>
        <head>
          <title>Hello World</title>
        </head>
        <body>
        <p>
          Welcome
        </p>
        </body>
</html>`
    fmt.Fprintln(w, html)
}

func main(){
    http.HandleFunc("/", index)

    server := &http.Server{
        Addr: ":8000", 
        ReadTimeout: 60 * time.Second, 
        WriteTimeout: 60 * time.Second, 
    }
    server.ListenAndServe()
}
```

## 中间件 Middleware

所谓中间件，就是连接上下级不同功能的函数或者软件，通常进行一些包裹函数的行为，为被包裹函数提供添加一些功能或行为。前文的HandleFunc就
能把签名为 func(w http.ResponseWriter, r *http.Reqeust)的函数包裹成handler。这个函数也算是中间件。

Go的HTTP中间件很简单，只要实现一个函数签名为func(http.Handler) http.Handler的函数即可。http.Handler是一个接口，接口方法我们
熟悉的为serveHTTP。返回也是一个handler。因为Go中的函数也可以当成变量传递或者或者返回，因此也可以在中间件函数中传递定义好的函数，
只要这个函数是一个handler即可，即实现或者被handlerFunc包裹成为handler处理器。

```go
func index(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")

    html := `<doctype html>
        <html>
        <head>
          <title>Hello World</title>
        </head>
        <body>
        <p>
          Welcome
        </p>
        </body>
</html>`
    fmt.Fprintln(w, html)
}

func middlewareHandler(next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        // 执行handler之前的逻辑
        next.ServeHTTP(w, r)
        // 执行完毕handler后的逻辑
    })
}

func loggingHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Started %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
        log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
    })
}

func main() {
    http.Handle("/", loggingHandler(http.HandlerFunc(index)))

    http.ListenAndServe(":8000", nil)
}
```

## 静态站点

下面代码通过指定目录，作为静态站点：
```go
package main

import (
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("D:/html/static/")))
	_ = http.ListenAndServe(":8080", nil)
}
```