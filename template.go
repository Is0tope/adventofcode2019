package main

import (
	"bufio"
	"os"
)

func main() {
	file, _ := os.Open("../test.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// text := scanner.Text()
	}

	// fmt.Printf("DONE: %s\n", test)
}
