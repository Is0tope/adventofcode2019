package main

import (
	"fmt"
	"strconv"
)

var lo int = 367479
var hi int = 893698

func validate(num int) bool {
	str := strconv.Itoa(num)
	last := 0
	isdouble := false
	for i := 0; i < len(str); i++ {
		d, _ := strconv.Atoi(str[i : i+1])
		if d < last {
			return false
		}
		if d == last {
			isdouble = true
		}
		last = d
	}
	return isdouble
}

func main() {
	// Brute force, why not
	counter := 0
	for i := lo; i <= hi; i++ {
		if validate(i) {
			counter++
		}
	}
	fmt.Printf("DONE: %d\n", counter)
}
