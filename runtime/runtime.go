package main

import (
	"fmt"
	"runtime"
)

type Stats struct {
	version       string
	numGoRoutines int
	stackTrace    []byte
}

func main() {
	currentStats := new(Stats)
	currentStats.numGoRoutines = runtime.NumGoroutine()
	currentStats.version = runtime.Version()
	currentStats.stackTrace = runtime.ReadTrace()

	anothervar := "hello"

	fmt.Printf("Go Version: %v\n", currentStats.version)
	fmt.Printf("# active goroutines: %v\n", currentStats.numGoRoutines)
	fmt.Printf("stack trace: %v\n", currentStats.stackTrace)
	fmt.Println(anothervar)
}
