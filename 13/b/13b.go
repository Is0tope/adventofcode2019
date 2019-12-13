package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type interpreter struct {
	halted  bool
	code    map[int64]int64
	pointer int64
	rbase   int64
	inchan  chan int64
	outchan chan int64
	reqchan chan bool
}

type point struct {
	x, y int64
}

type game struct {
	brain      interpreter
	screen     map[point]int
	score      int
	joystick   int
	width      int64
	height     int64
	refresh    int
	speed      int
	renderstop chan struct{}
	screenlock sync.RWMutex
}

var tiles = map[int]string{0: " ", 1: "#", 2: "X", 3: "=", 4: "O"}

func NewGame(code map[int64]int64, width int64, height int64) *game {
	g := game{}
	g.brain = NewInterpreter(code)
	g.screen = make(map[point]int)
	g.score = 0
	g.joystick = 0
	g.width = width
	g.height = height
	g.refresh = 14 // hz
	g.speed = 100  // hz
	g.renderstop = make(chan struct{})
	g.screenlock = sync.RWMutex{}
	return &g
}

func (g *game) run() {
	// start rendering
	go g.render()
	// start listening to input
	go g.input()
	// start the program
	g.brain.start()

	for {
		// Get the coords & type
		x, _ := g.brain.output()
		y, _ := g.brain.output()
		t, more := g.brain.output()

		// Check for scoring instruction
		if x == -1 && y == 0 {
			g.score = int(t)
		} else {
			g.screenlock.Lock()
			g.screen[point{x, y}] = int(t)
			g.screenlock.Unlock()
		}
		// Check if we are done
		if !more {
			close(g.renderstop)
			break
		}
	}
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("FINAL SCORE:", g.score)
}

// Not threadsafe
func (g *game) getPixel(x, y int64) string {
	p := point{x, y}
	t, ok := g.screen[p]
	if !ok {
		t = 0
	}
	return tiles[t]
}

func (g *game) printScreen() {
	g.screenlock.RLock()
	defer g.screenlock.RUnlock()
	// Clear terminal (apparently...)
	//fmt.Println("\033[2J")
	// Print the screen
	fmt.Println("SCORE:", g.score)
	for y := int64(0); y < g.height; y++ {
		for x := int64(0); x < g.width; x++ {
			fmt.Printf(g.getPixel(x, y))
		}
		fmt.Printf("\n")
	}
}

func (g *game) render() {
	for {
		select {
		case <-g.renderstop:
			break // exit
		default:
		}
		// render
		g.printScreen()
		// Calculate sleep time
		p := int64(math.Ceil(1000.0 / float64(g.refresh)))
		time.Sleep(time.Duration(p) * time.Millisecond)
	}
}

func (g *game) input() {
	for range g.brain.reqchan {
		// Slow the game down
		slp := int64(math.Ceil(1000.0 / float64(g.speed)))
		time.Sleep(time.Duration(slp) * time.Millisecond)
		// Find out where the ball is, and make the game autoplay itself
		p := g.getPaddlePos()
		b := g.getBallPos()
		g.joystick = 0
		if b.x > p.x {
			g.joystick = 1
		}
		if b.x < p.x {
			g.joystick = -1
		}
		g.brain.inchan <- int64(g.joystick)
	}
}

// So inefficient
func (g *game) getBallPos() point {
	for k, v := range g.screen {
		if v == 4 {
			return k
		}
	}
	return point{}
}

func (g *game) getPaddlePos() point {
	for k, v := range g.screen {
		if v == 3 {
			return k
		}
	}
	return point{}
}

func NewInterpreter(code map[int64]int64) interpreter {
	inter := interpreter{}
	inter.code = code
	inter.pointer = 0
	inter.rbase = 0
	inter.inchan = make(chan int64)
	inter.outchan = make(chan int64, 1000)
	inter.reqchan = make(chan bool)
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
			close(interp.reqchan)
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
			interp.reqchan <- true
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

	gam := NewGame(code, 40, 20)
	gam.run()
}
