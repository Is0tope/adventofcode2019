package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	file, _ := os.Open("../input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	weight := 0

	for scanner.Scan() {
		text := scanner.Text()
		num, _ := strconv.ParseFloat(text, 0)
		weight += int(math.Floor(num/3)) - 2
	}

	fmt.Printf("DONE: %d\n", weight)
}
