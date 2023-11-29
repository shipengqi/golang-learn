---
title: API 文档
weight: 6
---

# API 文档

# 使用 Swagger 生成 API 文档

主要的 Swagger 工具：

- **Swagger 编辑器**：基于浏览器的[编辑器](https://editor.swagger.io)，可以在其中编写 OpenAPI 规范，并实时预览 API 文档。
- **Swagger UI**：将 OpenAPI 规范呈现为交互式 API 文档，并可以在浏览器中尝试 API 调用。
- **Swagger Codegen**：根据 OpenAPI 规范，生成服务器存根和客户端代码库，目前已涵盖了 40 多种语言。

> OpenAPI 是一个 API 规范，它的前身叫 Swagger 规范，目前最新的OpenAPI规范是 [OpenAPI 3.0](https://swagger.io/docs/specification/about/)（也就是 Swagger 2.0 规范）。

## go-swagger

生成 API 文档一般有两种方式：

1. 如果熟悉 Swagger 语法的话，可以直接编写 JSON/YAML 格式的 Swagger 文档。
2. 通过工具 [swag](https://github.com/swaggo/swag) 或者 [go-swagger](https://github.com/go-swagger/go-swagger) 工具来生成。

使用 go-swagger 的原因：

1. go-swagger 提供了更灵活、更多的功能来描述 API。
2. 使用 swag，每一个 API 都需要有一个冗长的注释，有时候代码注释比代码还要长，但是通过 go-swagger 可以将代码和注释分开编写，一方面可以使代码保持简洁，
   清晰易读，另一方面可以在另外一个包中，统一管理这些 Swagger API 文档定义。

### 安装 go-swagger

```
$ go get -u github.com/go-swagger/go-swagger/cmd/swagger
$ swagger version
version: v0.30.3
commit: ecf6f05b6ecc1b1725c8569534f133fa27e9de6b
```

### 常用命令

|  子命令   | 描述  |
|  ----  | ----  |
| diff  | 对比两个 swagger 文档的差异 |
| expand  | 展开 swagger 定义文档中的 $ref |
| flatten  | 展平 swagger 文档 |
| generate  | 生成 swagger 文档，客户端，服务端代码 |
| ini  | 初始化一个 swagger 定义文档 |
| mix  | 合并 swagger 文档 |
| serv  | 启动 http 服务，用来查看 swagger 文档 |
| validate  | 验证 swagger 第一文件是否正确 |

### 使用

go-swagger 通过解析源码中的注释来生成 Swagger 文档，go-swagger 的详细注释语法可参考[官方文档](https://goswagger.io/)。常用的注释语法：

|  注释语法   | 描述  |
|  ----  | ----  |
| swagger:meta  | 定义 API 接口全局基本信息 |
| swagger:route  | 定义路由信息 |
| swagger:parameters  | API 请求参数 |
| swagger:response  | API 响应参数 |
| swagger:model  | 可以服用的 Go 数据结构 |
| swagger:allOf  | 嵌入其他 Go 结构体 |
| swagger:strfmt  | 格式化的字符串 |
| swagger:ignore  | 需要忽略的结构体 |


`swagger generate` 命令会找到 `main` 函数，然后遍历所有源码文件，解析源码中与 Swagger 相关的注释，然后自动生成 `swagger.json/swagger.yaml` 文件。

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/shipengqi/idm/swagger/api"
	// This line is necessary for go-swagger to find your docs!
	_ "github.com/shipengqi/idm/swagger/docs"
)

var users []*api.User

func main() {
	r := gin.Default()
	r.POST("/users", Create)
	r.GET("/users/:name", Get)

	log.Fatal(r.Run(":5555"))
}

// Create create a user in memory.
func Create(c *gin.Context) {
	var user api.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 10001})
		return
	}

	for _, u := range users {
		if u.Name == user.Name {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user %s already exist", user.Name), "code": 10001})
			return
		}
	}

	users = append(users, &user)
	c.JSON(http.StatusOK, user)
}

