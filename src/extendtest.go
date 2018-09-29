package main

import "fmt"

type user struct{ 
	name string
	age byte
}

func (u user) ToString() string {
	return fmt.Sprintf("%+v", u)
}

type manager struct {
  user
  title string
}

func main() {
	var m manager
	m.name = "Xiaoming"
	m.age = 18

	fmt.Println(m.ToString())
}