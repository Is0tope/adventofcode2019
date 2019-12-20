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

const (
	NORTH = iota
	SOUTH = iota
	WEST  = iota
	EAST  = iota
)

var opposite map[int]int = map[int]int{NORTH: SOUTH, SOUTH: NORTH, WEST: EAST, EAST: WEST}

func (p1 point) add(p2 point) point {
	return point{p1.x + p2.x, p1.y + p2.y}
}

func (p1 point) surroundingSpaces() map[int]point {
	ret := make(map[int]point)
	ret[NORTH] = p1.add(point{0, -1})
	ret[SOUTH] = p1.add(point{0, 1})
	ret[EAST] = p1.add(point{1, 0})
	ret[WEST] = p1.add(point{-1, 0})
	return ret
}

func (p1 point) move(dir int) point {
	if dir == NORTH {
		return p1.add(point{0, -1})
	}
	if dir == SOUTH {
		return p1.add(point{0, 1})
	}
	if dir == EAST {
		return p1.add(point{1, 0})
	}
	if dir == WEST {
		return p1.add(point{-1, 0})
	}
	return p1
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
			//fmt.Println(interp.pointer, "took input", input, "set at", dest)
			interp.pointer += 2
		}
		// Output
		if opcode == 4 {
			output = interp.getValue(1, modes)
			//fmt.Println(interp.pointer, "gave output", output, "from", instruction, interp.get(interp.pointer+1))
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

func printGrid(grid map[point]string, width int64, height int64) {
	for y := int64(0); y < height; y++ {
		for x := int64(0); x < width; x++ {
			fmt.Printf(grid[point{x, y}])
		}
		fmt.Printf("\n")
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

	parser := NewInterpreter(code)
	parser.start()

	//tiles := map[int64]string{46: ".", 35: "#"}
	grid := make(map[point]string)
	x, y := int64(0), int64(0)
	width, height := int64(0), int64(0)
	for c := range parser.outchan {
		if c == 10 {
			x = 0
			y++
			continue
		}
		grid[point{x, y}] = string(c)
		if (x + 1) > width {
			width = x + 1
		}
		if (y + 1) > height {
			height = y + 1
		}
		x++
	}
	printGrid(grid, width, height)

	seen := make(map[point]struct{})
	junctions := make(map[point]struct{})

	// Find starting point
	position := point{}
	direction := NORTH
	for k, v := range grid {
		if v == "^" {
			position = k
			direction = NORTH
		}
		if v == "v" {
			position = k
			direction = SOUTH
		}
		if v == ">" {
			position = k
			direction = EAST
		}
		if v == "<" {
			position = k
			direction = WEST
		}
	}
	// Start scanning
	for {
		fmt.Println("at", position, "dir", direction, "tile", grid[position])
		// Have we seen this spot before?
		if _, ok := seen[position]; ok {
			junctions[position] = struct{}{}
		}
		seen[position] = struct{}{}
		// check directions
		dirs := position.surroundingSpaces()
		options := make(map[int]point)
		for k, v := range dirs {
			if grid[v] == "#" && opposite[direction] != k {
				options[k] = v
			}
		}
		fmt.Println(options)
		// Are we at the end?
		if len(options) == 0 {
			fmt.Println("at end of track")
			break
		}
		// Check forward direction
		if _, ok := options[direction]; ok {
			position = position.move(direction)
			continue
		}
		// Decide where to head next (should be one choice)
		for k, v := range options {
			direction = k
			position = v
		}
	}
	fmt.Println("junctions", junctions)
	total := int64(0)
	for k, _ := range junctions {
		total += k.x * k.y
	}
	fmt.Println("total alignment:", total)
}
