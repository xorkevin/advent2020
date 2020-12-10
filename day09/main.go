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
	scanner := bufio.NewScanner(file)

	r := NewTargetReader(scanner)

	target, err := r.Scan()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", target)

	sum := 0
	buffer := []int{}
	for sum != target {
		if sum < target {
			k, err := r.Next()
			if err != nil {
				log.Fatal(err)
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
}

type (
	TargetReader struct {
		buf     []int
		scanner *bufio.Scanner
	}
)

func NewTargetReader(scanner *bufio.Scanner) *TargetReader {
	return &TargetReader{
		buf:     []int{},
		scanner: scanner,
	}
}

func (r *TargetReader) readLine() (int, error) {
	if !r.scanner.Scan() {
		if err := r.scanner.Err(); err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("End of file")
	}
	num, err := strconv.Atoi(r.scanner.Text())
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (r *TargetReader) Scan() (int, error) {
	preamble := 25
	b := []int{}
	for preamble > 0 {
		k, err := r.readLine()
		if err != nil {
			return 0, err
		}
		b = append(b, k)
		preamble--
	}

	idx := 0
	for {
		k, err := r.readLine()
		if err != nil {
			return 0, err
		}
		b = append(b, k)
		if !hasParts(k, b[idx:idx+25]) {
			break
		}
		idx++
	}

	r.buf = b
	return b[len(b)-1], nil
}

func (r *TargetReader) Next() (int, error) {
	if len(r.buf) > 0 {
		k := r.buf[0]
		r.buf = r.buf[1:]
		return k, nil
	}
	return r.readLine()
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
