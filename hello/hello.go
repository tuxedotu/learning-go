package hello

import "fmt"

func SayHello() {
	fmt.Println("Hello go :)")
}

func PrintNLines(n int) {
	for counter := 0; counter < n; counter++ {
		fmt.Println("Current count: ", counter)
	}
}

func 
