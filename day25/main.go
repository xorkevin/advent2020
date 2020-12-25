package main

import (
	"fmt"
)

const (
	puzzleInput1 = 10705932
	puzzleInput2 = 12301431
)

func main() {
	key2 := findLoopSize(puzzleInput2, 7)
	fmt.Println("Part 1:", deriveKey(key2, puzzleInput1))
}

func deriveKey(key int, subj int) int {
	val := 1
	for i := 0; i < key; i++ {
		val = deriveIter(val, subj)
	}
	return val
}

func deriveIter(val int, subj int) int {
	return (val * subj) % 20201227
}

func findLoopSize(target int, subj int) int {
	val := 1
	idx := 0
	for val != target {
		val = deriveIter(val, subj)
		idx++
	}
	return idx
}
