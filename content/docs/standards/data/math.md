---
title: math
---

# math
math 包实现的就是数学函数计算。

## 三角函数

正弦函数，反正弦函数，双曲正弦，反双曲正弦
```
- func Sin(x float64) float64
- func Asin(x float64) float64
- func Sinh(x float64) float64
- func Asinh(x float64) float64
```


一次性返回 `sin`,`cos`

- `func Sincos(x float64) (sin, cos float64)`

余弦函数，反余弦函数，双曲余弦，反双曲余弦
```
- func Cos(x float64) float64
- func Acos(x float64) float64
- func Cosh(x float64) float64
- func Acosh(x float64) float64
```

正切函数，反正切函数，双曲正切，反双曲正切
```
- func Tan(x float64) float64
- func Atan(x float64) float64 和 func Atan2(y, x float64) float64
- func Tanh(x float64) float64
- func Atanh(x float64) float64
```

## 幂次函数
```
- func Cbrt(x float64) float64 // 立方根函数
- func Pow(x, y float64) float64  // x 的幂函数
- func Pow10(e int) float64 // 10 根的幂函数
- func Sqrt(x float64) float64 // 平方根
- func Log(x float64) float64 // 对数函数
- func Log10(x float64) float64 // 10 为底的对数函数
- func Log2(x float64) float64  // 2 为底的对数函数
- func Log1p(x float64) float64 // log(1 + x)
- func Logb(x float64) float64 // 相当于 log2(x) 的绝对值
- func Ilogb(x float64) int // 相当于 log2(x) 的绝对值的整数部分
- func Exp(x float64) float64 // 指数函数
- func Exp2(x float64) float64 // 2 为底的指数函数
- func Expm1(x float64) float64 // Exp(x) - 1
```


## 特殊函数
```
- func Inf(sign int) float64  // 正无穷
- func IsInf(f float64, sign int) bool // 是否正无穷
- func NaN() float64 // 无穷值
- func IsNaN(f float64) (is bool) // 是否是无穷值
- func Hypot(p, q float64) float64 // 计算直角三角形的斜边长
```


## 类型转化函数
```
- func Float32bits(f float32) uint32  // float32 和 unit32 的转换
- func Float32frombits(b uint32) float32 // uint32 和 float32 的转换
- func Float64bits(f float64) uint64 // float64 和 uint64 的转换
- func Float64frombits(b uint64) float64 // uint64 和 float64 的转换
```


## 其他函数
```
- func Abs(x float64) float64 // 绝对值函数
- func Ceil(x float64) float64  // 向上取整
- func Floor(x float64) float64 // 向下取整
- func Mod(x, y float64) float64 // 取模
- func Modf(f float64) (int float64, frac float64) // 分解 f，以得到 f 的整数和小数部分
- func Frexp(f float64) (frac float64, exp int) // 分解 f，得到 f 的位数和指数
- func Max(x, y float64) float64  // 取大值
- func Min(x, y float64) float64  // 取小值
- func Dim(x, y float64) float64 // 复数的维数
- func J0(x float64) float64  // 0 阶贝塞尔函数
- func J1(x float64) float64  // 1 阶贝塞尔函数
- func Jn(n int, x float64) float64 // n 阶贝塞尔函数
- func Y0(x float64) float64  // 第二类贝塞尔函数 0 阶
- func Y1(x float64) float64  // 第二类贝塞尔函数 1 阶
- func Yn(n int, x float64) float64 // 第二类贝塞尔函数 n 阶
- func Erf(x float64) float64 // 误差函数
- func Erfc(x float64) float64 // 余补误差函数
- func Copysign(x, y float64) float64 // 以 y 的符号返回 x 值
- func Signbit(x float64) bool // 获取 x 的符号
- func Gamma(x float64) float64 // 伽玛函数
- func Lgamma(x float64) (lgamma float64, sign int) // 伽玛函数的自然对数
- func Ldexp(frac float64, exp int) float64 // value 乘以 2 的 exp 次幂
- func Nextafter(x, y float64) (r float64) // 返回参数 x 在参数 y 方向上可以表示的最接近的数值，若 x 等于 y，则返回 x
- func Nextafter32(x, y float32) (r float32) // 返回参数 x 在参数 y 方向上可以表示的最接近的数值，若 x 等于 y，则返回 x
- func Remainder(x, y float64) float64 // 取余运算
- func Trunc(x float64) float64 // 截取函数
```
