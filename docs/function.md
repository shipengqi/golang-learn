# 函数
`func`关键字声明函数：
```go
func 函数名(形式参数列表) (返回值列表) {
    函数体
}
```
如果函数返回一个无名变量或者没有返回值，返回值列表的括号可以省略。如果一个函数声明没有返回值列表，那么这个
函数不会返回任何值。

```go
// 两个int 类型参数 返回一个 int 类型的值
func max(num1, num2 int) int {
   /* 定义局部变量 */
   var result int

   if (num1 > num2) {
      result = num1
   } else {
      result = num2
   }
   return result 
}

// 返回多个 string 类型的值
func swap(x, y string) (string, string) {
   return y, x
}
```