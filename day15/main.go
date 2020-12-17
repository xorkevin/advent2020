package main

import (
	"fmt"
)

func main() {
	puzzleInput := []int{19, 20, 14, 0, 9, 1}
	boundary := 1 << 24
	hist := make([]int, boundary)
	hist2 := map[int]int{}
	for n, i := range puzzleInput[:len(puzzleInput)-1] {
		if i < boundary {
			hist[i] = n + 1
		} else {
			hist2[i] = n + 1
		}
	}
	idx := len(puzzleInput)
	prev := puzzleInput[len(puzzleInput)-1]
	for idx < 2020 {
		var v int
		if prev < boundary {
			v = hist[prev]
			hist[prev] = idx
		} else {
			v = hist2[prev]
			hist2[prev] = idx
		}
		if v == 0 {
			prev = 0
		} else {
			prev = idx - v
		}
		idx++
	}
	fmt.Println("Part 1:", prev)
	for idx < 30_000_000 {
		var v int
		if prev < boundary {
			v = hist[prev]
			hist[prev] = idx
		} else {
			v = hist2[prev]
			hist2[prev] = idx
		}
		if v == 0 {
			prev = 0
		} else {
			prev = idx - v
		}
		idx++
	}
	fmt.Println("Part 2:", prev)
}
