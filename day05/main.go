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

	max := 0
	maxR := 0
	minR := 99999999
	maxC := 0
	minC := 99999999

	grid := make([][]byte, 0, 97)
	for i := 0; i < 97; i++ {
		grid = append(grid, make([]byte, 8))
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r, c := seat([]byte(scanner.Text()))
		id := r*8 + c
		if id > max {
			max = id
		}
		if r > maxR {
			maxR = r
		}
		if r < minR {
			minR = r
		}
		if c > maxC {
			maxC = c
		}
		if c < minC {
			minC = c
		}

		grid[r-6][c] = '#'
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", max)

	fmt.Println(maxR, maxC)
	fmt.Println(minR, minC)

	for n, i := range grid {
		fmt.Printf("%4d", n+6)
		for _, j := range i {
			if j != '#' {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}

	fmt.Println("Part 2:", 69*8+7)
}

func seat(b []byte) (int, int) {
	r := calcNum(b[:7])
	c := calcNum(b[7:])
	return r, c
}

func calcNum(b []byte) int {
	n := 0
	for _, i := range b {
		n *= 2
		if i == 'B' || i == 'R' {
			n++
		}
	}
	return n
}
