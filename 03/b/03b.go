package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type M map[string][]int64

func getLinePoints(instructions []string) M {
	points := make(M)
	x, y, travel := int64(0), int64(0), int64(1)
	for _, t := range instructions {
		dir := t[:1]
		dist, _ := strconv.ParseInt(t[1:], 0, 0)
		dist = int64(dist)
		if dir == "U" || dir == "D" {
			for i := int64(1); i <= dist; i++ {
				newy := y
				if dir == "U" {
					newy += i
				} else {
					newy -= i
				}
				key := fmt.Sprintf("(%d,%d)", x, newy)
				if _, ok := points[key]; !ok {
					points[key] = []int64{x, newy, travel}
				}
				travel += 1
			}
			if dir == "U" {
				y += dist
			} else {
				y -= dist
			}
		}
		if dir == "R" || dir == "L" {
			for i := int64(1); i <= dist; i++ {
				newx := x
				if dir == "R" {
					newx += i
				} else {
					newx -= i
				}
				key := fmt.Sprintf("(%d,%d)", newx, y)
				if _, ok := points[key]; !ok {
					points[key] = []int64{newx, y, travel}
				}
				travel += 1
			}
			if dir == "R" {
				x += dist
			} else {
				x -= dist
			}
		}

	}
	return points
}

func manhattanDist(a []int64, b []int64) int64 {
	return int64(math.Abs(float64(a[0]-b[0])) + math.Abs(float64(a[1]-b[1])))
}

func main() {
	file, _ := os.Open("../input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var lines []M
	for scanner.Scan() {
		text := scanner.Text()
		tokens := strings.Split(text, ",")
		line := getLinePoints(tokens)
		lines = append(lines, line)
	}
	line1 := lines[0]
	line2 := lines[1]

	mintravel := int64(math.MaxInt64)
	var minpoint []int64
	for k, v := range line1 {
		if v2, ok := line2[k]; ok {
			newtravel := v[2] + v2[2]
			if newtravel < mintravel {
				mintravel = newtravel
				minpoint = v[:2]
			}
		}
	}
	fmt.Println("DONE: ", minpoint, mintravel)
}
