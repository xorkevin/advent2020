package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	puzzleInput  int64 = 1006605
	puzzleInput2       = "19,x,x,x,x,x,x,x,x,x,x,x,x,37,x,x,x,x,x,883,x,x,x,x,x,x,x,23,x,x,x,x,13,x,x,x,17,x,x,x,x,x,x,x,x,x,x,x,x,x,797,x,x,x,x,x,x,x,x,x,41,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,29"
)

func main() {
	nums := []int64{}
	rem := []int64{}

	for n, i := range strings.Split(puzzleInput2, ",") {
		if i == "x" {
			continue
		}
		num, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		nums = append(nums, num)
		rem = append(rem, num-int64(n))
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

func crt(nums, rem []int64) (int64, error) {
	if len(nums) != len(rem) {
		return 0, fmt.Errorf("Invalid pairs")
	}
	var p int64 = 1
	for _, i := range nums {
		p *= i
	}
	var k int64 = 0
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
func mulInv(a, n int64) (int64, error) {
	var t0 int64 = 0
	var t1 int64 = 1
	var r0 int64 = n
	var r1 int64 = a
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

func canTake(t int64, nums []int64) (int64, bool) {
	for _, i := range nums {
		if t%i == 0 {
			return i, true
		}
	}
	return 0, false
}
