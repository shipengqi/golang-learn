---
title: Json Unmarshal
---

# Json Unmarshal

Go 官方的 `encoding/json` 包可以实现 json 互转。如：
```go
type AuthResponse struct {
	Auth  AuthData `json:"auth"`
}

type AuthData struct {
	ClientToken  string `json:"client_token"`
}

var authResponse AuthResponse

err = json.Unmarshal(response.Body(), &authResponse)
```

`encoding/json` 包的 `json.Marshal/Unmarshal` 是非常慢的，因为是通过大量反射来实现的。

可以使用第三方库来替代标准库：
- [json-iterator/go](https://github.com/json-iterator/go)，完全兼容标准库，性能有很大提升。
- [ffjson](https://github.com/pquerna/ffjson)