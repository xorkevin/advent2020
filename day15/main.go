package main

import (
	"fmt"
)

func main() {
	puzzleInput := []int{19, 20, 14, 0, 9, 1}
	hist := make([]int, 30_000_000)
	for n, i := range puzzleInput[:len(puzzleInput)-1] {
		hist[i] = n + 1
	}
	idx := len(puzzleInput)
	prev := puzzleInput[len(puzzleInput)-1]
	for idx < 2020 {
		v := hist[prev]
		hist[prev] = idx
		if v == 0 {
			prev = 0
		} else {
			prev = idx - v
		}
		idx++
	}
	fmt.Println("Part 1:", prev)
	for idx < 30_000_000 {
		v := hist[prev]
		hist[prev] = idx
		if v == 0 {
			prev = 0
		} else {
			prev = idx - v
		}
		idx++
	}
	fmt.Println("Part 2:", prev)
}
