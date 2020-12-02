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

	valid := 0
	valid2 := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 3 {
			log.Fatal("Invalid line format")
		}
		r := strings.Split(fields[0], "-")
		if len(r) != 2 {
			log.Fatal("Invalid range")
		}
		a, err := strconv.Atoi(r[0])
		if err != nil {
			log.Fatal(err)
		}
		b, err := strconv.Atoi(r[1])
		if err != nil {
			log.Fatal(err)
		}
		ch := strings.TrimRight(fields[1], ":")
		if len(ch) != 1 {
			log.Fatal("Invalid char")
		}
		c := ch[0]
		pass := []byte(fields[2])

		count := 0
		for _, i := range pass {
			if i == c {
				count++
			}
		}
		if count >= a && count <= b {
			valid++
		}
		l := len(pass)
		if !inBounds(a-1, 0, l-1) || !inBounds(b-1, 0, l-1) {
			log.Fatal("Invalid indicies")
		}
		c1 := pass[a-1]
		c2 := pass[b-1]
		if c1 != c2 && (c1 == c || c2 == c) {
			valid2++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", valid)
	fmt.Println("Part 2:", valid2)
}

func inBounds(a, l, h int) bool {
	return a >= l && a <= h
}
