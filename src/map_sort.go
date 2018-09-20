package main

import (
	"sort"
	"fmt"
)

func main() {
	ages := make(map[string]int)
  ages["alice"] = 31
  ages["charlie"] = 34
	var names []string
	for name := range ages { // 第二个循环变量可以忽略
			names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
			fmt.Printf("%s\t%d\n", name, ages[name])
	}	
}