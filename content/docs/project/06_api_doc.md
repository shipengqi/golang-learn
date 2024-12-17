---
title: API 文档
weight: 6
---

## 使用 Swagger 生成 API 文档

Swagger 是基于 OpenAPI 规范的 API 文档工具。

> OpenAPI 是一个 API 规范，它的前身就是 Swagger 规范，目前最新的 OpenAPI
> 规范是 [OpenAPI 3.0](https://swagger.io/docs/specification/about/)（也就是 Swagger 2.0 规范）。

## Swagger 编辑器

Swagger 编辑是一个在线的 API 文档[编辑器](https://editor.swagger.io)，可以在其中编写 OpenAPI 规范，并实时预览 API 文档。

## 基于代码自动生成 Swagger 文档

Go 生成 Swagger 文档常用的工具有两个，分别是 [swag](https://github.com/swaggo/swag) 和 [go-swagger](https://github.com/go-swagger/go-swagger)。

推荐使用 go-swagger：

1. go-swagger 提供了更灵活、更多的功能来描述 API，可以生成客户端和服务器端代码。
2. 使用 swag 的话，每一个 API 都需要有一个冗长的注释，有时候代码注释比代码还要长，但是通过 go-swagger 可以将代码和注释分开编写，可以使代码保持简洁，清晰易读，而且可以把 API 定义放在一个目录中，方便管理。

### 安装 go-swagger

```
$ go get -u github.com/go-swagger/go-swagger/cmd/swagger
$ swagger version
version: v0.30.3
commit: ecf6f05b6ecc1b1725c8569534f133fa27e9de6b
```

命令格式为 `swagger [OPTIONS] <command>`。

`swagger` 提供的子命令：

| 子命令      | 描述                         |
|----------|----------------------------|
| diff     | 对比两个 swagger 文档的差异         |
| expand   | 展开 swagger 定义文档中的 $ref     |
| flatten  | 展平 swagger 文档              |
| generate | 生成 swagger 文档，客户端，服务端代码    |
| ini      | 初始化一个 swagger 定义文档         |
| mix      | 合并 swagger 文档              |
| serv     | 启动 http 服务，用来查看 swagger 文档 |
| validate | 验证 swagger 第一文件是否正确        |

### 使用

go-swagger 通过解析源码中的注释来生成 Swagger 文档。

注释语法：

| 注释语法               | 描述            |
|--------------------|---------------|
| swagger:meta       | 定义全局基本信息      |
| swagger:route      | 定义路由信息        |
| swagger:parameters | API 请求参数      |
| swagger:response   | API 响应参数      |
| swagger:model      | 可以复用的 Go 数据结构 |
| swagger:allOf      | 嵌入其他 Go 结构体   |
| swagger:strfmt     | 格式化的字符串       |
| swagger:ignore     | 需要忽略的结构体      |

`swagger generate` 命令会找到 `main` 函数，然后遍历所有源码文件，解析源码中与 Swagger
相关的注释，然后自动生成 `swagger.json/swagger.yaml` 文件。

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

	log.Fatal(r.Run(":8081"))
}

// Create a user.
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

// Get return user details.
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

`swagger/api/user.go`：

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

另外一个 Go 包中编写带 go-swagger 注释的 API 文档。 先创建一个目录 `swagger/docs`。在 `swagger/docs` 创建 `doc.go` 文件，提供基本的 API 信息：

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

注意最后以 `swagger:meta` 注释结束。

编写完 `doc.go` 文件后， 编写 API 的定义文件 `swagger/docs/user.go`：

```go
package docs

import (
  v1 "github.com/shipengqi/idm/api/apiserver/v1"
  metav1 "github.com/shipengqi/idm/api/meta/v1"
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

- `swagger:route`：描述一个 API，格式为 `swagger:route [method] [url path pattern] [?tag1 tag2 tag3] [operation id]`，tag 可以是多个，相同 tag 的 API 在 Swagger 文档中会被分为一组。operation id 会和
  `swagger:parameters` 的定义进行匹配，就是该 API 的请求参数。
  `swagger:route` 下面的一行是该 API 的描述，需要以 `.` 为结尾。`responses:` 定义了 API 的返回参数，例如当 HTTP 状态码是 200 时，返回 `createUserResponse`，`createUserResponse` 会和 `swagger:response` 的定义进行匹配，匹配成功的 `swagger:response` 就是该 API 返回 200 状态码时的返回。
- `swagger:response`：定义了 API 的返回，格式为 `swagger:response [?response name]`.例如 `getUserResponseWrapper` 中有一个 `Body` 字段，其注释为 `// in:body`，说明该参数是在 HTTP Body 中返回。
  `swagger:response` 上的注释是 response 的描述。`api.User` 会被 go-swagger 解析为 Example Value 和 Model，不需要重复编写。
- `swagger:parameters`：定义了 API 的请求参数，格式为 `swagger:parameters [operationid1 operationid2]`
  。例如 `userRequestParamsWrapper`。`userRequestParamsWrapper` 上的注释是请求参数的描述。

进入 `swagger` 目录，执行如下命令，生成 Swagger API 文档：

```
$ swagger generate spec -o swagger.yaml
```

- `-o`：指定要输出的文件名。swagger 会根据文件名后缀 `.yaml` 或者 `.json`，决定生成的文件格式为 YAML 或 JSON。

启动 HTTP 服务：

```
$ swagger serve --no-open -F=redoc --port 36666 swagger.yaml
```

- `–no-open`：`-–no-open` 禁止调用浏览器打开 URL。
- `-F`：指定文档的风格，可选 `swagger` 和 `redoc`。`redoc` 格式更加易读和清晰。
- `–port`：指定启动的 HTTP 服务监听端口。

在浏览器查看 API 文档：

![swagger-redoc](https://raw.githubusercontent.com/shipengqi/illustrations/0dad40be58e3b5959e8a9093a92ce49a9e280c6d/go/swagger-redoc.png)

还可以使用下面的命令将生成的 `swagger.yaml` 转换为 `swagger.json`：

```
$ swagger generate spec -i ./swagger.yaml -o ./swagger.json
```
