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

func runIntcode(code []int, input []int) []int {
	pointer := 0
	inpointer := 0
	output := make([]int, 0)
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
			//fmt.Println("Found exit code")
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
			code[code[pointer+1]] = input[inpointer]
			//fmt.Println("took input",input[inpointer])
			pointer += 2
			inpointer++
		}
		// Output
		if opcode == 4 {
			output = append(output, getValue(code, pointer, 1, modes))
			//.Println("gave output", output)
			pointer += 2
		}
		// Jump if true, Jump if false
		if opcode == 5 || opcode == 6 {
			cond := getValue(code, pointer, 1, modes)
			if (cond != 0 && opcode == 5) || (cond == 0 && opcode == 6) {
				pointer = getValue(code, pointer, 2, modes)
			} else {
				pointer += 3
			}
		}
		// less than
		if opcode == 7 {
			ret := 0
			if getValue(code, pointer, 1, modes) < getValue(code, pointer, 2, modes) {
				ret = 1
			}
			code[code[pointer+3]] = ret
			pointer += 4
		}
		// equals
		if opcode == 8 {
			ret := 0
			if getValue(code, pointer, 1, modes) == getValue(code, pointer, 2, modes) {
				ret = 1
			}
			code[code[pointer+3]] = ret
			pointer += 4
		}
	}
	return output
}

type amplifier struct {
	code   []int
	phase  int
	input  int
	output int
}

func (a *amplifier) Execute() {
	fmt.Println("Running AMP with phase: ", a.phase, " input: ", a.input)
	a.output = runIntcode(a.code, []int{a.phase, a.input})[0]
	fmt.Println("Finished running with output ", a.output)
}

func NewAmplifier(code []int, phase int, input int) *amplifier {
	amp := new(amplifier)
	amp.code = make([]int, len(code))
	copy(amp.code, code)
	amp.phase = phase
	amp.input = input
	amp.output = 0
	return amp
}

func Combination(list []int) [][]int {
	if len(list) <= 1 {
		return [][]int{list}
	}
	ret := make([][]int, 0)
	for i, x := range list {
		sub := make([]int, len(list))
		copy(sub, list)
		sub = append(sub[:i], sub[i+1:]...)
		for _, sl := range Combination(sub) {
			part := append([]int{x}, sl...)
			ret = append(ret, part)
		}
	}
	return ret
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

	// Test phase settings
	maxoutput := 0
	maxphases := make([]int, 5)
	phaselist := Combination([]int{0, 1, 2, 3, 4})
	for _, phases := range phaselist {
		input := 0
		output := 0
		for _, phase := range phases {
			amp := NewAmplifier(code, phase, input)
			amp.Execute()
			output = amp.output
			input = amp.output
		}
		if output > maxoutput {
			maxoutput = output
			copy(maxphases, phases)
		}
		fmt.Println("phases: ", phases, " got output: ", output)
	}
	fmt.Println("max output ", maxoutput, " at phases ", maxphases)
}
