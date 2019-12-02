package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func runIntcode(code []int64) []int64 {
	pointer := 0
	for {
		instruction := code[pointer]
		// Exit code
		if instruction == 99 {
			fmt.Println("Found exit code")
			break
		}
		// Addition
		if instruction == 1 {
			res := code[code[pointer+1]] + code[code[pointer+2]]
			dest := code[pointer+3]
			code[dest] = res
		}
		// Multiplication
		if instruction == 2 {
			res := code[code[pointer+1]] * code[code[pointer+2]]
			dest := code[pointer+3]
			code[dest] = res
		}
		pointer += 4
	}
	return code
}

func main() {
	file, _ := os.Open("../input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	text := scanner.Text()
	var code []int64

	for _, num := range strings.Split(text, ",") {
		num2, _ := strconv.ParseInt(num, 0, 0)
		code = append(code, num2)
	}

	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			codecopy := make([]int64, len(code))
			copy(codecopy, code)
			codecopy[1] = int64(i)
			codecopy[2] = int64(j)
			out := runIntcode(codecopy)
			if out[0] == 19690720 {
				fmt.Println("Number 1: ", i)
				fmt.Println("Number 1: ", j)
				fmt.Printf("DONE: %d\n", 100*i+j)
				os.Exit(0)
			}
		}
	}
	fmt.Println("Couldn't find it!")
}
