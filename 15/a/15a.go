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
	out     int64
	rbase   int64
}

type point struct {
	x, y int64
}

type state struct {
	brain    interpreter
	position point
	distance int64
	status   int64
}

var tiles = map[int64]string{0: "#", 1: ".", 2: "X"}

var GRID map[point]string // " " = empty, "#" = wall, "X" = goal

func NewState(code map[int64]int64, position point, distance int64) *state {
	s := state{}
	s.brain = NewInterpreter(code)
	s.position = position
	s.distance = distance
	s.status = -1
	return &s
}

func (s state) CloneState() *state {
	ret := state{}
	ret.brain = *s.brain.CloneInterpreter()
	ret.position = s.position
	ret.distance = s.distance
	ret.status = s.status
	return &ret
}

func copyMap(source map[int64]int64) map[int64]int64 {
	res := make(map[int64]int64)
	for k, v := range source {
		res[k] = v
	}
	return res
}

func (s *state) move(dir int64) int64 {
	// Send input
	s.brain.input(dir)
	// Get status
	status := s.brain.output()
	fmt.Println("status", status)
	// Set the status of this state to status
	s.status = status
	// Make the adjustments if we moved
	if status > 0 {
		if dir == 1 {
			s.position.y++
		}
		if dir == 2 {
			s.position.y--
		}
		if dir == 3 {
			s.position.x--
		}
		if dir == 4 {
			s.position.x++
		}
	}
	return status
}

func getTile(p point) string {
	t, ok := GRID[p]
	if !ok {
		t = "_"
	}
	return t
}

func printGrid(pos point, n int64) {
	for y := int64(40); y > -20; y-- {
		for x := int64(-n); x < n; x++ {
			if pos.x == x && pos.y == y {
				fmt.Printf("@")
			} else {
				t := getTile(point{x, y})
				fmt.Printf(t)
			}
		}
		fmt.Printf("\n")
	}
}

func NewInterpreter(code map[int64]int64) interpreter {
	inter := interpreter{}
	inter.code = code
	inter.pointer = 0
	inter.rbase = 0
	return inter
}

func (source interpreter) CloneInterpreter() *interpreter {
	inter := interpreter{}
	inter.code = copyMap(source.code)
	inter.pointer = source.pointer
	inter.rbase = source.rbase
	return &inter
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

func (inter *interpreter) input(in int64) {
	inter.run(in)
}

func (inter *interpreter) output() int64 {
	return inter.out
}

func (interp *interpreter) run(input int64) {
	output := int64(0)
	firstinput := true
	for {
		// Instruction
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
			if !firstinput {
				break
			}
			dest := interp.getAddress(1, modes)
			interp.set(dest, input)
			fmt.Println(interp.pointer, "took input", input, "set at", dest)
			interp.pointer += 2
			firstinput = false
		}
		// Output
		if opcode == 4 {
			output = interp.getValue(1, modes)
			fmt.Println(interp.pointer, "gave output", output, "from", instruction, interp.get(interp.pointer+1))
			interp.out = output
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

	// Set up GRID
	GRID = make(map[point]string)
	// Set up queue
	QUEUE := make([]state, 0)

	// Initial state
	init := NewState(code, point{0, 0}, 0)
	init.status = 1
	QUEUE = append(QUEUE, *init)

	cnt := 0
	for len(QUEUE) > 0 {
		cnt++
		// pop off a state
		s := QUEUE[0]
		QUEUE = QUEUE[1:]
		// Have we seen this state before?
		_, seen := GRID[s.position]
		if seen {
			continue
		}
		// Is this the final state?
		if s.status == 2 {
			fmt.Println("FOUND OXYGEN SYSTEM @", s.position, " and distance", s.distance)
			break
		}
		// Start the computer
		fmt.Println("processing state s at", s.position, "distance", s.distance)
		//time.Sleep(100 * time.Millisecond)

		// Update the grid
		GRID[s.position] = tiles[s.status]
		// Test all of the directions but move back
		states := []state{}
		opposite := map[int64]int64{1: 2, 2: 1, 3: 4, 4: 3}
		for _, dir := range []int64{1, 2, 3, 4} {
			if s.move(dir) > 0 {
				clone := *s.CloneState()
				clone.distance++
				states = append(states, clone)
				s.move(opposite[dir])
			} else {
				p := s.position
				if dir == 1 {
					p.y++
				}
				if dir == 2 {
					p.y--
				}
				if dir == 3 {
					p.x--
				}
				if dir == 4 {
					p.x++
				}
				GRID[p] = "#"
			}
			//fmt.Println("s at", s.position, "direction", dir, "status", s.status)
		}
		if len(states) > 0 {
			for _, n := range states {
				QUEUE = append(QUEUE, n)
			}
		}
		fmt.Println("MAP:")
		printGrid(s.position, 40)
		// if cnt == 100 {
		// 	break
		// }
	}

}