// Get return the detail information for a user.
func Get(c *gin.Context) {
	username := c.Param("name")
	for _, u := range users {
		if u.Name == username {
			c.JSON(http.StatusOK, u)
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user %s not exist", username), "code": 10002})
}
```

`user.go`：

```go
package api

// User represents body of User request and response.
type User struct {
    // User's name.
    // Required: true
    Name string `json:"name"`

    // User's nickname.
    // Required: true
    Nickname string `json:"nickname"`

    // User's address.
    Address string `json:"address"`

    // User's email.
    Email string `json:"email"`
}
```

`Required: true` 说明字段是必须的。

另外一个Go包中编写带go-swagger注释的API文档。 先创建一个目录 `swagger/docs`。在 `swagger/docs` 创建 `doc.go` 文件，提供基本的 API 信息：

```go
// Package docs awesome.
//
// Documentation of our awesome API.
//
//     Schemes: http, https
//     BasePath: /
//     Version: 0.1.0
//     Host: some-url.com
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - basic
//
//    SecurityDefinitions:
//    basic:
//      type: basic
//
// swagger:meta
package docs
```
[swagger](..%2F..%2Fswagger)
注意最后以 `swagger:meta` 注释结束。[swagger](..%2F..%2Fswagger)

编写完 `doc.go` 文件后，进入 `swagger` 目录，执行如下命令，生成 Swagger API 文档，并启动 HTTP 服务，在浏览器查看 Swagger：

```
$ swagger generate spec -o swagger.yaml
$ swagger serve --no-open -F=redoc --port 36666 swagger.yaml
```

- `-o`：指定要输出的文件名。swagger 会根据文件名后缀 `.yaml` 或者 `.json`，决定生成的文件格式为 YAML 或 JSON。
- `–no-open`：`-–no-open` 禁止调用浏览器打开 URL。
- `-F`：指定文档的风格，可选 swagger 和 redoc。redoc 格式更加易读和清晰。
- `–port`：指定启动的 HTTP 服务监听端口。

生成的 `swagger.yaml` 转换为 `swagger.json`：
```
$ swagger generate spec -i ./swagger.yaml -o ./swagger.json
```

编写 API 接口的定义文件 `swagger/docs/user.go`：

```go
package docs

import (
	v1 "github.com/shipengqi/idm/api/apiserver/v1"
	metav1 "github.com/shipengqi/idm/api/meta/v1"
	"github.com/shipengqi/idm/internal/apiserver/contorller/v1/user"
)

// swagger:route GET /users/{name} Users getUserRequest
//
// Get details for specified user.
//
// Get details for specified user according to input parameters.
//
//     Responses:
//       default: errResponse
//       200: getUserResponse

// swagger:route GET /users Users listUserRequest
//
// List users.
//
// List users.
//
//     Responses:
//       default: errResponse
//       200: listUserResponse

// List users request.
// swagger:parameters listUserRequest
type listUserRequestParamsWrapper struct {
	// in:query
	metav1.ListOptions
}

// List users response.
// swagger:response listUserResponse
type listUserResponseWrapper struct {
	// in:body
	Body v1.UserList
}

// User response.
// swagger:response getUserResponse
type getUserResponseWrapper struct {
	// in:body
	Body v1.User
}

// swagger:parameters createUserRequest updateUserRequest
type userRequestParamsWrapper struct {
	// User information.
	// in:body
	Body v1.User
}

// swagger:parameters deleteUserRequest getUserRequest updateUserRequest
type userNameParamsWrapper struct {
	// Username.
	// in:path
	Name string `json:"name"`
}

// ErrResponse defines the return messages when an error occurred.
// swagger:response errResponse
type errResponseWrapper struct {
	// in:body
	Body response.Response
}

// Return nil json object.
// swagger:response okResponse
type okResponseWrapper struct{}
```

- `swagger:route`：代表一个 API 接口描述的开始，后面的字符串格式为 HTTP 方法，URL，Tag，ID。可以填写多个 tag，相同 tag 的 API 接口在 Swagger 文档中会被分为一组。ID 是一个标识符，
  `swagger:parameters` 是具有相同 ID 的 `swagger:route` 的请求参数。`swagger:route` 下面的一行是该 API 接口的描述，需要以英文点号为结尾。`responses:` 定义了API接口的返回参数，
  例如当 HTTP 状态码是 200 时，返回 `createUserResponse`，`createUserResponse` 会跟 `swagger:response` 进行匹配，匹配成功的 `swagger:response` 就是该API接口返回200状态码时的返回。
- `swagger:response`：定义了 API 接口的返回，例如 `getUserResponseWrapper`，关于名字，我们可以根据需要自由命名，并不会带来任何不同。`getUserResponseWrapper` 中有一个 `Body` 字段，
  其注释为 `// in:body`，说明该参数是在 HTTP Body 中返回。`swagger:response` 之上的注释会被解析为返回参数的描述。`api.User` 自动被 go-swagger 解析为 Example Value 和 Model。我们不用
  再去编写重复的返回字段，只需要引用已有的 Go 结构体即可。
- `swagger:parameters`：定义了 API 接口的请求参数，例如 `userParamsWrapper`。`userParamsWrapper` 之上的注释会被解析为请求参数的描述，`// in:body` 代表该参数是位于 HTTP Body 中。
  同样，`userParamsWrapper` 结构体名我们也可以随意命名。`swagger:parameters` 之后的 `createUserReques` t会跟 `swagger:route` 的 ID 进行匹配，匹配成功则说明是该 ID 所在 API 接口的请求参数。



## 使用 swag

Download swag:
```bash
go get -u github.com/swaggo/swag/cmd/swag

$ swag -v
swag version v1.6.5
```

Install gin-swagger:
```bash
go get -u github.com/swaggo/gin-swagger@v1.2.0 
go get -u github.com/swaggo/files
go get -u github.com/alecthomas/template
```

```go
package router

import (
    "github.com/gin-gonic/gin"
    _ "github.com/shipengqi/idm/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
    "github.com/swaggo/gin-swagger/swaggerFiles"
)

func Init() *gin.Engine {
	r := gin.New()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
```

Run `swag init` in the project's root folder which contains the `main.go` file.
This will parse your comments and generate the required files (`docs` folder and `docs/docs.go`).

Generate:
```bash
docs/
├── docs.go
└── swagger
    ├── swagger.json
    └── swagger.yaml
```

Run your app, and browse to `http://localhost:<port>/swagger/index.html`. 