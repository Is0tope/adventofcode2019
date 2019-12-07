package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type entity struct {
	name  string
	from  string
	depth int
}

func findOrbits(orbits []string, c chan string) {
	// BFS
	for ent := range c {
		fmt.Println(ent)
	}
}

func main() {
	file, _ := os.Open("../input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// Mappings. orbited -> orbitee
	orbits := make(map[string][]string)
	for scanner.Scan() {
		text := scanner.Text()
		tokens := strings.Split(text, ")")
		// Do it both ways
		if val, ok := orbits[tokens[0]]; ok {
			orbits[tokens[0]] = append(val, tokens[1])
		} else {
			orbits[tokens[0]] = []string{tokens[1]}
		}
		if val, ok := orbits[tokens[1]]; ok {
			orbits[tokens[1]] = append(val, tokens[0])
		} else {
			orbits[tokens[1]] = []string{tokens[0]}
		}

	}
	fmt.Println(orbits)
	// Set up queue
	queue := make(chan entity, 1000)
	queue <- entity{"YOU", "", 0}
	// BFS
	ok := true
	distance := 0
	for ok {
		select {
		case ent := <-queue:
			fmt.Println("processing ", ent)
			if ent.name == "SAN" {
				fmt.Println("FOUND SANTA @ : ", ent.depth-2)
				ok = false
				continue
			}
			orbiters := orbits[ent.name]
			distance += ent.depth
			for _, o := range orbiters {
				// Don't go back the way you came
				if o == ent.from {
					continue
				}
				queue <- entity{o, ent.name, ent.depth + 1}
			}
		default:
			fmt.Println("NO MORE ENTITIES")
			ok = false
		}
	}
}
