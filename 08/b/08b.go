package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var width = 25
var height = 6

type layer [][]int
type image []layer

func coords(i int) (int, int, int) {
	i2 := i % (width * height)
	l := (i - i2) / (width * height)
	x := i2 % width
	y := (i2 - x) / width
	return l, x, y
}

func makeLayer() layer {
	l := make([][]int, height)
	for i, _ := range l {
		l[i] = make([]int, width)
	}
	return l
}

func countDigits(l layer, d int) int {
	cnt := 0
	for _, row := range l {
		for _, cell := range row {
			if cell == d {
				cnt++
			}
		}
	}
	return cnt
}

func getColor(pixels []int) int {
	last := pixels[0]
	for _, p := range pixels {
		if p == 2 {
			continue
		}
		last = p
	}
	return last
}

func main() {
	file, _ := os.Open("../input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	text := scanner.Text()

	// Initialise
	layercnt := len(text) / (width * height)
	layers := make([]layer, layercnt)
	for i, _ := range layers {
		layers[i] = makeLayer()
	}
	// Fill in
	for i, val := range text {
		l, x, y := coords(i)
		num, _ := strconv.Atoi(string(val))
		layers[l][y][x] = num
	}
	fmt.Println(layers)
	// Find the pixel colors
	fmt.Println("IMAGE: ")
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			stack := make([]int, 0)
			for l := len(layers) - 1; l >= 0; l-- {
				stack = append(stack, layers[l][i][j])
			}
			color := getColor(stack)
			if color == 1 {
				fmt.Printf("x")
			} else {
				fmt.Printf(" ")
			}

		}
		fmt.Printf("\n")
	}
}
