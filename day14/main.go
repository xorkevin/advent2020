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
	var maskPos []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " = ")
		if len(s) != 2 {
			log.Fatal("Invalid line")
		}
		if s[0] == "mask" {
			maskZeros, maskOnes, maskFloats, maskPos = processMask([]byte(s[1]))
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

		for _, i := range maskAddrs(maskPos, (addr&maskFloats)|maskOnes) {
			mem2[i] = num
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

func maskAddrs(pos []int, base int) []int {
	if len(pos) == 0 {
		return []int{base}
	}
	addrs := []int{}
	addrs = append(addrs, maskAddrs(pos[1:], base)...)
	addrs = append(addrs, maskAddrs(pos[1:], base|1<<pos[0])...)
	return addrs
}

func processMask(b []byte) (int, int, int, []int) {
	pos := []int{}
	zeros := 0
	ones := 0
	fls := 0
	for n, i := range b {
		p := len(b) - n - 1
		switch i {
		case '0':
			zeros |= 1 << p
		case '1':
			ones |= 1 << p
		case 'X':
			fls |= 1 << p
			pos = append(pos, p)
		}
	}
	return ^zeros, ones, ^fls, pos
}
