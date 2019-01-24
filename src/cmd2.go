package main

import (
	"fmt"
	"os/exec"
	"sync"
)

func exe_cmd(cmd string, wg *sync.WaitGroup) {
	fmt.Println(cmd)
	out, err := exec.Command("sh","-c", cmd).Output()
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
	wg.Done()
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(3)

	x := []string{"echo newline", "echo newline1", "echo newline2"}
	go exe_cmd(x[0], wg)
	go exe_cmd(x[1], wg)
	go exe_cmd(x[2], wg)

	wg.Wait()
}