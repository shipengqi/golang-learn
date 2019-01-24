package main

import (
	"fmt"
	"os/exec"
)


func main() {
	out, err := exec.Command("echo", "00000").Output()
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
}