package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type entity struct {
	name  string
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
		if val, ok := orbits[tokens[0]]; ok {
			orbits[tokens[0]] = append(val, tokens[1])
		} else {
			orbits[tokens[0]] = []string{tokens[1]}
		}
	}
	fmt.Println(orbits)
	// Set up queue
	queue := make(chan entity, 1000)
	queue <- entity{"COM", 0}
	// BFS
	ok := true
	orbitCount := 0
	for ok {
		select {
		case ent := <-queue:
			fmt.Println("processing ", ent)
			orbiters := orbits[ent.name]
			orbitCount += ent.depth
			for _, o := range orbiters {
				queue <- entity{o, ent.depth + 1}
			}
		default:
			fmt.Println("NO MORE ENTITIES")
			ok = false
		}
	}
	fmt.Printf("DONE: %d\n", orbitCount)
}
