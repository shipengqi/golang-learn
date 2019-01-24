package main

import (
	"fmt"
	"os/exec"
)

func main(){
	cmd := exec.Command("touch", "test_file")

	err := cmd.Run()
	if err != nil {
		fmt.Println("Execute Command failed:" + err.Error())
		return
	}

	fmt.Println("Execute Command finished.")
}