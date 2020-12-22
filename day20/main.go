package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

const (
	puzzleInput = "input.txt"
)

type (
	Edge struct {
		f, b int
	}

	Tile struct {
		id         int
		r, t, l, b Edge
	}
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

	tiles := []Tile{}
	tileID := 0
	grid := [][]byte{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if tileID == 0 {
			num, err := strconv.Atoi(line[5 : len(line)-1])
			if err != nil {
				log.Fatal(err)
			}
			tileID = num
			continue
		}
		if line == "" {
			h := len(grid)
			w := len(grid[0])
			rf := 0
			rb := 0
			lf := 0
			lb := 0
			for i := 0; i < h; i++ {
				rf <<= 1
				if grid[i][w-1] == '#' {
					rf += 1
				}
				rb <<= 1
				if grid[h-i-1][w-1] == '#' {
					rb += 1
				}
				lf <<= 1
				if grid[h-i-1][0] == '#' {
					lf += 1
				}
				lb <<= 1
				if grid[i][0] == '#' {
					lb += 1
				}
			}
			tf := 0
			tb := 0
			bf := 0
			bb := 0
			for i := 0; i < w; i++ {
				tf <<= 1
				if grid[0][i] == '#' {
					tf += 1
				}
				tb <<= 1
				if grid[0][w-i-1] == '#' {
					tb += 1
				}
				bf <<= 1
				if grid[h-1][w-i-1] == '#' {
					bf += 1
				}
				bb <<= 1
				if grid[h-1][i] == '#' {
					bb += 1
				}
			}
			tiles = append(tiles, Tile{
				id: tileID,
				r:  Edge{f: rf, b: rb},
				t:  Edge{f: tf, b: tb},
				l:  Edge{f: lf, b: lb},
				b:  Edge{f: bf, b: bb},
			})
			tileID = 0
			grid = [][]byte{}
			continue
		}
		grid = append(grid, []byte(line))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	squareSize := math.Sqrt(float64(len(tiles)))
	if math.Floor(squareSize) != squareSize {
		log.Fatal("Not a square")
	}

	fmt.Println(squareSize)
}
