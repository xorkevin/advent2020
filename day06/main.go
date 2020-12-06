package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	count := 0
	s := map[byte]struct{}{}

	count2 := 0
	k := map[byte]int{}
	people := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			count += len(s)
			s = map[byte]struct{}{}
			for _, v := range k {
				if v == people {
					count2++
				}
			}
			k = map[byte]int{}
			people = 0
			continue
		}
		people++
		for _, i := range []byte(line) {
			s[i] = struct{}{}
			if _, ok := k[i]; !ok {
				k[i] = 0
			}
			k[i]++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	count += len(s)

	for _, v := range k {
		if v == people {
			count2++
		}
	}

	fmt.Println("Part 1:", count)
	fmt.Println("Part 2:", count2)
}
