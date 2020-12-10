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

	nums := []int{0}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		nums = append(nums, num)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Ints(nums)
	target := nums[len(nums)-1] + 3
	nums = append(nums, target)

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

	count := countPaths(nums)
	fmt.Println("Part 2:", count)
}

func countPaths(nums []int) int {
	cache := make([]int, len(nums))
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

func getVal(start int, n int, nums []int, cache []int) int {
	if n < 0 || n >= len(cache) {
		return 0
	}
	if start-nums[n] > 3 {
		return 0
	}
	return cache[n]
}
