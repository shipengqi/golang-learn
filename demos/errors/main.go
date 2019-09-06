package errors

import "fmt"

func soleTitle() (title string, err error) {
	type bailout struct {}
	defer func() {
		switch p := recover(); p {
		case nil: // no panic
		case bailout{}: // "expected" panic
			err = fmt.Errorf("multiple title elements")
			fmt.Println(err.Error())
		default:
			panic(p) // unexpected panic; carry on panicking
		}
	}()
	panic(bailout{})
}

func main() {
	fmt.Printf("Calling test\r\n")
	_, _ = soleTitle()
	fmt.Printf("Test completed\r\n")
}
