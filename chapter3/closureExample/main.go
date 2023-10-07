package main

import "fmt"

// closure function to generate integers
func generator() func() int {
	i := 0
	return func() int {
		i++
		return i
	}

}

func main() {

	numGenerator := generator()
	for i := 0; i < 5; i++ {
		fmt.Print(numGenerator(), "\t")
	}
	fmt.Println()
}
