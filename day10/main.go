package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
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

	nums := []uint64{0}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		num, err := strconv.ParseUint(scanner.Text(), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		nums = append(nums, num)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
	nums = append(nums, nums[len(nums)-1]+3)

	start := nums[0]

	diff1 := 0
	diff3 := 0
	for _, i := range nums[1:] {
		k := i - start
		if k == 1 {
			diff1++
		} else if k == 3 {
			diff3++
		}
		start = i
	}
	fmt.Println("Part 1:", diff1*diff3)
	fmt.Println("Part 2:", countPaths(nums))
}

func countPaths(nums []uint64) uint64 {
	cache := make([]uint64, len(nums))
	for n, i := range nums {
		if n == 0 {
			cache[n] = 1
			continue
		}
		a := getVal(i, n-1, nums, cache)
		b := getVal(i, n-2, nums, cache)
		c := getVal(i, n-3, nums, cache)
		cache[n] = a + b + c
	}
	return cache[len(cache)-1]
}

func getVal(start uint64, n int, nums []uint64, cache []uint64) uint64 {
	if n < 0 || n >= len(cache) {
		return 0
	}
	if start-nums[n] > 3 {
		return 0
	}
	return cache[n]
}
