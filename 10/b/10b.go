package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type point struct {
	x, y int
}

type distpoint struct {
	x, y     int
	distance float64
}

type vec2 struct {
	x, y  int
	angle float64
}

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

func rad2deg(x float64) float64 {
	if x < 0 {
		x = math.Pi + (math.Pi + x)
	}
	return (x * 180 / math.Pi)
}

func distance(a, b point) float64 {
	dx := math.Abs(float64(a.x - b.x))
	dy := math.Abs(float64(a.y - b.y))
	return math.Sqrt((dx * dx) + (dy * dy))
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
	// Find angle using atan2, switch x & y coords in order to give angular origin at y axis
	at := rad2deg(math.Atan2(float64(vx), float64(-1*vy)))
	return vec2{vx, vy, at}
}

func main() {
	file, _ := os.Open("../input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	STATION := point{26, 29}

	vectors := make(map[vec2][]distpoint)
	y := 0
	for scanner.Scan() {
		text := scanner.Text()
		for x, c := range text {
			if string(c) == "#" {
				p := point{x, y}
				pd := distpoint{x, y, distance(p, STATION)}
				if p == STATION {
					continue
				}
				v := NewVector(STATION, p)
				if list, ok := vectors[v]; ok {
					vectors[v] = append(list, pd)
				} else {
					vectors[v] = []distpoint{pd}
				}
			}
		}
		y++
	}
	// Presort the distances
	for _, v := range vectors {
		sort.Slice(v, func(i, j int) bool {
			return v[i].distance < v[j].distance
		})
	}
	// Get the ordering of the vectors
	keys := []vec2{}
	for k, _ := range vectors {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].angle < keys[j].angle
	})
	// Start spinning
	counter := 0
	destroyed := 0
	for destroyed < 200 {
		asteroids := vectors[keys[counter%len(keys)]]
		if len(asteroids) == 0 {
			counter++
			continue
		}
		tokill := asteroids[0]
		vectors[keys[counter%len(keys)]] = asteroids[1:]
		fmt.Println("asteroid", destroyed+1, tokill, "vector", keys[counter%len(keys)])
		destroyed++
		counter++
	}
}
