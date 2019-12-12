package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

type planet struct {
	x, y, z    int
	vx, vy, vz int
}

var GRAVITY = 1

func NewPlanet(x, y, z int) planet {
	p := planet{x, y, z, 0, 0, 0}
	return p
}

func abs(x int) int {
	if x < 0 {
		x *= -1
	}
	return x
}

func getGravity(src, dest int) int {
	if src > dest {
		return -GRAVITY
	}
	if src < dest {
		return GRAVITY
	}
	return 0
}

func applyGravity(current *planet, target planet) {
	//fmt.Println("calculating gravity", current, target)
	current.vx += getGravity(current.x, target.x)
	current.vy += getGravity(current.y, target.y)
	current.vz += getGravity(current.z, target.z)
}

func applyVelocity(current *planet) {
	current.x += current.vx
	current.y += current.vy
	current.z += current.vz
}

func (p planet) kineticEnergy() int {
	return abs(p.vx) + abs(p.vy) + abs(p.vz)
}

func (p planet) potentialEnergy() int {
	return abs(p.x) + abs(p.y) + abs(p.z)
}

func runOverPlanets(planets []planet, f func(p1 *planet, p2 planet)) {
	for i, _ := range planets {
		for j, _ := range planets {
			if planets[i] == planets[j] {
				continue
			}
			f(&planets[i], planets[j])
		}
	}
}

func stateString(planets []planet, dim string) string {
	str := ""
	for _, p := range planets {
		if dim == "x" {
			str += fmt.Sprintf("(%d,%d)", p.x, p.vx)
		}
		if dim == "y" {
			str += fmt.Sprintf("(%d,%d)", p.y, p.vy)
		}
		if dim == "z" {
			str += fmt.Sprintf("(%d,%d)", p.z, p.vz)
		}
	}
	return str
}

// Least Common Multiple (dirty method)
func DirtyLCM(numbers []int) int64 {
	// find smallest number (assume > 0)
	min := int64(math.MaxInt64)
	for _, n := range numbers {
		if int64(n) < min {
			min = int64(n)
		}
	}
	// Start iterating until it all works
	mul := int64(1)
	for {
		bad := false
		for _, n := range numbers {
			if (mul*min)%int64(n) != 0 {
				bad = true
			}
		}
		if !bad {
			return mul * min
		}
		mul++
	}
}

func main() {
	file, _ := os.Open("../input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	r, _ := regexp.Compile("<x=(-?\\d+), y=(-?\\d+), z=(-?\\d+)>")
	planets := []planet{}
	for scanner.Scan() {
		text := scanner.Text()
		matches := r.FindStringSubmatch(text)
		x, _ := strconv.Atoi(matches[1])
		y, _ := strconv.Atoi(matches[2])
		z, _ := strconv.Atoi(matches[3])
		planets = append(planets, NewPlanet(x, y, z))
	}
	fmt.Println("planets:", planets)

	// Run simulations in each dimension
	statesX := make(map[string]bool)
	statesY := make(map[string]bool)
	statesZ := make(map[string]bool)
	iterX := 0
	iterY := 0
	iterZ := 0
	time := 0
	for {
		if 0 == time%100000 {
			fmt.Println("time is", time)
		}
		// Apply gravity first
		runOverPlanets(planets, applyGravity)
		// Apply velocity
		for i, _ := range planets {
			applyVelocity(&planets[i])
		}
		strX := stateString(planets, "x")
		strY := stateString(planets, "y")
		strZ := stateString(planets, "z")
		// Is this seen before
		if _, ok := statesX[strX]; ok {
			if iterX == 0 {
				iterX = time
			}
		}
		if _, ok := statesY[strY]; ok {
			if iterY == 0 {
				iterY = time
			}
		}
		if _, ok := statesZ[strZ]; ok {
			if iterZ == 0 {
				iterZ = time
			}
		}
		// Check if all are done
		if iterX != 0 && iterY != 0 && iterZ != 0 {
			break
		}
		// Add to states
		statesX[strX] = true
		statesY[strY] = true
		statesZ[strZ] = true
		time++
	}
	fmt.Println("All duplicate state time have been found at time", time)
	fmt.Println("X", iterX)
	fmt.Println("Y", iterY)
	fmt.Println("Z", iterZ)

	// Find the LCM of the three numbers
	fmt.Println("LCM is", DirtyLCM([]int{iterX, iterY, iterZ}))
}
