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
		text string
		num1 int
		num2 int
		num3 int
		num4 int
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
		num1, err := strconv.Atoi(m[2])
		if err != nil {
			log.Fatal(err)
		}
		num2, err := strconv.Atoi(m[3])
		if err != nil {
			log.Fatal(err)
		}
		num3, err := strconv.Atoi(m[4])
		if err != nil {
			log.Fatal(err)
		}
		num4, err := strconv.Atoi(m[5])
		if err != nil {
			log.Fatal(err)
		}
		rules[text] = Rule{
			text: text,
			num1: num1,
			num2: num2,
			num3: num3,
			num4: num4,
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

	countInvalid := 0

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
			countInvalid += i
		} else {
			otherTickets = append(otherTickets, ticket)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", countInvalid)

	possible := make([]map[string]struct{}, 0, len(ownTicket))
	for i := 0; i < len(ownTicket); i++ {
		possible = append(possible, copyMap(ruleNames))
	}

	determined := map[string]int{}

	for len(determined) < len(ownTicket) {
		for {
			changed := false
			for n, i := range possible {
				if len(i) < 0 {
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
				delete(i, k)
			}
		}
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

func notInRange(i int, rule Rule) bool {
	return i < rule.num1 || i > rule.num4 || (i > rule.num2 && i < rule.num3)
}

func inRange(i int, rule Rule) bool {
	return i >= rule.num1 && i <= rule.num4 && (i <= rule.num2 || i >= rule.num3)
}

func copyMap(m map[string]struct{}) map[string]struct{} {
	other := make(map[string]struct{}, len(m))
	for k, v := range m {
		other[k] = v
	}
	return other
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
