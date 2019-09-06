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
// 这表示所有的 cookie 都写在同一行, 用分号分隔（cookie 内部用逗号分隔属性） 
func (r *Request) AddCookie(c *Cookie)

// 返回 request 中指定名 name 的 cookie，如果没有发现，返回 ErrNoCookie 
func (r *Request) Cookie(name string) (*Cookie, error)

// 返回该请求的所有cookies 
func (r *Request) Cookies() []*Cookie

// 利用提供的用户名和密码给 http 基本权限提供具有一定权限的 header。
// 当使用 http 基本授权时，用户名和密码是不加密的 
func (r *Request) SetBasicAuth(username, password string)

// 如果在 request 中发送，该函数返回客户端的 user-Agent
func (r *Request) UserAgent() string

// 对于指定格式的 key，FormFile 返回符合条件的第一个文件，如果有必要的话，
// 该函数会调用 ParseMultipartForm 和 ParseForm。 
func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)

// 返回 key 获取的队列中第一个值。在查询过程中 post 和 put 中的主题参数优先级
// 高于 url 中的 value。为了访问相同 key 的多个值，调用 ParseForm 然后直接
// 检查 RequestForm。 
func (r *Request) FormValue(key string) string

// 如果这是一个有多部分组成的 post 请求，该函数将会返回一个 MIME 多部分 reader，
// 否则的话将会返回一个 nil 和 error。使用本函数代替 ParseMultipartForm
// 可以将请求 body 当做流 stream 来处理。 
func (r *Request) MultipartReader() (*multipart.Reader, error)

// 解析URL中的查询字符串，并将解析结果更新到 r.Form 字段。对于POST 或 PUT
// 请求，ParseForm 还会将 body 当作表单解析，并将结果既更新到 r.PostForm 也
// 更新到 r.Form。解析结果中，POST 或 PUT 请求主体要优先于 URL 查询字符串
// （同名变量，主体的值在查询字符串的值前面）。如果请求的主体的大小没有被
// MaxBytesReader 函数设定限制，其大小默认限制为开头 10MB。
// ParseMultipartForm 会自动调用 ParseForm。重复调用本方法是无意义的。
func (r *Request) ParseForm() error 

// ParseMultipartForm 将请求的主体作为 multipart/form-data 解析。
// 请求的整个主体都会被解析，得到的文件记录最多 maxMemery 字节保存在内存，
// 其余部分保存在硬盘的 temp 文件里。如果必要，ParseMultipartForm 会
// 自行调用 ParseForm。重复调用本方法是无意义的。
func (r *Request) ParseMultipartForm(maxMemory int64) error 

// 返回 post 或者 put 请求 body 指定元素的第一个值，其中 url 中的参数被忽略。
func (r *Request) PostFormValue(key string) string 

// 检测在 request 中使用的 http 协议是否至少是 major.minor 
func (r *Request) ProtoAtLeast(major, minor int) bool

// 如果 request 中有 refer，那么 refer 返回相应的 url。Referer 在 request
// 中是拼错的(Referrer)，这个错误从 http 初期就已经存在了。该值也可以从 Headermap 中
// 利用 Header["Referer"] 获取；在使用过程中利用 Referer 这个方法而
// 不是 map 的形式的好处是在编译过程中可以检查方法的错误，而无法检查 map 中
// key 的错误。
func (r *Request) Referer() string 

// Write 方法以有线格式将 HTTP/1.1 请求写入 w（用于将请求写入下层 TCPConn 等）
// 。本方法会考虑请求的如下字段：Host URL Method (defaults to "GET")
//  Header ContentLength TransferEncoding Body如果存在Body，
// ContentLength字段<= 0且TransferEncoding字段未显式设置为
// ["identity"]，Write 方法会显式添加 ”Transfer-Encoding: chunked”
// 到请求的头域。Body 字段会在发送完请求后关闭。
func (r *Request) Write(w io.Writer) error 

