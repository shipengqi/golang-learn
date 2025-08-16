---
title: Casbin
weight: 16
---

[Casbin](https://github.com/casbin/casbin) 是基于 Go 语言的开源权限控制库。支持 ACL，RBAC，ABAC 等常用的访问控制模型。

- Casbin **只负责访问控制**，不负责验证用户的用户名、密码，应该有专门的组件负责身份认证，再配合 Casbin 进行访问控制。
- Casbin **只存储用户和角色之间的映射关系**，不存储用户、角色等信息。

Casbin 的访问控制模型核心叫做 PERM（Policy、Effect、Request、Matcher） 元模型。

## PERM 元模型

PERM（Policy、Effect、Request、Matcher）可以简单理解为，当一个请求（Request）进来，需要通过策略匹配器（Matcher）去匹配存储在数据库中或者
csv 文件中的策略规则，拿到所有匹配的策略规则的结果（`eft`）之后，在使用 `Effect` 定义中的表达式进行计算，最终返回一个 `true`
（通过）或 `false`（拒绝）。

### Request

Request 代表请求。

请求的定义，这个定义的是 `e.Enforce(...)` 函数中的参数：

```
[request_definition]        
r = sub, obj, act
```

- `sub`：代表访问资源的实体，一般就是指用户。
- `obj`：要访问的资源。
- `act`：对资源执行的操作。

这里是定义你的请求格式。如果你不需要指定资源，你可以定义：

```
[request_definition]        
r = sub, act
```

或者如果你有两个访问实体，你可以定义：

```
[request_definition]        
r = sub, sub2, obj, act
```

### Policy

Policy 代表策略。

**定义访问策略模型**，这个定义的是数据库中或者 csv 文件中策略规则的格式：

```
[policy_definition]
p = sub, obj, act, eft
p2 = sub, act
```

- `eft`：表示策略的结果，一般为空，默认是 `allow`。`eft` **只接受两个值是 `allow` 或 `deny`**。

**策略规则**，策略规则一般存储在数据库中，也可以在 csv 文件中：

```
p, alice, data1, read
p2, bob, write-all-objects
```

策略中的每一行都被称为**策略规则**。 每个策略规则都以策略类型（Policy type）开始，如 `p` 或 `p2`。用于匹配策略定义中的 `p`、
`p2`。匹配器中使用时：

```
e = some(where (p2.eft == allow))
```

### Matcher

Matcher 代表策略匹配器。用来验证一个请求是否匹配某个策略规则。

定义匹配规则：

```
[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
```

上面定义的匹配规则表示一个请求如果满足下面的三个条件，匹配器就会认为匹配到了该策略规则，并返回策略结果 `p.eft`：

- `r.sub == p.sub`：请求的 `sub` 和策略的 `sub` 相等
- `r.obj == p.obj`：请求的 `obj` 和策略的 `obj` 相等
- `r.act == p.act`：请求的 `act` 和策略的 `act` 相等

> 匹配器中可以使用算术运算符如 `+`、`-`、`*`、`/` 和逻辑运算符如 `&&`、`||`、`!`。

### Effect

Effect 代表策略效果。它可以被理解为一种模型，在这种模型中，对匹配结果再次作出逻辑组合判断。

例如一个请求匹配到了两条策略规则，一条规则允许，另一条规则拒绝。

策略效果的定义为：

```
[policy_effect]
e = some(where (p.eft == allow))
```

这意味着，如果匹配的策略规则有一条的策略结果为 `allow`，那么最终结果为 `allow`。

策略效果的另一个例子：

```
[policy_effect]
e = !some(where (p.eft == deny))
```

这意味着，如果匹配的策略规则有一条的策略结果为 `deny`，那么最终结果为 `deny`。否则就是 `allow`。

策略效果甚至可以用逻辑表达式连接：

```
[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))
```

这意味着必须至少有一个匹配的 `allow` 策略规则，并且不能有任何匹配的 `deny` 策略规则。

目前 Casbin **只支持下面的内置策略效果**：

- `some(where (p.eft == allow))`
- `!some(where (p.eft == deny))`
- `some(where (p.eft == allow)) && !some(where (p.eft == deny))`
- `priority(p.eft) || deny`
- `subjectPriority(p.eft)`

## ACL Model

Casbin 提供了一个在线的[编辑器](https://casbin.org/editor/)，可以用来验证你的 PERM 模型和策略规则。

一个简单的 ACL Model：

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
```

策略规则（Policy）：

```
p, alice, data1, read
p, bob, data2, write
p, bob, data1, write, deny
```

请求（Request）：

```
alice, data1, read
```

上面的请求，可以匹配到策略 `p, alice, data1, read`，并且默认的 `eft` 是 `allow`，经过 Effect 表达式
`some(where (p.eft == allow))`，得到执行结果：

```
true Reason: [alice, data1, read]
```

接着修改请求：

```
alice, data1, write
bob, data1, read
```

上面的两个请求得到执行结果：

```
false
false
```

修改请求：

```
bob, data1, write
```

会匹配到策略 `p, bob, data1, write` 得到结果 `deny`，所以执行结果是 `false`。

如果策略中有两条相同的策略，但是 `eft` 不同，例如：

```
p, bob, data1, write, allow
p, bob, data1, write, deny
```

请求 `bob, data1, write` 得到的执行结果是 `true Reason: [bob, data1, write]`，因为 Effect 的表达式只要有一条 `allow` 就返回
`true`。

如果想要它不通过就修改 Effect 为：

```
[policy_effect]
some(where (p.eft == allow)) && !some(where (p.eft == deny))
```

表示必须至少有一个匹配的 `allow` 策略规则，并且不能有任何匹配的 `deny` 策略规则。

## RBAC Model

RBAC 有一个新的概念角色域 `role_definition`，`[role_definition]` 用于表示 `sub` 与角色的关系，角色的继承关系。

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

- `g = _, _` 中的第一个参数 `_` 表示**用户**，第二个 `_` 表示**角色**。
- `g(r.sub, p.sub)` 表示请求中的 `sub` 是否与策略规则中的（`p.sub`）存在角色关系。

> `g` 函数：是 Casbin 中用于检查角色关系的函数。例如 alice 是 admin，策略规则中定义了 `p.sub = admin`，那么
`g(alice, admin)` 返回 `true`。

策略文件：

```
// 角色 admin 可以对 data1 资源执行 read 和 write 操作
p, admin, data1, read
p, admin, data1, write
// 角色 user 可以对 data2 资源执行 read 操作
p, user, data2, read

// alice 被分配为 admin 角色，bob 被分配为 user 角色。
g, alice, admin
g, bob, user
```

请求：

```
alice, data1, write
bob, data1, read
```

上面的两个请求得到执行结果：

```
true Reason: [admin, data1, write]
false
```

- alice 是 admin 角色，可以对 data1 资源执行 read 和 write 操作，所以第一个请求结果为 `true`。
- bob 是 user 角色，只有 data2 资源的 read 权限，没有 data1 的任何权限，所以请求结果为 `false`。

> 请求 `alice, data1, write` 经过 `g(alice, admin)` 处理，用户 alice 和角色 admin 都会用来去匹配策略规则中的 `sub`。

### 角色层次

Casbin 的 RBAC 支持 RBAC1 的角色层次结构特性，这意味着如果 alice 用有角色 role1，并且角色 role1 拥有角色 role2，那么 alice
也将拥有 role2 的所有权限。

Casbin 中的内置角色管理器，可以指定**最大层次级别**。默认值是 `10`。

```go
// NewRoleManager is the constructor for creating an instance of the
// default RoleManager implementation.
func NewRoleManager(maxHierarchyLevel int) rbac.RoleManager {
    rm := RoleManager{}
    rm.allRoles = &sync.Map{}
    rm.maxHierarchyLevel = maxHierarchyLevel
    rm.hasPattern = false

    return &rm
}
```

### 区分角色和用户

在 RBAC 系统内**用户和角色不能使用相同的名称**，因为对于 Casbin 来说，不管是用户还是角色都只是一个**字符串**，Casbin
是没有办法知道你指定的是用户 alice 还是角色 alice。可以通过使用前缀来解决这个问题，例如 `role::alice` 表示角色 alice。

### 查询隐式角色或权限

当**用户通过 RBAC 层次结构继承角色或权限，而不是在策略规则中直接分配它们时**，这种类型的分配为**隐式**。 要查询此类隐式关系，需要使用这两个API：`GetImplicitRolesForUser()` 和 `GetImplicitPermissionsForUser()`，而不是 `GetRolesForUser()` 和 `GetPermissionsForUser()`。

### 菜单权限

RBAC 通常只需要用户的角色，只使用 `g` 就可以了。但是当你需要为资源定义继承关系时，可以同时使用 `g` 和 `g2`，下面的例子是定义菜单权限的模型：

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _
g2 = _, _

[policy_effect]
e = some(where (p_eft == allow)) && !some(where (p_eft == deny))

[matchers]
m = g(r.sub, p.sub) && r.act == p.act && (g2(r.obj, p.obj) || r.obj == p.obj)
```

`g2` 是为了定义菜单项的层次结构。

策略规则：

```
p, ROLE_ROOT, SystemMenu, read, allow
p, ROLE_ROOT, AdminMenu, read, allow
p, ROLE_ROOT, UserMenu, read, deny
p, ROLE_ADMIN, UserMenu, read, allow
p, ROLE_ADMIN, AdminMenu, read, allow
p, ROLE_ADMIN, AdminSubMenu_deny, read, deny
p, ROLE_USER, UserSubMenu_allow, read, allow

g, user, ROLE_USER
g, admin, ROLE_ADMIN
g, root, ROLE_ROOT
g, ROLE_ADMIN, ROLE_USER

g2, UserSubMenu_allow, UserMenu
g2, UserSubMenu_deny, UserMenu
g2, UserSubSubMenu, UserSubMenu_allow
g2, AdminSubMenu_allow, AdminMenu
g2, AdminSubMenu_deny, AdminMenu
g2, (NULL), SystemMenu
```

- `g2, UserSubMenu_allow, UserMenu` 表示 `UserSubMenu_allow` 是 `UserMenu` 的子菜单。
- `g2, (NULL), SystemMenu` 表示 `SystemMenu` 没有子菜单项，意味着它是顶级菜单项。

| 菜单名称 | ROLE_ROOT | ROLE_ADMIN | ROLE_USER |
| --- | --- |------------| --- |
| SystemMenu | ✅ | ❌          | ❌ |
| UserMenu | ❌ | ✅ | ❌ |
| UserSubMenu_allow | ❌ | ✅ | ✅ |
| UserSubSubMenu | ❌ | ✅ | ✅ |
| UserSubMenu_deny | ❌ | ✅ | ❌ |
| AdminMenu | ✅ | ✅ | ❌ |
| AdminSubMenu_allow | ✅ | ✅ | ❌ |
| AdminSubMenu_deny | ✅ | ❌ | ❌ |

### 菜单权限继承的两个重要规则

**父菜单权限的继承**：

如果一个父菜单明确地被授予 `allow` 权限，那么它的所有子菜单也默认为 `allow` 权限，除非特别标记为 `deny`
。这意味着一旦一个父菜单是可访问的，它的子菜单默认也是可访问的。

**没有直接设置权限的父菜单**：

如果一个父菜单没有直接设置权限（既没有明确允许也没有明确拒绝），但**如果有一个子菜单明确被授予 `allow` 权限，那么父菜单被隐式地认为具有 `allow` 权限**。这**确保了用户可以导航到这些子菜单**。

## RBAC with Domains

Casbin 中的 RBAC 是可以支持多租户的。

多租户的 Model 定义：

```
[request_definition]
// 这里的 dom 参数就表示域/租户。
r = sub, dom, obj, act

[policy_definition]
// 这里的 dom 参数就表示域/租户。
p = sub, dom, obj, act

[role_definition]
// 这里的但三个参数 `_` 就表示域/租户。
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
```

对应的策略规则是：

```
// 域 dom/租户 在第二个参数的位置
p, admin, tenant1, data1, read
p, admin, tenant2, data2, read

g, alice, admin, tenant1
g, alice, user, tenant2
```

- tenant1 中的 admin 角色可以读取 data1。
- alice 在 tenant1 中拥有 admin 角色，在 tenant2 中拥有 user 角色。

请求 `alice, tenant1, data1, read` 的执行结果为 `true Reason: [admin, tenant1, data1, read]`，因为 alice 在 tenant1 中拥有
admin 角色。

请求 `alice, tenant2, data2, read` 的执行结果为 `false`，因为 alice 在 tenant2 中不是 admin，所以她无法读取 data2。

匹配器：

```
[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
```

- `g(r.sub, p.sub, r.dom)`：验证请求中的 `r.sub` 是否与策略规则中的 `p.sub` 在特定域（`r.dom`）中存在角色关系。
- `r.dom == p.dom`：确保请求中的域（`r.dom`）与策略规则中的域（`p.dom`）一致。在 RBAC with domains 模型中，*
  *不同的域可能有相同的角色或用户，需要明确区分所属域**。

## ABAC

ABAC (Attribute-Based Access Control) 模型，是基于如用户、资源、环境等属性来控制权限。

**ABAC 模型核心概念**：

- **属性**（Attributes）：描述请求主体（用户）、对象（资源）和环境的相关信息。
    - **用户属性**：如用户名、年龄、部门、角色等。
    - **资源属性**：如资源类型、标签、创建者等。
    - **环境属性**：如时间、地点、设备类型等。

官方的 ABAC 示例：

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == r.obj.Owner
```

在匹配器中，使用 `r.obj.Owner` 代替 `r.obj`。传递给 `Enforce()` 函数的 `r.obj` 将是一个结构体或类实例，而不是一个字符串。
Casbin 将使用反射来检索该结构体或类中的 `obj` 成员变量。

```go
type testResource struct {
Name  string
Owner string
}

e, _ := NewEnforcer("examples/abac_model.conf")
e.EnableAcceptJsonRequest(true)

data1 := testResource{Name: "data1", Owner: "bob"}

ok, _ := e.Enforce("alice", data1, "read")
```

可以使用 `e.EnableAcceptJsonRequest(true)` 启用向 `enforcer` 传递 JSON 参数的功能。但是注意**启用接受 JSON 参数的功能可能会导致性能下降
1.1 到 1.5 倍**。

```go
e, _ := NewEnforcer("examples/abac_model.conf")
e.EnableAcceptJsonRequest(true)

data1Json := `{ "Name": "data1", "Owner": "bob"}`

ok, _ := e.Enforce("alice", data1Json, "read")
```

### 使用 ABAC

要使用 ABAC，需要做两件事：

- 在模型匹配器中指定属性。
- 将结构体或类实例作为参数传递给 `Enforce()` 函数。

> **只有请求元素支持 ABAC（`r.sub`、`r.obj`、`r.act` 等）**，`p.sub` 这样的策略元素是不支持的，因为在策略规则中没有办法定义一个结构体或类，只支持字符串。

### ABAC with Policy

在许多情况下，授权系统需要复杂和大量的 ABAC 规则，所以上面修改匹配器的方式是满足不了需求的。

Casbin 提供了 `eval()` 函数，用于动态执行字符串形式的表达式。这个函数可以实现在策略中添加规则，而不是在模型中。

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub_rule, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = eval(p.sub_rule) && r.obj == p.obj && r.act == p.act
```

策略规则：

```
p, r.sub.Age > 18, /data1, read
p, r.sub.Age < 60, /data2, write
```

请求 `{Age: 30}, /data1, read` 的执行结果为 `true Reason [r.sub.Age > 18, /data1, read]`。 因为匹配器通过
`eval(p.sub_rule)` 执行策略中的条件逻辑 `r.sub.Age > 18`。

## Super Admin

超级管理员是整个系统的管理员。可以用于 RBAC、ABAC 等模型中。

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act || r.sub == "root"
```

`r.sub == "root"` 检查 `sub` 是否为 `root`，如果是，返回 `allow`。

## 函数

Casbin 提供了一系列内置函数，用于**扩展和简化匹配器的逻辑处理**。这些函数可以支持字符串匹配、路径匹配、正则表达式匹配等功能。

`keyMatch/keyMatch2/keyMatch3/keyMatch4/keyMatch5` 都是匹配 URL 路径的，`regexMatch` 使用正则匹配，`ipMatch` 匹配 IP 地址。参见 [Functions](https://casbin.org/zh/docs/function/)。

### keyMatch

用于字符串通配符匹配，`*` 表示通配。当需要对路径规则进行通配时使用，例如 `/data/*` 可以匹配 `/data/1` 和 `/data/2`。

```
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act
```

策略文件：

```
p, alice, /data/*, read
```

请求校验：

```go
e.Enforce("alice", "/data/123", "read")  // 返回 True
e.Enforce("alice", "/data", "read")     // 返回 False
```

### keyMatch2

支持更复杂的路径匹配规则，`:parameter` 表示路径变量。当需要匹配 RESTful 风格的 URL 时使用，例如 `/data/:id` 可以匹配 `/data/1`。

```
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act
```

策略文件：

```
p, alice, /data/:id, read
```

请求校验：

```go
e.Enforce("alice", "/data/123", "read")  // 返回 True
e.Enforce("alice", "/data/abc", "read") // 返回 True
e.Enforce("alice", "/data/", "read")    // 返回 False
```

### 自定义函数

例如定义一个函数，判断两个字符串是否是反向（回文）关系：

```go
package main

import (
	"errors"
	"strings"
)

// 自定义函数：判断两个字符串是否为反向字符串
func ReverseMatch(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New("ReverseMatch requires exactly 2 arguments")
	}

	str1, ok1 := args[0].(string)
	str2, ok2 := args[1].(string)

	if !ok1 || !ok2 {
		return nil, errors.New("arguments must be strings")
	}

	// 判断 str1 是否为 str2 的反向
	return str1 == ReverseString(str2), nil
}

// 辅助函数：反转字符串
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

```

通过 `e.AddFunction` 将自定义函数注册到 Casbin 的模型中：

```go
package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
)

func main() {
	// 加载模型和策略
	e, _ := casbin.NewEnforcer("model.conf", "policy.csv")

	// 注册自定义函数
	e.AddFunction("ReverseMatch", ReverseMatch)

	// 测试权限校验
	testCases := []struct {
		sub string
		obj string
		act string
	}{
		{"alice", "ecila", "read"}, // 匹配反向字符串
		{"bob", "bob", "read"},    // 不匹配
	}

	for _, tc := range testCases {
		result, _ := e.Enforce(tc.sub, tc.obj, tc.act)
		fmt.Printf("Enforce(%s, %s, %s) = %v\n", tc.sub, tc.obj, tc.act, result)
	}
}
```

在模型中使用函数：

```
[matchers]
m = r.sub == p.sub && ReverseMatch(r.obj, p.obj) && r.act == p.act
```


## Enforcer

Casbin [Enforcer](https://casbin.org/zh/docs/enforcers) （执行器）往内存中缓存策略数据时并**不是并发安全**的，
所以 Casbin 还实现了多种 `Enforcer`，可以分为并发安全的 `SyncedEnforcer`，带缓存的 `CachedEnforcer` 等等。

- `Enforcer`：基础的执行器，内存中的操作不是并发安全的。
- `CachedEnforcer`：基于 `Enforcer` 实现，使用内存中的映射来缓存请求的评估结果，这样可以提高访问控制决策的性能。
- `SyncedEnforcer`：基于 `Enforcer` 实现的并发安全的执行器。
- `SyncedCachedEnforcer`：是 `SyncedEnforcer` 和 `CachedEnforcer` 的结合体，支持缓存请求的评估结果，同时是并发安全的。
- `DistributedEnforcer`：基于 `SyncedEnforcer` 实现，支持分布式集群中的多个实例。

## Adapter

[Adapter](https://casbin.org/zh/docs/adapters) 是用来持久化策略数据的（支持文件，数据库等各种持久化方式）。用户可以使用 `LoadPolicy()` 从持久化存储中加载策略规则，
使用 `SavePolicy()` 将策略规则保存到持久化存储中。

1. 从 `Adapter` A 加载策略规则到内存中：

   ```go
   e, _ := NewEnforcer(m, A)
   e.LoadPolicy()
   ```

2. 将适配器从 A 转换到 B
   ```go
   e.SetAdapter(B)
   ```
   
3. 将策略从内存保存到持久化存储中：

   ```go
   e.SavePolicy()
   ```

Casbin 初始化后是可以重新加载模型，重新加载策略或保存策略的：

```go
// Reload the model from the model CONF file.
e.LoadModel()

// Reload the policy from file/database.
e.LoadPolicy()

// Save the current policy (usually after changed with Casbin API) back to file/database.
e.SavePolicy()
```

### AutoSave

由于 Casbin 为了提高性能，在 `Enforcer` 中内置了一个**内存缓存**，所有的操作都是在操作内存中的数据，直到调用 `SavePolicy()`。
`AutoSave` 功能，是将策略的操作自动保存到存储中，这意味着单个策略规则添加，删除或者更新，都会自动保存，而不需要再调用 `SavePolicy()`。

{{< callout type="info" >}}
AutoSave 功能需要适配器的支持，如果适配器支持，可以使用 `EnableAutoSave()` 开启或者禁用。**默认启用**。
{{< /callout >}}

```go

// By default, the AutoSave option is enabled for an enforcer.
a := xormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/")
e := casbin.NewEnforcer("examples/basic_model.conf", a)

// Disable the AutoSave option.
e.EnableAutoSave(false)

// Because AutoSave is disabled, the policy change only affects the policy in Casbin enforcer,
// it doesn't affect the policy in the storage.
e.AddPolicy(...)
e.RemovePolicy(...)

// Enable the AutoSave option.
e.EnableAutoSave(true)

// Because AutoSave is enabled, the policy change not only affects the policy in Casbin enforcer,
// but also affects the policy in the storage.
e.AddPolicy(...)
e.RemovePolicy(...)
```

{{< callout type="warning" >}}
`SavePolicy()` 会删除持久化存储中的所有策略规则，然后将 Enforcer 内存中的策略规则保存到持久化存储中。因此，当策略规则的数量较大时，可能会有性能问题。
{{< /callout >}}

## Watcher

在使用 Casbin 作为权限管理时，一般会增加多个 Casbin 实例来保证请求高并发和稳定性，但是与此同时会出现多个实例之间数据不一致的问题。

这是因为在 Casbin 内置的**内存缓存**，每当执行读数据的动作时，会先获取内存缓存中的数据，而不是数据库中的数据。

因此在多实例场景下，当一个实例的 Casbin 数据发生变更时，其它实例的 Casbin 内存缓存并没有及时同步刷新，这就导致了当请求被分发到对应的实例上时，获取到的依然还是旧的数据。

对此 Casbin 官方提供了一种解决方案 **[Watcher 机制](https://casbin.org/zh/docs/watchers)**，来维持多个 Casbin 实例之间数据的一致性。

## Dispatcher

Casbin 多节点之间的策略数据同步，可以通过 Watcher 机制来实现，也可通过 Dispatcher 调度器来实现。

官方实现的 [hraft-dispatcher](https://github.com/casbin/hraft-dispatcher) 的架构：

![dispatcher-architecture](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/dispatcher-architecture.svg)

在 Dispatcher 中能使用适配器，因为 Dispatcher 自带一个适配器。所有策略都由 Dispatcher 维护。不能调用 `LoadPolicy` 和 `SavePolicy` 方法，
因为这会影响数据的一致性。可以理解为 `Dispatcher = Adapter + Watcher`。

使用示例：

```go
m, err := model.NewModelFromString(modelText)
if err != nil {
    log.Fatal(err)
}

// Adapter is not required here.
// Dispatcher = Adapter + Watcher
e, err := casbin.NewDistributedEnforcer(m)
if err != nil {
    log.Fatal(err)
}

// New a Dispatcher
dispatcher, err := hraftdispatcher.NewHRaftDispatcher(&hraftdispatcher.Config{
    Enforcer:      e,
    JoinAddress:   config.JoinAddress,
    ListenAddress: config.ListenAddress,
    TLSConfig:     tlsConfig,
    DataDir:       config.DataDir,
})
if err != nil {
    log.Fatal(err)
}

e.SetDispatcher(dispatcher)
```

## 性能优化

### 优化匹配器

匹配器是 Casbin 权限检查的核心，优化匹配器可以显著提升性能。

**通过合并或简化条件减少计算开销**：

```
// 优化前
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act && r.dom == p.dom

// 优化后
// 将多个条件合并为字符串比较
m = g(r.sub, p.sub) && r.dom == p.dom && r.obj + ":" + r.act == p.obj + ":" + p.act
```

- **字符串拼接的计算成本通常低于多次逻辑比较**。
- 合并后匹配器计算次数减少。

如果匹配器的条件多且顺序不当，可能会导致不必要的计算。**将匹配概率较高的条件放在前面**：

```
# 优化前
m = r.obj == p.obj && g(r.sub, p.sub) && r.act == p.act

# 优化后（更高概率条件在前）
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

### 尽量减少角色继承层级

角色继承关系会影响匹配器的效率。深度过大的继承链会增加计算成本。

### CachedEnforcer

`CachedEnforcer` 基于 `Enforcer`，但是增加了缓存机制。

1. **缓存评估结果**：`CachedEnforcer` 使用内存中的映射来缓存请求的评估结果，这样可以提高访问控制决策的性能，减少对策略存储的频繁访问。
2. **指定过期时间清除缓存**：它允许在指定的过期时间内清除缓存，这有助于确保策略的变更能够及时生效 。
3. **通过读写锁保证缓存的并发安全**。
4. **启用缓存**：可以使用 `EnableCache` 方法来启用评估结果的缓存，**默认是启用的**。
5. **API 一致性**：`CachedEnforcer` 的其他 API 方法与 `Enforcer` 相同。

`CachedEnforcer` 的部分源码：

```go
// Enforce decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (sub, obj, act).
// if rvals is not string , ignore the cache.
func (e *CachedEnforcer) Enforce(rvals ...interface{}) (bool, error) {
    // 没有开启缓存
    if atomic.LoadInt32(&e.enableCache) == 0 {
        return e.Enforcer.Enforce(rvals...)
    }

	// 查询缓存 key 是否存在
    key, ok := e.getKey(rvals...)
    if !ok {
        return e.Enforcer.Enforce(rvals...)
    }
    
	// 查询缓存的评估结果
    if res, err := e.getCachedResult(key); err == nil {
        return res, nil
    } else if err != cache.ErrNoSuchKey {
        return res, err
    }
    
	// 没有找到缓存的评估结果
    res, err := e.Enforcer.Enforce(rvals...)
    if err != nil {
        return false, err
    }

	// 缓存评估结果
    err = e.setCachedResult(key, res, e.expireTime)
    return res, err
}

// LoadPolicy，策略重新加载，会清除缓存的所有评估结果
func (e *CachedEnforcer) LoadPolicy() error {
	if atomic.LoadInt32(&e.enableCache) != 0 {
		if err := e.cache.Clear(); err != nil {
			return err
		}
	}
	return e.Enforcer.LoadPolicy()
}

// RemovePolicy 策略删除会清除响应的 key 的评估结果
func (e *CachedEnforcer) RemovePolicy(params ...interface{}) (bool, error) {
	if atomic.LoadInt32(&e.enableCache) != 0 {
		key, ok := e.getKey(params...)
		if ok {
			if err := e.cache.Delete(key); err != nil && err != cache.ErrNoSuchKey {
				return false, err
			}
		}
	}
	return e.Enforcer.RemovePolicy(params...)
}

// getCachedResult setCachedResult 对缓存的操作是并发安全的
func (e *CachedEnforcer) getCachedResult(key string) (res bool, err error) {
    e.locker.Lock()
    defer e.locker.Unlock()
    return e.cache.Get(key)
}

func (e *CachedEnforcer) setCachedResult(key string, res bool, extra ...interface{}) error {
    e.locker.Lock()
    defer e.locker.Unlock()
    return e.cache.Set(key, res, extra...)
}
```

{{< callout type="warning" >}}
`CachedEnforcer` 并没有针对 `UpdatePolicy` 做特殊处理，如果**使用 `UpdatePolicy` 修改策略，那么缓存的评估结果不会被清理**。

所以尽量使用 `AddPolicy` 和 `RemovePolicy` 来修改策略。
{{< /callout >}}


### 分片

**进行分片，让 Casbin 执行器只加载一小部分策略规则**。例如，`执行器_0` 可以服务于 `租户_0` 到 `租户_99`，而 `执行器_1` 可以服务于 `租户_100` 到 `租户_199`。

可以利用 [LoadFilteredPolicy ](https://casbin.org/zh/docs/policy-subset-loading) 来实现。

```go
import (
    "github.com/casbin/casbin/v2"
    fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

enforcer, _ := casbin.NewEnforcer()

adapter := fileadapter.NewFilteredAdapter("examples/rbac_with_domains_policy.csv")
enforcer.InitWithAdapter("examples/rbac_with_domains_model.conf", adapter)

filter := &fileadapter.Filter{
    P: []string{"", "domain1"},
    G: []string{"", "", "domain1"},
}
enforcer.LoadFilteredPolicy(filter)

// The loaded policy now only contains the entries pertaining to "domain1".
```