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

	re := regexp.MustCompile(`^(\d+)-(\d+) ([a-z]): ([a-z]+)$`)

	valid := 0
	valid2 := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		matches := re.FindStringSubmatch(scanner.Text())
		if matches == nil {
			log.Fatal("Failed to match line")
		}
		a, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Fatal(err)
		}
		b, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatal(err)
		}
		c := matches[3][0]
		pass := []byte(matches[4])

		count := 0
		for _, i := range pass {
			if i == c {
				count++
			}
		}
		if count >= a && count <= b {
			valid++
		}
		l := len(pass)
		if !inBounds(a-1, 0, l-1) || !inBounds(b-1, 0, l-1) {
			log.Fatal("Invalid indicies")
		}
		c1 := pass[a-1]
		c2 := pass[b-1]
		if c1 != c2 && (c1 == c || c2 == c) {
			valid2++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", valid)
	fmt.Println("Part 2:", valid2)
}

func inBounds(a, l, h int) bool {
	return a >= l && a <= h
}