// 该函数与 Write 方法类似，但是该方法写的 request 是按照 http 代理的格式去写。
// 尤其是，按照 RFC 2616 Section 5.1.2，WriteProxy 会使用绝对 URI
// （包括协议和主机名）来初始化请求的第1行（Request-URI行）。无论何种情况，
// WriteProxy 都会使用 r.Host 或 r.URL.Host 设置 Host 头。
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
// ReadResponse 从 r 读取并返回一个 HTTP 回复。req 参数是可选的，指定该回复
// 对应的请求（即是对该请求的回复）。如果是 nil，将假设请 求是 GET 请求。
// 客户端必须在结束 resp.Body 的读取后关闭它。读取完毕并关闭后，客户端可以
// 检查 resp.Trailer 字段获取回复的 trailer 的键值对。
func ReadResponse(r *bufio.Reader, req *Request) (*Response, error)

// 解析 cookie 并返回在 header 中利用 set-Cookie 设定的 cookie 值。
func (r *Response) Cookies() []*Cookie 

// 返回 response 中 Location 的 header 值的 url。如果该值存在的话，则对于
// 请求问题可以解决相对重定向的问题，如果该值为nil，则返回ErrNOLocation。
func (r *Response) Location() (*url.URL, error) 

// 判定在 response 中使用的 http 协议是否至少是 major.minor 的形式。
func (r *Response) ProtoAtLeast(major, minor int) bool 

// 将 response 中信息按照线性格式写入 w 中。
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

// 接受 l Listener 的连接，创建一个新的服务协程。该服务协程读取请求然后调用
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

`Handler` 接口是 `server.go` 中最关键的接口，如果我们仔细看这个文件的源代码，将会发现很多结构体实现了这个接口的 `ServeHTTP` 方法。

注意这个接口的注释：`Handler` 响应 HTTP 请求。没错，最终我们的 HTTP 服务是通过实现 `ServeHTTP(ResponseWriter, *Request)` 来达
到服务端接收客户端请求并响应。

```go
func main() {
	http.HandleFunc("/", MyHandler)
	_ = http.ListenAndServe(":8080", nil)
}
```

以上两行代码，就成功启动了一个 HTTP 服务器。我们通过 `net/http` 包源代码分析发现，调用 `Http.HandleFunc`，按顺序做了几件事：


1. `Http.HandleFunc` 调用了 `DefaultServeMux` 的 `HandleFunc`

```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    DefaultServeMux.HandleFunc(pattern, handler)
}
```

2. `DefaultServeMux.HandleFunc` 调用了 `DefaultServeMux` 的 `Handle`，`DefaultServeMux` 是一个 `ServeMux` 指针变量。
而 `ServeMux` 是 Go 语言中的 `Multiplexer`（多路复用器），通过 `Handle` 匹配 `pattern` 和我们定义的 `handler`
（其实就是 `http.HandlerFunc` 函数类型变量）。

```go
type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]muxEntry // 保存路由规则 和 handler
	hosts bool // whether any patterns contain hostnames
}

type muxEntry struct {
	h       Handler
	pattern string
}

var DefaultServeMux = &defaultServeMux
var defaultServeMux ServeMux

func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, HandlerFunc(handler)) // 这个 handler 就是 MyHandler
}
```
注意：
上面的方法命名 `Handle`，`HandleFunc` 和 `HandlerFunc`，`Handler`（接口），他们很相似，容易混淆。记住 **`Handle` 和 `HandleFunc` 
和 `pattern` 匹配有关，也即往 `DefaultServeMux` 的 `map[string]muxEntry` 中增加对应的 `handler` 和路由规则**。

接着我们看看 `MyHandler` 的声明和定义：

```go
func MyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}
```
而 `type HandlerFunc func(ResponseWriter, *Request)` 是一个函数类型，而我们定义的 `MyHandler` 的函数签名刚好符合这个函数类型。

所以 `http.HandleFunc("/", MyHandler)`，实际上是 `mux.Handle("/", HandlerFunc(MyHandler))`。

`HandlerFunc(MyHandler)` 让 `MyHandler` 成为了 `HandlerFunc` 类型，我们称 `MyHandler` 为 `handler`。而 `HandlerFunc` 类型是
具有 `ServeHTTP` 方法的，而有了 `ServeHTTP` 方法也就是实现了 `Handler` 接口。

```go
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r) // 这相当于自身的调用
}

```
现在 `ServeMux` 和 `Handler` 都和我们的 `MyHandler` 联系上了，`MyHandler` 是一个 `Handler` 接口变量也是 `HandlerFunc` 类型变量，
接下来和结构体 `server` 有关了。

