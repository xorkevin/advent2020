package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	m := map[int]struct{}{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		k := 2020 - num
		if _, ok := m[k]; ok {
			fmt.Println("Part 1:", k*num)
		}
		m[num] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
pt2loop:
	for i := range m {
		for j := range m {
			k := 2020 - i - j
			if _, ok := m[k]; ok {
				fmt.Println("Part 2:", i*j*k)
				break pt2loop
			}
		}
	}
}
