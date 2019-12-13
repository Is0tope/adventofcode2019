package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type interpreter struct {
	halted  bool
	code    map[int64]int64
	pointer int64
	rbase   int64
	inchan  chan int64
	outchan chan int64
}

type point struct {
	x, y int64
}

type game struct {
	brain  interpreter
	screen map[point]int
}

func NewGame(code map[int64]int64) game {
	g := game{}
	g.brain = NewInterpreter(code)
	g.screen = make(map[point]int)
	return g
}

func (g *game) run() {
	g.brain.start()

	for {
		// Get the coords & type
		x, _ := g.brain.output()
		y, _ := g.brain.output()
		t, more := g.brain.output()
		g.screen[point{x, y}] = int(t)
		fmt.Println(x, y, t)
		// Check if we are done
		if !more {
			fmt.Println("got all coords")
			break
		}
	}
}

func NewInterpreter(code map[int64]int64) interpreter {
	inter := interpreter{}
	inter.code = code
	inter.pointer = 0
	inter.rbase = 0
	inter.inchan = make(chan int64)
	inter.outchan = make(chan int64, 1000)

	return inter
}

func (interp interpreter) getAddress(offset int64, modes []int) int64 {
	mode := 1
	if (offset - 1) < int64(len(modes)) {
		mode = modes[offset-1]
	}
	if mode == 1 {
		return interp.get(interp.pointer + offset)
	} else if mode == 2 {
		return interp.rbase + interp.get(interp.pointer+offset)
	} else {
		return interp.get(interp.code[interp.pointer+offset])
	}
}

func (interp interpreter) getValue(offset int64, modes []int) int64 {
	mode := 0
	if (offset - 1) < int64(len(modes)) {
		mode = modes[offset-1]
	}
	if mode == 0 {
		return interp.get(interp.code[interp.pointer+offset])
	} else if mode == 1 {
		return interp.get(interp.pointer + offset)
	} else {
		return interp.get(interp.rbase + interp.code[interp.pointer+offset])
	}
}

func (inter *interpreter) get(pointer int64) int64 {
	if cur, ok := inter.code[pointer]; ok {
		return cur
	} else {
		return 0
	}
}

func (inter *interpreter) set(pointer int64, val int64) {
	inter.code[pointer] = val
}

func (inter *interpreter) start() {
	go inter.run()
}

func (inter *interpreter) input(in int64) {
	inter.inchan <- in
}

func (inter *interpreter) output() (int64, bool) {
	val, more := <-inter.outchan
	return val, more
}

func (interp *interpreter) run() {
	output := int64(0)
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
		//fmt.Println(interp.pointer, []int64{instruction, interp.get(interp.pointer + 1)}, opcode, modes, interp.rbase)
		// Exit code
		if opcode == 99 {
			fmt.Println("Found exit code")
			interp.halted = true
			close(interp.outchan)
			break
		}
		// Addition
		if opcode == 1 {
			res := interp.getValue(1, modes) + interp.getValue(2, modes)
			dest := interp.getAddress(3, modes)
			//fmt.Println(interp.pointer, "add", interp.getValue(1, modes), "to", interp.getValue(2, modes), "with result", res, "send to", dest)
			interp.code[dest] = res
			interp.pointer += 4
		}
		// Multiplication
		if opcode == 2 {
			res := interp.getValue(1, modes) * interp.getValue(2, modes)
			dest := interp.getAddress(3, modes)
			//fmt.Println(interp.pointer, "multiply", interp.getValue(1, modes), "to", interp.getValue(2, modes), "with result", res, "send to", dest)
			interp.code[dest] = res
			interp.pointer += 4
		}
		// Input
		if opcode == 3 {
			//interp.commchan <- true
			input := <-interp.inchan
			dest := interp.getAddress(1, modes)
			interp.set(dest, input)
			fmt.Println(interp.pointer, "took input", input, "set at", dest)
			interp.pointer += 2
		}
		// Output
		if opcode == 4 {
			output = interp.getValue(1, modes)
			fmt.Println(interp.pointer, "gave output", output, "from", instruction, interp.get(interp.pointer+1))
			//interp.commchan <- true
			interp.outchan <- output
			interp.pointer += 2
		}
		// Jump if true, Jump if false
		if opcode == 5 || opcode == 6 {
			cond := interp.getValue(1, modes)
			if (cond != 0 && opcode == 5) || (cond == 0 && opcode == 6) {
				interp.pointer = interp.getValue(2, modes)
				//fmt.Println(interp.pointer, "jumping because", opcode, "to", interp.pointer)
			} else {
				interp.pointer += 3
				//fmt.Println(interp.pointer, "not jumping")
			}
		}
		// less than
		if opcode == 7 {
			ret := int64(0)
			if interp.getValue(1, modes) < interp.getValue(2, modes) {
				ret = int64(1)
			}
			//fmt.Println(interp.pointer, "evaluating", interp.getValue(1, modes), "<", interp.getValue(2, modes), "to", ret)
			interp.code[interp.getAddress(3, modes)] = ret
			interp.pointer += 4
		}
		// equals
		if opcode == 8 {
			ret := int64(0)
			if interp.getValue(1, modes) == interp.getValue(2, modes) {
				ret = int64(1)
			}
			//fmt.Println(interp.pointer, "evaluating", interp.getValue(1, modes), "==", interp.getValue(2, modes), "to", ret)
			interp.code[interp.getAddress(3, modes)] = ret
			interp.pointer += 4
		}
		// adjust relative base
		if opcode == 9 {
			//fmt.Println(interp.pointer, "incrementing rbase", interp.rbase, "by", interp.getValue(1, modes))
			interp.rbase += interp.getValue(1, modes)
			interp.pointer += 2
		}
	}
}

func main() {
	file, _ := os.Open("../input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	text := scanner.Text()
	code := make(map[int64]int64)

	for i, num := range strings.Split(text, ",") {
		num2, _ := strconv.ParseInt(num, 0, 0)
		code[int64(i)] = num2
	}

	gam := NewGame(code)
	gam.run()
	// Count the number of block (2) tiles
	blocks := 0
	for _, v := range gam.screen {
		if v == 2 {
			blocks++
		}
	}
	fmt.Println("Total blocks:", blocks)
}
