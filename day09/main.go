package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	target := 731031916
	buffer := make([]int, 0, 25)
	buffer2 := []int{}
	noPart1 := true

	ctr := 0
	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		if noPart1 {
			if lineCount < 25 {
				buffer = append(buffer, num)
				lineCount++
			} else {
				if !hasParts(num, buffer) {
					fmt.Println("Part 1:", num)
					noPart1 = false
				} else {
					buffer[ctr] = num
				}
				ctr = (ctr + 1) % len(buffer)
			}
		}
		buffer2 = append(buffer2, num)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sum := 0
	i := 0
	j := 0
	for sum != target {
		if sum < target {
			sum += buffer2[j]
			j++
		} else {
			sum -= buffer2[i]
			i++
		}
	}

	smallest := 99999999
	largest := 0
	for k := i; k < j; k++ {
		c := buffer2[k]
		if c < smallest {
			smallest = c
		}
		if c > largest {
			largest = c
		}
	}
	fmt.Println("Part 2:", smallest+largest)
}

func hasParts(a int, buffer []int) bool {
	for i := 0; i < len(buffer); i++ {
		for j := i + 1; j < len(buffer); j++ {
			if buffer[i]+buffer[j] == a {
				return true
			}
		}
	}
	return false
}
