package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type req map[string]int64
type quant struct {
	amount int64
	level  int64
	reqs   req
}
type set map[string]struct{}

var FORMULAS map[string]quant
var REVERSE map[string]set

func mergeReqs(target req, source req) req {
	for k, v := range source {
		if _, ok := target[k]; ok {
			target[k] += v
		} else {
			target[k] = v
		}
	}
	return target
}

func multiplyReqs(target req, multiplier int64) req {
	ret := make(req)
	for k, _ := range target {
		ret[k] = target[k] * multiplier
	}
	return ret
}

func getOreRequirement(bill req) req {
	for {
		// Scan over each item in the bill & find the max level
		maxlevel := int64(0)
		for prec, _ := range bill {
			if prec == "ORE" {
				continue
			}
			lvl := getLevel(prec)
			if lvl > maxlevel {
				maxlevel = lvl
			}
		}
		// Exit out if done
		if maxlevel == 0 {
			return bill
		}
		// Find all precursors that are at the highest level, and break them down
		// If not at highest level just keep them there
		newbill := make(req)
		for prec, amount := range bill {
			lvl := getLevel(prec)
			if lvl == maxlevel {
				formula := FORMULAS[prec]
				multiplier := int64(math.Ceil(float64(amount) / float64(formula.amount)))
				//fmt.Println(prec, amount, formula.amount, multiplier)
				r := multiplyReqs(formula.reqs, multiplier)
				newbill = mergeReqs(newbill, r)
			} else {
				newbill = mergeReqs(newbill, req{prec: amount})
			}
		}
		// Overwrite and loop
		bill = newbill
	}
}

func markRequirementLevels(item string, level int64) {
	lvl := getLevel(item)
	if lvl < level {
		f := FORMULAS[item]
		f.level = level
		FORMULAS[item] = f
	}
	children := REVERSE[item]
	for k, _ := range children {
		markRequirementLevels(k, level+1)
	}
}

func getLevel(item string) int64 {
	if item == "ORE" {
		return 0
	}
	return FORMULAS[item].level
}

func getPrecursorsByLevel(level int64) []string {
	ret := []string{}
	for k, v := range FORMULAS {
		if v.level == level {
			ret = append(ret, k)
		}
	}
	return ret
}
func main() {
	file, _ := os.Open("../input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// Get all the requirements
	FORMULAS = make(map[string]quant)
	REVERSE = make(map[string]set)
	for scanner.Scan() {
		text := scanner.Text()
		strs := strings.Split(text, " => ")
		target := strings.Split(strs[1], " ")
		reqs := strings.Split(strs[0], ", ")
		formula := make(req)
		for _, r := range reqs {
			comp := strings.Split(r, " ")
			amount, _ := strconv.Atoi(comp[0])
			formula[comp[1]] = int64(amount)
			if _, ok := REVERSE[comp[1]]; !ok {
				REVERSE[comp[1]] = make(set)
			}
			REVERSE[comp[1]][target[1]] = struct{}{}
		}
		amount, _ := strconv.Atoi(target[0])
		FORMULAS[target[1]] = quant{int64(amount), int64(0), formula}
	}
	// Mark the levels
	markRequirementLevels("ORE", 0)

	// Educated brute force
	// Unsure if ore use is strictly sorted with respect to fuel requirement, so wary of using binary search.
	cargo := int64(1000000000000)
	fuel := int64(1)
	step := int64(100000)
	for {
		ore := getOreRequirement(req{"FUEL": fuel})["ORE"]
		fmt.Println("Used", ore, "ore to make", fuel, "fuel")
		if ore > cargo {
			fuel -= step
			break
		}
		fuel += step
	}
	fmt.Println("Switching to more accurate tracking")
	step = 1
	for {
		ore := getOreRequirement(req{"FUEL": fuel})["ORE"]
		fmt.Println("Used", ore, "ore to make", fuel, "fuel")
		if ore > cargo {
			fuel -= step
			break
		}
		fuel += step
	}
	fmt.Println("MAX FUEL GENERATED FROM CARGO:", fuel)

}
