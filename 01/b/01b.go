package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func calcFuel(weight int64) int64 {
	fuel := int64(0)
	for weight > 0 {
		weight = int64(math.Floor(float64(weight)/3)) - 2
		if weight < 0 {
			weight = 0
		}
		fuel += weight
	}
	return fuel
}

func main() {
	file, _ := os.Open("../input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	weight := int64(0)

	for scanner.Scan() {
		text := scanner.Text()
		num, _ := strconv.ParseInt(text, 0, 0)
		weight += calcFuel(num)
	}

	fmt.Printf("DONE: %d\n", weight)
}
