package main

import "fmt"

func main()  {
	urls := make([]string, 3)
	urls = append(urls, "hello")
	fmt.Println(len(urls)) // 4

	urls2 := make([]string, 0)
	urls2 = append(urls2, "hello")
	fmt.Println(len(urls2)) // 1
}
