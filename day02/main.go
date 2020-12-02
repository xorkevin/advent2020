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
		r := strings.Split(fields[0], "-")
		a, err := strconv.Atoi(r[0])
		if err != nil {
			log.Fatal(err)
		}
		b, err := strconv.Atoi(r[1])
		if err != nil {
			log.Fatal(err)
		}
		c := strings.TrimRight(fields[1], ":")[0]
		count := 0
		for _, i := range fields[2] {
			if i == rune(c) {
				count++
			}
		}
		if count >= a && count <= b {
			valid++
		}
		if fields[2][a-1] == c && fields[2][b-1] != c {
			valid2++
		} else if fields[2][a-1] != c && fields[2][b-1] == c {
			valid2++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", valid)
	fmt.Println("Part 2:", valid2)
}
