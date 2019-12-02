package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

	pointer := 0

	for {
		instruction := code[pointer]
		fmt.Println("Instruction: ", instruction)
		fmt.Println("Pointer: ", pointer)

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
	fmt.Printf("DONE: %d\n", code[0])
}
