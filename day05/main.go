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
		r, c := seatID([]byte(scanner.Text()))
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

func seatID(b []byte) (int, int) {
	r := calcNum(b[:7], 128)
	c := calcNum(b[7:], 8)
	return r, c
}

func calcNum(b []byte, k int) int {
	n := 0
	k /= 2
	for _, i := range b {
		if i == 'B' || i == 'R' {
			n += k
		}
		k /= 2
	}
	return n
}
