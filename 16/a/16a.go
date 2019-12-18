package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func patternOffset(element, length, repeat int) int {
	return ((1 + element) / repeat) % length
}

func abs(x int64) int64 {
	if x < 0 {
		x *= -1
	}
	return x
}

func phase(signal, pattern []int64) []int64 {
	ret := make([]int64, len(signal))
	for i := range signal {
		total := int64(0)
		for j, b := range signal {
			total += b * pattern[patternOffset(j, len(pattern), i+1)]
			//fmt.Printf("%d*%d ", b, pattern[patternOffset(j, len(pattern), i+1)])
		}
		ret[i] = abs(total) % 10
		//fmt.Printf("= %d\n", ret[i])
	}
	return ret
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

	pattern := []int64{0, 1, 0, -1}

	input := digits
	for i := 0; i < 100; i++ {
		input = phase(input, pattern)
	}
	fmt.Println(input[:8])

}
