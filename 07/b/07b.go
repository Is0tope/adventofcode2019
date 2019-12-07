package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type amplifier struct {
	phase  int
	input  int
	output int
	interp interpreter
}

type interpreter struct {
	halted  bool
	code    []int
	pointer int
}

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

func (interp *interpreter) run(input int) int {
	output := 0
	inputCalled := false
	for {
		instruction := interp.code[interp.pointer]
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
			interp.halted = true
			break
		}
		// Addition
		if opcode == 1 {
			res := getValue(interp.code, interp.pointer, 1, modes) + getValue(interp.code, interp.pointer, 2, modes)
			dest := interp.code[interp.pointer+3]
			interp.code[dest] = res
			interp.pointer += 4
		}
		// Multiplication
		if opcode == 2 {
			res := getValue(interp.code, interp.pointer, 1, modes) * getValue(interp.code, interp.pointer, 2, modes)
			dest := interp.code[interp.pointer+3]
			interp.code[dest] = res
			interp.pointer += 4
		}
		// Input
		if opcode == 3 {
			if inputCalled {
				break
			}
			interp.code[interp.code[interp.pointer+1]] = input
			//fmt.Println("took input",input)
			interp.pointer += 2
			inputCalled = true
		}
		// Output
		if opcode == 4 {
			output = getValue(interp.code, interp.pointer, 1, modes)
			//.Println("gave output", output)
			interp.pointer += 2
			break
		}
		// Jump if true, Jump if false
		if opcode == 5 || opcode == 6 {
			cond := getValue(interp.code, interp.pointer, 1, modes)
			if (cond != 0 && opcode == 5) || (cond == 0 && opcode == 6) {
				interp.pointer = getValue(interp.code, interp.pointer, 2, modes)
			} else {
				interp.pointer += 3
			}
		}
		// less than
		if opcode == 7 {
			ret := 0
			if getValue(interp.code, interp.pointer, 1, modes) < getValue(interp.code, interp.pointer, 2, modes) {
				ret = 1
			}
			interp.code[interp.code[interp.pointer+3]] = ret
			interp.pointer += 4
		}
		// equals
		if opcode == 8 {
			ret := 0
			if getValue(interp.code, interp.pointer, 1, modes) == getValue(interp.code, interp.pointer, 2, modes) {
				ret = 1
			}
			interp.code[interp.code[interp.pointer+3]] = ret
			interp.pointer += 4
		}
	}
	return output
}

func (a *amplifier) Execute(input int) {
	fmt.Println("Running AMP with phase: ", a.phase, " input: ", input)
	a.output = a.interp.run(input)
	fmt.Println("Finished running with output ", a.output)
}

func NewAmplifier(code []int, phase int) *amplifier {
	amp := new(amplifier)
	amp.interp.code = make([]int, len(code))
	copy(amp.interp.code, code)
	amp.phase = phase
	amp.output = 0
	amp.Execute(amp.phase)
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
	phaselist := Combination([]int{5, 6, 7, 8, 9})
	for _, phases := range phaselist {
		input := 0
		output := 0
		amps := make([]*amplifier, 0)
		// Set it up
		for _, phase := range phases {
			amp := NewAmplifier(code, phase)
			amps = append(amps, amp)
		}
		fmt.Println(amps)
		// Run it
		for !amps[len(amps)-1].interp.halted {
			for _, a := range amps {
				a.Execute(input)
				if !a.interp.halted {
					output = a.output
				}
				input = a.output
			}
		}
		if output > maxoutput {
			maxoutput = output
			copy(maxphases, phases)
		}
		fmt.Println("phases: ", phases, " got output: ", output)
	}
	fmt.Println("max output ", maxoutput, " at phases ", maxphases)
}