从 `http.ListenAndServe` 的源码可以看出，它创建了一个 `server` 对象，并调用 `server` 对象的 `ListenAndServe` 方法：

```go
func ListenAndServe(addr string, handler Handler) error {
    server := &Server{Addr: addr, Handler: handler}
    return server.ListenAndServe()
}
```
而我们 HTTP 服务器中第二行代码：

```go
http.ListenAndServe(":8080", nil)
```

创建了一个 `server` 对象，并调用 `server` 对象的 `ListenAndServe` 方法，这里没有直接传递 `Handler`，而是默认
使用 `DefautServeMux` 作为 `multiplexer`。

`Server` 的 `ListenAndServe` 方法中，会初始化监听地址 `Addr`，同时调用 `Listen` 方法设置监听。

```go
for {
    rw, e := l.Accept()
    if e != nil {
        select {
        case <-srv.getDoneChan():
            return ErrServerClosed
        default:
        }
        if ne, ok := e.(net.Error); ok && ne.Temporary() {
            if tempDelay == 0 {
                tempDelay = 5 * time.Millisecond
            } else {
                tempDelay *= 2
            }
            if max := 1 * time.Second; tempDelay > max {
                tempDelay = max
            }
            srv.logf("http: Accept error: %v; retrying in %v", e, tempDelay)
            time.Sleep(tempDelay)
            continue
        }
        return e
    }
    tempDelay = 0
    c := srv.newConn(rw)
    c.setState(c.rwc, StateNew) // before Serve can return
    go c.serve(ctx)
}
```
监听开启之后，一旦客户端请求过来，Go 就开启一个协程 `go c.serve(ctx)` 处理请求，主要逻辑都在 `serve` 方法之中。

`func (c *conn) serve(ctx context.Context)`，这个方法很长，里面主要的一句：`serverHandler{c.server}.ServeHTTP(w, w.req)`。
其中 `w` 由 `w, err := c.readRequest(ctx)` 得到，因为有传递 `context`。

还是来看源代码：

```go
// serverHandler delegates to either the server's Handler or
// DefaultServeMux and also handles "OPTIONS *" requests.
type serverHandler struct {
	srv *Server
}


func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
    // 此 handler 即为 http.ListenAndServe 中的第二个参数
    handler := sh.srv.Handler 
    if handler == nil {
        // 如果 handler 为空则使用内部的 DefaultServeMux 进行处理
        handler = DefaultServeMux
    }
    if req.RequestURI == "*" && req.Method == "OPTIONS" {
        handler = globalOptionsHandler{}
    }
    // 这里就开始处理 http 请求
    // 如果需要使用自定义的 mux，就需要实现 ServeHTTP 方法，即实现 Handler 接口。
    // ServeHTTP(rw, req) 默认情况下是 func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request)
    handler.ServeHTTP(rw, req)
}
```
从 `http.ListenAndServe(":8080", nil)` 开始，`handler` 是 `nil`，所以最后实际 `ServeHTTP` 方法
是 `DefaultServeMux.ServeHTTP(rw, req)`。

