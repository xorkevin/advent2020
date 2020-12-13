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
	nums := []int{}
	rem := []int{}

	for n, i := range strings.Split(puzzleInput2, ",") {
		if i == "x" {
			continue
		}
		num, err := strconv.Atoi(i)
		if err != nil {
			log.Fatal(err)
		}
		nums = append(nums, num)
		rem = append(rem, num-n)
	}

	for i := puzzleInput; ; i++ {
		k, ok := canTake(i, nums)
		if ok {
			fmt.Println("Part 1:", k*(i-puzzleInput))
			break
		}
	}

	x, err := crt(nums, rem)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 2:", x)
}

func crt(nums, rem []int) (int, error) {
	if len(nums) != len(rem) {
		return 0, fmt.Errorf("Invalid pairs")
	}
	p := 1
	for _, i := range nums {
		p *= i
	}
	k := 0
	for i := range nums {
		n := p / nums[i]
		t, err := mulInv(n, nums[i])
		if err != nil {
			return 0, err
		}
		k += rem[i] * t * n
		k %= p
	}
	return (k + p) % p, nil
}

// mulInv returns t, where a*t = 1 (mod n)
func mulInv(a, n int) (int, error) {
	t0 := 0
	t1 := 1
	r0 := n
	r1 := a
	for r1 != 0 {
		q := r0 / r1
		r0, r1 = r1, r0-q*r1
		t0, t1 = t1, t0-q*t1
	}
	if r0 != 1 {
		return 0, fmt.Errorf("No mul inverse for %d mod %d", a, n)
	}
	if t0 < 0 {
		return t0 + n, nil
	}
	return t0, nil
}

func canTake(t int, nums []int) (int, bool) {
	for _, b := range nums {
		if t%b == 0 {
			return b, true
		}
	}
	return 0, false
}
