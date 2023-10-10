---
title: 常见错误
weight: 9
---

# 常见错误

## go mod tidy error message: "but go 1.16 would select"

```bash
$ go mod tidy
github.com/shipengqi/crtctl/internal/secret-writer imports
	github.com/shipengqi/kube imports
	k8s.io/client-go/kubernetes imports
	k8s.io/client-go/kubernetes/typed/admissionregistration/v1 imports
	k8s.io/client-go/applyconfigurations/admissionregistration/v1 imports
	k8s.io/apimachinery/pkg/util/managedfields imports
	k8s.io/kube-openapi/pkg/util/proto tested by
	k8s.io/kube-openapi/pkg/util/proto.test imports
	github.com/onsi/ginkgo imports
	github.com/onsi/ginkgo/internal/remote imports
	github.com/nxadm/tail imports
	github.com/nxadm/tail/winfile loaded from github.com/nxadm/tail@v1.4.4,
	but go 1.16 would select v1.4.8

To upgrade to the versions selected by go 1.16:
	go mod tidy -go=1.16 && go mod tidy -go=1.17
If reproducibility with go 1.16 is not needed:
	go mod tidy -compat=1.17
For other options, see:
	https://golang.org/doc/modules/pruning
```

[Go 1.17 release notes](https://go.dev/doc/go1.17#go-command):

> By default, go mod tidy verifies that the selected versions of dependencies relevant to the main module are the same
> versions that would be used by the prior Go release (Go 1.16 for a module that specifies go 1.17) 

错误信息提示了两种修复方法。

`go mod tidy -go=1.16 && go mod tidy -go=1.17` - 这将选择依赖版本为 Go 1.16，然后为 Go 1.17。

`go mod tidy -compat=1.17` - 这只是删除了 Go 1.16 的校验和（因此提示 "不需要用 Go 1.16 进行重现"）。

升级到 Go 1.18 后，这个错误应该不会再出现，因为那时模块图的加载与 Go 1.17 相同。
