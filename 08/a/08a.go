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
	// Find the right layer
	min := 1000000000
	minlayer := 0
	for i, l := range layers {
		cnt := countDigits(l, 0)
		if cnt < min {
			min = cnt
			minlayer = i
		}
		fmt.Println("layer", i, "has", cnt, "0's")
	}
	fmt.Println("layer", minlayer, "has minimum", min, "0's")
	fmt.Println("DONE: ", countDigits(layers[minlayer], 1)*countDigits(layers[minlayer], 2))
}
