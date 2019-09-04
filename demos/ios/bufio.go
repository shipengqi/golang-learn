package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	reader := bufio.NewReader(strings.NewReader("Hello \nworld"))
	line, _ := reader.ReadSlice('\n')
	fmt.Printf("the line:%s\n", line) // the line:Hello
	line2, _ := reader.ReadBytes('\n')
	fmt.Printf("the line:%s\n", line) // the line:world
	fmt.Printf("the line2:%s\n", line2) // the line:world
	fmt.Println(string(line2)) // world
}