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

	var mask []byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " = ")
		if len(s) != 2 {
			log.Fatal("Invalid line")
		}
		if s[0] == "mask" {
			mask = []byte(s[1])
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

		mem[addr] = applyMask(num, mask)

		for _, i := range maskAddrs(orMask(addr, mask)) {
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

func maskAddrs(b []byte) []int {
	base := 0
	for n, i := range b {
		p := len(b) - n - 1
		if i == '1' {
			base = setBit(base, p)
		}
	}
	pos := []int{}
	for n, i := range b {
		if i != 'X' {
			continue
		}
		pos = append(pos, len(b)-n-1)
	}
	return maskAddrsRec(pos, base)
}

func maskAddrsRec(pos []int, base int) []int {
	if len(pos) == 0 {
		return []int{base}
	}
	addrs := []int{}
	addrs = append(addrs, maskAddrsRec(pos[1:], base)...)
	addrs = append(addrs, maskAddrsRec(pos[1:], setBit(base, pos[0]))...)
	return addrs
}

func orMask(v int, b []byte) []byte {
	next := make([]byte, len(b))
	for n, i := range b {
		p := len(b) - n - 1
		switch i {
		case '0':
			if hasBit(v, p) {
				next[n] = '1'
			} else {
				next[n] = '0'
			}
		case '1':
			next[n] = '1'
		case 'X':
			next[n] = 'X'
		}
	}
	return next
}

func applyMask(v int, b []byte) int {
	for n, i := range b {
		p := len(b) - n - 1
		switch i {
		case '0':
			v = clearBit(v, p)
		case '1':
			v = setBit(v, p)
		}
	}
	return v
}

func setBit(n int, pos int) int {
	return n | (1 << pos)
}

func clearBit(n int, pos int) int {
	return n & ^(1 << pos)
}

func hasBit(n int, pos int) bool {
	val := n & (1 << pos)
	return (val > 0)
}
