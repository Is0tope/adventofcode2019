package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type req map[string]int
type quant struct {
	amount int
	level  int
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

func multiplyReqs(target req, multiplier int) req {
	for k, _ := range target {
		target[k] *= multiplier
	}
	return target
}

func getOreRequirement(bill req) req {
	for {
		// Scan over each item in the bill & find the max level
		maxlevel := 0
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
				multiplier := int(math.Ceil(float64(amount) / float64(formula.amount)))
				fmt.Println(prec, amount, formula.amount, multiplier)
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

func markRequirementLevels(item string, level int) {
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

func getLevel(item string) int {
	if item == "ORE" {
		return 0
	}
	return FORMULAS[item].level
}

func getPrecursorsByLevel(level int) []string {
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
			formula[comp[1]] = amount
			if _, ok := REVERSE[comp[1]]; !ok {
				REVERSE[comp[1]] = make(set)
			}
			REVERSE[comp[1]][target[1]] = struct{}{}
		}
		amount, _ := strconv.Atoi(target[0])
		FORMULAS[target[1]] = quant{amount, 0, formula}
	}
	// Mark the levels
	markRequirementLevels("ORE", 0)
	//fmt.Println(formulas)
	oreReqs := getOreRequirement(req{"FUEL": 1})
	fmt.Println("ORE REQUIRED FOR 1 FUEL:", oreReqs)
	//fmt.Println(FORMULAS)

}
