package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

const (
	puzzleInput = "input.txt"
)

func main() {
	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	re := regexp.MustCompile(`^([a-z ]+) bags contain`)
	re2 := regexp.MustCompile(`([0-9]+) ([a-z ]+) bags?[,.]`)
	re3 := regexp.MustCompile(`no other bags.$`)

	graph := map[string]map[string]int{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		m := re.FindStringSubmatch(line)
		if m == nil {
			log.Fatal("Invalid line format")
		}
		m2 := re2.FindAllStringSubmatch(line, -1)
		m3 := re3.MatchString(line)
		graph[m[1]] = map[string]int{}
		if m3 {
			continue
		}
		if m2 == nil {
			log.Fatal("Invalid line format")
		}
		for _, i := range m2 {
			num, err := strconv.Atoi(i[1])
			if err != nil {
				log.Fatal(err)
			}
			graph[m[1]][i[2]] = num
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	count := 0
	cache := map[string]bool{}
	for k := range graph {
		if hasPathTo(k, "shiny gold", graph, cache) {
			count++
		}
	}

	fmt.Println("Part 1:", count)

	cache2 := map[string]int{}
	fmt.Println("Part 2:", countChildren("shiny gold", graph, cache2)-1)
}

func hasPathTo(a, b string, graph map[string]map[string]int, cache map[string]bool) bool {
	if v, ok := cache[a]; ok {
		return v
	}
	for k := range graph[a] {
		if b == k {
			cache[a] = true
			return true
		}
	}
	for k := range graph[a] {
		if hasPathTo(k, b, graph, cache) {
			cache[a] = true
			return true
		}
	}
	cache[a] = false
	return false
}

func countChildren(a string, graph map[string]map[string]int, cache map[string]int) int {
	if v, ok := cache[a]; ok {
		return v
	}
	count := 1
	for k, v := range graph[a] {
		count += v * countChildren(k, graph, cache)
	}
	cache[a] = count
	return count
}
