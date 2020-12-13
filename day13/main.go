package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	puzzleInput  = 1006605
	puzzleInput2 = "19,x,x,x,x,x,x,x,x,x,x,x,x,37,x,x,x,x,x,883,x,x,x,x,x,x,x,23,x,x,x,x,13,x,x,x,17,x,x,x,x,x,x,x,x,x,x,x,x,x,797,x,x,x,x,x,x,x,x,x,41,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,29"
)

func main() {
	buses := []int{}
	rem := []int{}

	for n, i := range strings.Split(puzzleInput2, ",") {
		if i == "x" {
			continue
		}
		num, err := strconv.Atoi(i)
		if err != nil {
			log.Fatal(err)
		}
		buses = append(buses, num)
		rem = append(rem, num-n)
	}

	for i := puzzleInput; ; i++ {
		k, ok := canTake(i, buses)
		if ok {
			fmt.Println("Part 1:", k*(i-puzzleInput))
			break
		}
	}

	fmt.Println("Part 2:", crt(buses, rem))
}

func crt(nums, rem []int) int {
	p := 1
	for _, i := range nums {
		p *= i
	}
	k := 0
	for i := range nums {
		t := p / nums[i]
		k += rem[i] * mulInv(t, nums[i]) * t
	}
	return k % p
}

func mulInv(a, b int) int {
	if b == 1 {
		return 1
	}
	b0 := b
	x0 := 0
	x1 := 1
	for a > 1 {
		q := a / b
		a, b = b, a%b
		x1, x0 = x0, x1-q*x0
	}
	if x1 < 0 {
		return x1 + b0
	}
	return x1
}

func canTake(t int, buses []int) (int, bool) {
	for _, b := range buses {
		if canTakeBus(t, b) {
			return b, true
		}
	}
	return 0, false
}

func canTakeBus(t int, b int) bool {
	return t%b == 0
}
