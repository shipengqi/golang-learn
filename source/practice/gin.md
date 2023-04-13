# Gin

## Gin 打包 Angular (React/Vue) 项目

Gin 作为 Web 框架提供 API 接口非常方便，但是在同一个项目中，既提供 API 接口，又要作为前端网页的静态服务器，就比较麻烦。通常 Angular (React/Vue) 
项目需要在 Nginx 或者 Tomcat 转发才可以。有些小项目并不需要前后端分离，如何解决？


### 1.16 版本的 embed

Go 的 1.16 版本增加了 `embed` 的标签，可以利用这个标签将静态资源打包到二进制文件中。

```bash
.
├── config
├── controller
├── model
├── options
├── pkg
│         └── response
│             └── response.go
├── resources
│   ├── dist
│   └── html.go
├── html.go
├── resource.go
├── router.go
├── server.go
└── store
    ├── audited.go
    ├── groups.go
    ├── mysql.go
    ├── settings.go
    ├── store.go
    └── tokens.go
```

上面项目的目录结构中注意这几个文件：

```bash
├── resources
│   ├── dist
│   └── html.go
├── html.go
├── resource.go
├── router.go
```

`dist` 是打包好的静态资源。

`html.go` 为了后面渲染 `index.html` 和静态资源提供的变量：

```go
package resources

import "embed"

//go:embed dist/stat-web/index.html
var Html []byte

//go:embed dist/stat-web
var Static embed.FS
```

`resource.go` 实现了 `FS` 接口：

`FS` 接口：
```go
type FS interface {
  // Open opens the named file.
  //
  // When Open returns an error, it should be of type *PathError
  // with the Op field set to "open", the Path field set to name,
  // and the Err field describing the problem.
  //
  // Open should reject attempts to open names that do not satisfy
  // ValidPath(name), returning a *PathError with Err set to
  // ErrInvalid or ErrNotExist.
  Open(name string) (File, error)
}
```

`resource.go`：

```go
package apiserver

import (
	"embed"
	"io/fs"
	"path"

	"project/resources"
)

type Resource struct {
	fs   embed.FS
	path string
}

func NewResource(staticPath string) *Resource {
	return &Resource{
		fs:   resources.Static, // resources/html.go 中定义的 Static
		path: staticPath,
	}
}

func (r *Resource) Open(name string) (fs.File, error) {
	// rewrite the static files path
	fullName := path.Join(r.path, name) // 这里拼出静态资源的完整路径，注意 windows 下使用 filepath.Join，会导致找不到文件
	return r.fs.Open(fullName)
}
```

`html.go` 中实现了 `HtmlHandler` 用来渲染 `index.html`：

```go
package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"project/resources"
)

type HtmlHandler struct{}

func NewHtmlHandler() *HtmlHandler {
	return &HtmlHandler{}
}

// RedirectIndex 重定向
func (h *HtmlHandler) RedirectIndex(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
	return
}

func (h *HtmlHandler) Index(c *gin.Context) {
	c.Header("content-type", "text/html;charset=utf-8")
	c.String(200, string(resources.Html))
	return
}
```

`router.go` 中配置路由：

```go
func installController(g *gin.Engine) {

    html := NewHtmlHandler()
    g.GET("/", html.Index)
    g.StaticFS("/static", http.FS(NewResource("dist/stat-web")))
    g.StaticFS("/assets", http.FS(NewResource("dist/stat-web/assets")))
    g.NoRoute(html.RedirectIndex)

    // API 接口
	v1 := g.Group("/api/v1")
    {
	   // ...
    }
}
```

上面的路由 `g.StaticFS("/static", http.FS(NewResource("dist/stat-web")))` ，路径之所以是 `/static` 是因为在打包 Angular 项目时使用了 `--deploy-url`：

```bash
ng build <project> --configuration production --deploy-url /static/
```

`--deploy-url` 将被弃用，之后需要考虑其他方式。暂时不适用 `--base-href` 是因为 `--base-href` 会将整个项目的路径全部加上 `--base-href` 指定的路径，例如

`--base-href /stat/` 访问路径是 `localhost:8080/stat`，这在配置子项目时比较有用。但是这里作为一个独立的项目，就不合适，需要给所有的 API 都加上 `stat`。
