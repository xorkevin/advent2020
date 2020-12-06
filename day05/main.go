package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
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

	ids := []int{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		n := calcID([]byte(scanner.Text()))
		ids = append(ids, n)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Ints(ids)

	if len(ids) == 0 {
		log.Fatal("No ids")
	}

	fmt.Println("Part 1:", ids[len(ids)-1])

	prev := ids[0] - 1
	for _, i := range ids {
		if i-prev > 1 {
			fmt.Println("Part 2:", i-1)
			break
		}
		prev = i
	}
}

func calcID(b []byte) int {
	n := 0
	for _, i := range b {
		n *= 2
		if i == 'B' || i == 'R' {
			n++
		}
	}
	return n
}
