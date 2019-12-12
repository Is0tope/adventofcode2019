package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type planet struct {
	x, y, z    int
	vx, vy, vz int
}

var GRAVITY = 1
var STEPS = 1000

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

	// Run simulations
	for time := 0; time < STEPS; time++ {
		fmt.Println("time is", time)
		// Apply gravity first
		runOverPlanets(planets, applyGravity)
		// Apply velocity
		for i, _ := range planets {
			applyVelocity(&planets[i])
		}
		fmt.Println(planets)
	}

	// Calculate total energy
	totalEnergy := 0
	for _, p := range planets {
		totalEnergy += p.kineticEnergy() * p.potentialEnergy()
	}
	fmt.Println("TOTAL ENERGY:", totalEnergy)
}
