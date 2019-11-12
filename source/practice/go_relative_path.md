# Go 的相对路径
在构建 Go 项目时，有没有碰到 `go build` 编译好的二进制文件（或者 `go run main.go`），在不同的目录下执行，得到的结果却不一样？

例如，我的目录结构是下面这样的：
```
backend
├── app
│   ├── cmd
│   │   └── cmd.go
│   ├── conf
│   │   └── conf.yaml
│   ├── config
│   │   └── config.go
│   ├── dao
│   │   └── dao.go
│   ├── http
│   │   └── http.go
│   └── main.go
├── go.mod
├── go.sum
└── suiteinstaller
```

`suiteinstaller` 是构建好的二进制文件，在 `backend` 目录下运行或者执行 `go run ./app/main.go` 可以正常运行。但是如果在 `app`
目录下执行同样的命令则会报错：
```sh
# go run ./app/main.go
Fail to parse 'app/conf/conf.yaml': open app/conf/conf.yaml: no such file or directory
exit status 1
```

这时因为 **`Golang` 的相对路径是相对于执行命令时的目录**。而且在代码中使用 `os.Getwd()` 相对路径也是相对于执行命令时的目录。

## 解决
### go build
```go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	path := GetAppPath()
	fmt.Println("--------------------", path)
}
func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}
```

上面的代码中的 `GetAppPath` 方法可以获取**二进制文件所在的路径**。比如;
```sh
# backend 目录下
$ ./app/main
-------------------- /root/code/newui/backend/app

# backend/app 目录下
$ ./main
-------------------- /root/code/newui/backend/app
```

可以看到得到路径是一致的，并不会因为执行命令的路径改变而改变。

### go run
上面的解决方案，对于 `go build` 生成的二进制文件是没问题的，但是如果运行 `go run main.go` 就不行了。
```sh
$ go run ./app/main.go
-------------------- /tmp/go-build424563838/b001/exe
```

这时因为 `go run`  执行时会将文件放到一个临时目录 `/tmp/go-build...` 目录下，编译并运行。

对于 `go run` 可以通过传参，或者环境变量来指定项目路径，再进行拼接。

可以参考 [beego](https://github.com/astaxie/beego/blob/master/config.go#L133-L146) 读取配置文件的代码，
可以兼容 `go build` 和在项目根目录执行 `go run` ，但是若跨目录执行 `go run` 就不行。