```go
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(StatusBadRequest)
		return
	}
	h, _ := mux.Handler(r) // 会匹配路由，h 就是 MyHandler
	h.ServeHTTP(w, r) // 调用自己
}

func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) {

	// CONNECT requests are not canonicalized.
	if r.Method == "CONNECT" {
		// If r.URL.Path is /tree and its handler is not registered,
		// the /tree -> /tree/ redirect applies to CONNECT requests
		// but the path canonicalization does not.
		if u, ok := mux.redirectToPathSlash(r.URL.Host, r.URL.Path, r.URL); ok {
			return RedirectHandler(u.String(), StatusMovedPermanently), u.Path
		}

		return mux.handler(r.Host, r.URL.Path)
	}

	// All other requests have any port stripped and path cleaned
	// before passing to mux.handler.
	host := stripHostPort(r.Host)
	path := cleanPath(r.URL.Path)

	// If the given path is /tree and its handler is not registered,
	// redirect for /tree/.
	if u, ok := mux.redirectToPathSlash(host, path, r.URL); ok {
		return RedirectHandler(u.String(), StatusMovedPermanently), u.Path
	}

	if path != r.URL.Path {
		_, pattern = mux.handler(host, path)
		url := *r.URL
		url.Path = path
		return RedirectHandler(url.String(), StatusMovedPermanently), pattern
	}

	return mux.handler(host, r.URL.Path)
}

// handler is the main implementation of Handler.
// The path is known to be in canonical form, except for CONNECT methods.
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

通过 `func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string)`，我们得到 `Handler h`，然后执
行 `h.ServeHTTP(w, r)` 方法，也就是执行我们的 `MyHandler` 函数（别忘了 `MyHandler` 是HandlerFunc类型，而他的 `ServeHTTP(w, r)` 
方法这里其实就是自己调用自己），把 `response` 写到 `http.ResponseWriter` 对象返回给客户端，`fmt.Fprintf(w, "hello")`，我们在客
户端会接收到 "hello" 。至此整个 HTTP 服务执行完成。


总结下，HTTP 服务整个过程大概是这样：
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

上面代码实际上就是 `server.ListenAndServe()` 执行的实际效果，只不过简单声明了一个结构体 `Server{Addr: addr, Handler: handler}` 实例。
如果我们声明一个 `Server` 实例，完全可以达到深度自定义 `http.Server` 的目的：


```go
package main

import (
	"fmt"
	"net/http"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "hi")
}

func main() {
	// 更多http.Server的字段可以根据情况初始化
	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  0,
		WriteTimeout: 0,
	}
	http.HandleFunc("/", MyHandler)
	_ = server.ListenAndServe()
}
```
我们完全可以根据情况来自定义我们的 `Server`。

还可以指定 `Servemux` 的用法:

```go
package main

import (
	"fmt"
	"net/http"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "hi")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", MyHandler)
	_ = http.ListenAndServe(":8080", mux)
}
```

如果既指定 `Servemux` 又自定义 `http.Server`，因为 `Server` 中有字段 `Handler`，所以我们可以直接把 `Servemux` 变量作
为 `Server.Handler`：


```go
package main

import (
	"fmt"
	"net/http"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {
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

	mux.HandleFunc("/", MyHandler)
	_ = server.ListenAndServe()
}
```

## 自定义处理器

自定义的 `Handler`：

标准库 http 提供了 `Handler` 接口，用于开发者实现自己的 `handler`。只要实现接口的 `ServeHTTP` 方法即可。

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
我们知道，`NewServeMux` 可以创建一个 `ServeMux` 实例，`ServeMux` 同时也实现了 `ServeHTTP` 方法，因此代码中的 `mux` 也是
一种 `handler`。把它当成参数传给 `http.ListenAndServe` 方法，后者会把 `mux` 传给 `Server` 实例。因为指定了 `handler`，
因此整个 `http` 服务就不再是 `DefaultServeMux`，而是 `mux`，无论是在注册路由还是提供请求服务的时候。

任何有 `func(http.ResponseWriter，*http.Request)` 签名的函数都能转化为一个 `HandlerFunc` 类型。这很有用，因为 `HandlerFunc` 对象
内置了 `ServeHTTP` 方法，后者可以聪明又方便的调用我们最初提供的函数内容。

## 中间件 Middleware

所谓中间件，就是连接上下级不同功能的函数或者软件，通常进行一些包裹函数的行为，为被包裹函数提供添加一些功能或行为。前文的 `HandleFunc` 就
能把签名为 `func(w http.ResponseWriter, r *http.Reqeust)` 的函数包裹成 `handler`。这个函数也算是中间件。

Go 的 HTTP 中间件很简单，只要实现一个函数签名为 `func(http.Handler) http.Handler` 的函数即可。`http.Handler` 是一个接口，
接口方法我们熟悉的为 `serveHTTP`。返回也是一个 `handler`。因为 Go 中的函数也可以当成变量传递或者或者返回，因此也可以在中间件函数
中传递定义好的函数，只要这个函数是一个 `handler` 即可，即实现或者被 `handlerFunc` 包裹成为 `handler` 处理器。

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
        // 执行 handler 之前的逻辑
        next.ServeHTTP(w, r)
        // 执行完毕 handler 后的逻辑
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