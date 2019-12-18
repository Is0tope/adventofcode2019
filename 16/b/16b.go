package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func abs(x int64) int64 {
	if x < 0 {
		x *= -1
	}
	return x
}

func main() {
	file, _ := os.Open("../input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	text := scanner.Text()
	digits := make([]int64, len(text))
	for i := 0; i < len(digits); i++ {
		d, _ := strconv.Atoi(text[i : i+1])
		digits[i] = int64(d)
	}
	digits2 := make([]int64, 10000*len(digits))
	for i := 0; i < len(digits2); i++ {
		digits2[i] = digits[i%len(digits)]
	}

	offset, _ := strconv.Atoi(text[:7])

	input := digits2
	for i := 0; i < 100; i++ {
		fmt.Println("phase", i)
		for j := len(input) - 2; j >= offset; j-- {
			input[j] = (input[j] + input[j+1]) % 10
		}
	}
	fmt.Println(input[offset : offset+8])

}
