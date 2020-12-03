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

	m := [][]byte{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m = append(m, []byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(getCount(3, 1, m))
	fmt.Println(getCount(1, 1, m) * getCount(3, 1, m) * getCount(5, 1, m) * getCount(7, 1, m) * getCount(1, 2, m))
}

func getCount(x, y int, m [][]byte) int {
	count := 0
	i := 0
	j := 0
	for i < len(m) {
		if m[i][j] == '#' {
			count++
		}
		i += y
		j = (j + x) % len(m[0])
	}
	return count
}
