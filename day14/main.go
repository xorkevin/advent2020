package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

	mem := map[int]int{}
	mem2 := map[int]int{}

	maskZeros := 0
	maskOnes := 0
	maskFloats := 0
	maskPos := []int{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " = ")
		if len(s) != 2 {
			log.Fatal("Invalid line")
		}
		if s[0] == "mask" {
			maskZeros, maskOnes, maskFloats, maskPos = processMask([]byte(strings.TrimLeft(s[1], "0")), maskPos[:0])
			continue
		}
		addr, err := strconv.Atoi(s[0][4 : len(s[0])-1])
		if err != nil {
			log.Fatal(err)
		}
		num, err := strconv.Atoi(s[1])
		if err != nil {
			log.Fatal(err)
		}

		mem[addr] = (num & maskZeros) | maskOnes

		base := (addr & maskFloats) | maskOnes
		for _, i := range maskPos {
			mem2[base|i] = num
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sum := 0
	for _, v := range mem {
		sum += v
	}
	fmt.Println("Part 1:", sum)

	sum2 := 0
	for _, v := range mem2 {
		sum2 += v
	}
	fmt.Println("Part 2:", sum2)
}

func processMask(b []byte, posMask []int) (int, int, int, []int) {
	pos := []int{}
	zeros := 0
	ones := 0
	fls := 0
	for n, i := range b {
		p := len(b) - n - 1
		k := 1 << p
		switch i {
		case '0':
			zeros |= k
		case '1':
			ones |= k
		case 'X':
			fls |= k
			pos = append(pos, k)
		}
	}
	return ^zeros, ones, ^fls, maskDfs(0, pos, posMask)
}

func maskDfs(base int, pos []int, posMask []int) []int {
	if len(pos) == 0 {
		return append(posMask, base)
	}
	posMask = maskDfs(base, pos[1:], posMask)
	return maskDfs(base|pos[0], pos[1:], posMask)
}
