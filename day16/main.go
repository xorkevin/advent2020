package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	Rule struct {
		n1 int
		n2 int
		n3 int
		n4 int
	}
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

	re := regexp.MustCompile(`^([a-z ]+): ([0-9]+)-([0-9]+) or ([0-9]+)-([0-9]+)$`)

	rules := map[string]Rule{}
	ruleNames := map[string]struct{}{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		m := re.FindStringSubmatch(line)
		if m == nil {
			log.Fatal("Invalid rule line")
		}
		text := m[1]
		n1, err := strconv.Atoi(m[2])
		if err != nil {
			log.Fatal(err)
		}
		n2, err := strconv.Atoi(m[3])
		if err != nil {
			log.Fatal(err)
		}
		n3, err := strconv.Atoi(m[4])
		if err != nil {
			log.Fatal(err)
		}
		n4, err := strconv.Atoi(m[5])
		if err != nil {
			log.Fatal(err)
		}
		rules[text] = Rule{
			n1: n1,
			n2: n2,
			n3: n3,
			n4: n4,
		}
		ruleNames[text] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	ownTicket := []int{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		if line == "your ticket:" {
			continue
		}
		for _, i := range strings.Split(line, ",") {
			num, err := strconv.Atoi(i)
			if err != nil {
				log.Fatal("Invalid ticket format")
			}
			ownTicket = append(ownTicket, num)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	otherTickets := [][]int{}

	part1 := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "nearby tickets:" {
			continue
		}
		ticket := []int{}
		for _, i := range strings.Split(line, ",") {
			num, err := strconv.Atoi(i)
			if err != nil {
				log.Fatal("Invalid ticket format")
			}
			ticket = append(ticket, num)
		}
		if i, ok := isInvalid(ticket, rules); ok {
			part1 += i
		} else {
			otherTickets = append(otherTickets, ticket)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", part1)

	possible := make([]map[string]struct{}, 0, len(ownTicket))
	for i := 0; i < len(ownTicket); i++ {
		possible = append(possible, copyMap(ruleNames))
	}

	determined := map[string]int{}

	for {
		changed := false
		for n, i := range possible {
			if len(i) < 2 {
				continue
			}
			for _, j := range otherTickets {
				rm := []string{}
				for k := range i {
					if notInRange(j[n], rules[k]) {
						rm = append(rm, k)
					}
				}
				for _, r := range rm {
					delete(i, r)
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}
	for {
		changed := false
		for n, i := range possible {
			if len(i) == 0 {
				log.Fatal("Invalid constraints")
			}
			if len(i) == 1 {
				for k := range i {
					determined[k] = n
				}
			}
		}
		for k, v := range determined {
			for n, i := range possible {
				if n == v {
					continue
				}
				if _, ok := i[k]; ok {
					delete(i, k)
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}
	if len(determined) < len(ownTicket) {
		log.Fatal("Invalid constraints")
	}

	part2 := 1
	for k, v := range determined {
		if startsWith(k, "departure") {
			part2 *= ownTicket[v]
		}
	}
	fmt.Println("Part 2:", part2)
}

func isInvalid(ticket []int, rules map[string]Rule) (int, bool) {
outer:
	for _, i := range ticket {
		for _, v := range rules {
			if inRange(i, v) {
				continue outer
			}
		}
		return i, true
	}
	return 0, false
}

func inRange(i int, r Rule) bool {
	return i >= r.n1 && i <= r.n4 && (i <= r.n2 || i >= r.n3)
}

func notInRange(i int, r Rule) bool {
	return i < r.n1 || i > r.n4 || (i > r.n2 && i < r.n3)
}

func copyMap(m map[string]struct{}) map[string]struct{} {
	other := make(map[string]struct{}, len(m))
	for k, v := range m {
		other[k] = v
	}
	return other
}

func startsWith(s, p string) bool {
	return len(s) >= len(p) && s[:len(p)] == p
}
