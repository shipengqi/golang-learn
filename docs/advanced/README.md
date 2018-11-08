# 高级编程
## WEB 编程
### 搭建一个web
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

#### http包
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

#### 自己实现一个简易的路由器
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

### 表单
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

#### 预防跨站脚本
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

#### 文件上传
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

### 访问数据库
Go定义了`database/sql`接口。