package main

import "fmt"

func main() {
	x := []int{100, 101, 102}

	for key, value := range x {
		if key == 0 {
			x[0] += 100
			x[1] += 200
			x[2] += 300
		}
		fmt.Printf("value: %d, x: %d\n", value, x[key])
	}

	for key, value := range x[:] {
		if key == 0 {
			x[0] += 100
			x[1] += 200
			x[2] += 300
		}
		fmt.Printf("value: %d, x: %d\n", value, x[key])
	}

	for _, value := range x {
		fmt.Println(value)
	}
}