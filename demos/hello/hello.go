package hello

import "fmt"

const spanish = "Spanish"
const french = "French"
const helloPrefix = "Hello, "
const spanishHelloPrefix = "Hola, "
const frenchHelloPrefix = "Bonjour, "

func createGreetingPrefix(language string) (prefix string) {
	switch language {
	case spanish:
		prefix = spanishHelloPrefix
	case french:
		prefix = frenchHelloPrefix
	default:
		prefix = helloPrefix
	}
	return prefix
}

func Hello(language string, name string) string {
	if name == "" {
		name = "world"
	}
 	return createGreetingPrefix(language) + name
}

func main() {
	fmt.Println(Hello("", ""))
}
