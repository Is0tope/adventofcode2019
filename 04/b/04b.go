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
	groups := make(map[int]int)
	for i := 0; i < len(str); i++ {
		d, _ := strconv.Atoi(str[i : i+1])
		if d < last {
			return false
		}
		if d == last {
			_, ok := groups[d]
			if ok {
				groups[d]++
			} else {
				groups[d] = 2
			}
		}
		last = d
	}
	for _, v := range groups {
		if v == 2 {
			return true
		}
	}
	return false
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
