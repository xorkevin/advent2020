package main

import (
	"bufio"
	"context"
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
	scanner := bufio.NewScanner(file)

	ctx, cancel := context.WithCancel(context.Background())

	numStream := make(chan int)
	doneReadNums := make(chan struct{})
	go readNums(ctx, doneReadNums, scanner, numStream)
	targetStream := make(chan int)
	throughStream := make(chan int)
	doneScanTarget := make(chan struct{})
	go scanTarget(ctx, doneScanTarget, numStream, targetStream, throughStream)

	target := <-targetStream
	fmt.Println("Part 1:", target)

	sum := 0
	buffer := []int{}
	for sum != target {
		if sum < target {
			k, ok := <-throughStream
			if !ok {
				break
			}
			sum += k
			buffer = append(buffer, k)
		} else {
			if len(buffer) == 0 {
				break
			}
			sum -= buffer[0]
			buffer = buffer[1:]
		}
	}
	if sum == target {
		min, max, ok := minMax(buffer)
		if ok {
			fmt.Println("Part 2:", min+max)
		}
	}

	cancel()
	<-doneReadNums
	<-doneScanTarget
}

func minMax(b []int) (int, int, bool) {
	if len(b) == 0 {
		return 0, 0, false
	}
	min := b[0]
	max := b[0]
	for _, i := range b {
		if i < min {
			min = i
		}
		if i > max {
			max = i
		}
	}
	return min, max, true
}

func readNums(ctx context.Context, done chan<- struct{}, scanner *bufio.Scanner, s chan<- int) {
	defer close(done)
	defer close(s)

	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Print(err)
			return
		}
		select {
		case <-ctx.Done():
			return
		case s <- num:
		}
	}
	if err := scanner.Err(); err != nil {
		log.Print(err)
	}
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

func scanTarget(ctx context.Context, done chan<- struct{}, s <-chan int, target chan<- int, through chan<- int) {
	defer close(done)
	defer close(through)
	defer close(target)

	{
		preamble := 25
		buffer := []int{}
		for preamble > 0 {
			k, ok := <-s
			if !ok {
				return
			}
			buffer = append(buffer, k)
			preamble--
		}

		idx := 0
		for {
			k, ok := <-s
			if !ok {
				return
			}
			buffer = append(buffer, k)
			if !hasParts(k, buffer[idx:idx+25]) {
				select {
				case <-ctx.Done():
					return
				case target <- k:
				}
				break
			}
			idx++
		}

		for _, i := range buffer {
			select {
			case <-ctx.Done():
				return
			case through <- i:
			}
		}
	}

	for {
		k, ok := <-s
		if !ok {
			return
		}
		select {
		case <-ctx.Done():
			return
		case through <- k:
		}
	}
}
