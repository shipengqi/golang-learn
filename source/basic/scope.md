---
title: 作用域
---
# 作用域
声明语句的作用域是指源代码中可以有效使用这个名字的范围。
- 局部变量 在函数体内或代码块内声明的变量称之为局部变量，它们的作用域只在代码块内，参数和返回值变量也是局部变量。
- 全局变量 作用域都是全局的（在本包范围内） 在函数体外声明的变量称之为全局变量，**全局变量可以在整个包甚至外部包
（被导出后 首字母大写）使用**。 全局变量可以在任何函数中使用。


Go 的标识符作用域是基于代码块的。代码块就是包裹在一对大括号内部的声明和语句，并且是可嵌套的。代码块内部声明的名字是无法
被外部块访问的。

声明语句作用域范围的大小。
- 内置的类型、函数和常量，比如 `int`、`len` 和 `true` 是全局作用域
- 在函数外部（也就是包级语法域）声明的名字可以在同一个包的任何源文件中访问
- 导入的包，如 `import "packages/test"`，是对应源文件级的作用域，只能在当前的源文件中访问
- 在函数内部声明的名字，只能在函数内部访问

**一个程序可能包含多个同名的声明，只要它们在不同的词法域就可以。内层的词法域会屏蔽外部的声明。** 