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

	grid := [][]byte{}
	grid2 := [][]byte{}
	tmpGrid := [][]byte{}
	tmpGrid2 := [][]byte{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		k := []byte(scanner.Text())
		grid = append(grid, k)
		tmpGrid = append(tmpGrid, make([]byte, len(k)))
		k2 := []byte(scanner.Text())
		grid2 = append(grid2, k2)
		tmpGrid2 = append(tmpGrid2, make([]byte, len(k2)))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	untilStable(false, grid, tmpGrid)
	fmt.Println("Part 1:", countSeats(grid))

	untilStable(true, grid2, tmpGrid2)
	fmt.Println("Part 2:", countSeats(grid2))
}

func countSeats(grid [][]byte) int {
	count := 0
	for _, i := range grid {
		for _, j := range i {
			if j == '#' {
				count++
			}
		}
	}
	return count
}

func untilStable(mode bool, grid, tmpGrid [][]byte) {
	for {
		if !nextGrid(mode, grid, tmpGrid) {
			break
		}
		tmp := grid
		grid = tmpGrid
		tmpGrid = tmp
	}
}

func nextGrid(mode bool, grid [][]byte, next [][]byte) bool {
	change := false
	for i, r := range grid {
		for j, c := range r {
			var k byte
			if mode {
				k = nextSeat2(i, j, grid)
			} else {
				k = nextSeat(i, j, grid)
			}
			if k != c {
				change = true
			}
			next[i][j] = k
		}
	}
	return change
}

func nextSeat2(i, j int, grid [][]byte) byte {
	c := grid[i][j]
	s1 := seatVal2(i, j, -1, 0, grid)
	s2 := seatVal2(i, j, -1, 1, grid)
	s3 := seatVal2(i, j, 0, 1, grid)
	s4 := seatVal2(i, j, 1, 1, grid)
	s5 := seatVal2(i, j, 1, 0, grid)
	s6 := seatVal2(i, j, 1, -1, grid)
	s7 := seatVal2(i, j, 0, -1, grid)
	s8 := seatVal2(i, j, -1, -1, grid)
	k := s1 + s2 + s3 + s4 + s5 + s6 + s7 + s8
	if c == 'L' && k == 0 {
		return '#'
	}
	if c == '#' && k >= 5 {
		return 'L'
	}
	return c
}

func seatVal2(i, j int, di, dj int, grid [][]byte) int {
	i += di
	j += dj
	for inBounds(i, 0, len(grid)-1) && inBounds(j, 0, len(grid[0])-1) {
		if grid[i][j] == 'L' {
			return 0
		}
		if grid[i][j] == '#' {
			return 1
		}
		i += di
		j += dj
	}
	return 0
}

func nextSeat(i, j int, grid [][]byte) byte {
	c := grid[i][j]
	s1 := seatVal(i-1, j, grid)
	s2 := seatVal(i-1, j+1, grid)
	s3 := seatVal(i, j+1, grid)
	s4 := seatVal(i+1, j+1, grid)
	s5 := seatVal(i+1, j, grid)
	s6 := seatVal(i+1, j-1, grid)
	s7 := seatVal(i, j-1, grid)
	s8 := seatVal(i-1, j-1, grid)
	k := s1 + s2 + s3 + s4 + s5 + s6 + s7 + s8
	if c == 'L' && k == 0 {
		return '#'
	}
	if c == '#' && k >= 4 {
		return 'L'
	}
	return c
}

func seatVal(i, j int, grid [][]byte) int {
	if !inBounds(i, 0, len(grid)-1) {
		return 0
	}
	if !inBounds(j, 0, len(grid[0])-1) {
		return 0
	}
	if grid[i][j] == '#' {
		return 1
	}
	return 0
}

func inBounds(a, l, r int) bool {
	return a >= l && a <= r
}
