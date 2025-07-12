package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Test program starting...")

	for i := 0; i < 5; i++ {
		result := fibonacci(i)
		fmt.Printf("fibonacci(%d) = %d\n", i, result)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("Test program finished")
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}
