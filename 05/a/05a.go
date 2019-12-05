package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getValue(code []int, pointer int, offset int, modes []int) int {
	mode := 0
	if (offset - 1) < len(modes) {
		mode = modes[offset-1]
	}
	if mode == 0 {
		return code[code[pointer+offset]]
	} else {
		return code[pointer+offset]
	}
}

func runIntcode(code []int, input int, output *int) []int {
	pointer := 0
	for {
		instruction := code[pointer]
		// get opcode
		opcode := instruction % 100
		// Handle modes
		modes := make([]int, 0)
		if instruction > 99 {
			val := int((instruction - opcode) / 100)
			for val > 0 {
				m := val % 10
				modes = append(modes, m)
				val = int((val - m) / 10)
			}
		}
		// Exit code
		if opcode == 99 {
			fmt.Println("Found exit code")
			break
		}
		// Addition
		if opcode == 1 {
			res := getValue(code, pointer, 1, modes) + getValue(code, pointer, 2, modes)
			dest := code[pointer+3]
			code[dest] = res
			pointer += 4
		}
		// Multiplication
		if opcode == 2 {
			res := getValue(code, pointer, 1, modes) * getValue(code, pointer, 2, modes)
			dest := code[pointer+3]
			code[dest] = res
			pointer += 4
		}
		// Input
		if opcode == 3 {
			code[code[pointer+1]] = input
			fmt.Println("took input")
			pointer += 2
		}
		// Output
		if opcode == 4 {
			*output = getValue(code, pointer, 1, modes)
			fmt.Println("gave output", *output)
			pointer += 2
		}
	}
	return code
}

func main() {
	file, _ := os.Open("../input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	text := scanner.Text()
	var code []int

	for _, num := range strings.Split(text, ",") {
		num2, _ := strconv.ParseInt(num, 0, 0)
		code = append(code, int(num2))
	}

	output := 0
	runIntcode(code, 1, &output)

}
