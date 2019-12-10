package main

import (
	"bufio"
	"fmt"
	"os"
)

type point struct {
	x, y int
}

type vec2 point

// Greatest Common Divisor
func GCD(a, b int) int {
	if b > a {
		tmp := a
		a = b
		b = tmp
	}
	for {
		if b == 0 {
			return a
		}
		if a == 0 {
			return b
		}
		a, b = b, a%b
	}
}

func abs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}

func sign(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

func incrementOrOne(m map[vec2]int, p vec2) {
	if _, ok := m[p]; ok {
		m[p]++
	} else {
		m[p] = 1
	}
}

func NewVector(a, b point) vec2 {
	vy := b.y - a.y
	vx := b.x - a.x
	// Find the greatest common divisor to get the lowest fraction
	div := GCD(abs(vy), abs(vx))
	vy /= div
	vx /= div
	// Normalise case of 0
	if vy == 0 {
		vx = sign(vx)
	}
	if vx == 0 {
		vy = sign(vy)
	}
	return vec2{vx, vy}
}

func main() {
	file, _ := os.Open("../input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	coords := make(map[point]map[vec2]int)
	y := 0
	for scanner.Scan() {
		text := scanner.Text()
		for x, c := range text {
			if string(c) == "#" {
				coords[point{x, y}] = make(map[vec2]int)
			}
		}
		y++
	}
	// Should make this look at only n^2/2 items but annoying to do with map
	for c1, _ := range coords {
		for c2, _ := range coords {
			if c1 == c2 {
				continue
			}
			v := NewVector(c1, c2)
			incrementOrOne(coords[c1], v)
		}
	}
	//fmt.Println(coords)
	maxVisible := 0
	var maxLocation point
	for c, v := range coords {
		//fmt.Println(c, len(v), v)
		l := len(v)
		if l > maxVisible {
			maxLocation = c
			maxVisible = l
		}
	}
	fmt.Println("Asteroid at", maxLocation, "can see", maxVisible, "other asteroids")
}
