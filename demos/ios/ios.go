package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := strings.NewReader("Hello world")
	p := make([]byte, 6)
	n, err := reader.ReadAt(p, 2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s, %d\n", p, n) // llo wo, 6

	file, err := os.Create("writeAt.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, _ = file.WriteString("Hello world----ignore")
	n, err = file.WriteAt([]byte("Golang"), 15)
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
}
