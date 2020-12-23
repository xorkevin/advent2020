package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	StringSet struct {
		elements map[string]struct{}
	}
)

func NewStringSet(items ...string) *StringSet {
	elements := map[string]struct{}{}
	for _, i := range items {
		elements[i] = struct{}{}
	}
	return &StringSet{
		elements: elements,
	}
}

func (s *StringSet) Len() int {
	return len(s.elements)
}

func (s *StringSet) First() string {
	for i := range s.elements {
		return i
	}
	return ""
}

func (s *StringSet) Has(k string) bool {
	_, ok := s.elements[k]
	return ok
}

func (s *StringSet) Add(k string) {
	s.elements[k] = struct{}{}
}

func (s *StringSet) Rm(k string) {
	delete(s.elements, k)
}

func (s *StringSet) Intersect(other *StringSet) {
	for i := range other.elements {
		if !s.Has(i) {
			s.Rm(i)
		}
	}
	toRm := map[string]struct{}{}
	for i := range s.elements {
		if !other.Has(i) {
			toRm[i] = struct{}{}
		}
	}
	for i := range toRm {
		s.Rm(i)
	}
}

func (s *StringSet) Union(other *StringSet) {
	for i := range other.elements {
		s.Add(i)
	}
}

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

	// allergens to ingredients
	possible := map[string]*StringSet{}

	ingredients := map[string]int{}
	allergens := NewStringSet()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " (contains ")
		if len(line) != 2 {
			log.Fatal("Invalid line")
		}
		l0 := strings.Split(line[0], " ")
		l1 := strings.Split(line[1][:len(line[1])-1], ", ")

		ing := NewStringSet()
		for _, i := range l0 {
			if _, ok := ingredients[i]; !ok {
				ingredients[i] = 0
			}
			ingredients[i]++
			ing.Add(i)
		}

		for _, i := range l1 {
			allergens.Add(i)
			if _, ok := possible[i]; !ok {
				possible[i] = NewStringSet(l0...)
			} else {
				possible[i].Intersect(ing)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	possible1 := NewStringSet()
	for _, v := range possible {
		possible1.Union(v)
	}
	count1 := 0
	for k, v := range ingredients {
		if !possible1.Has(k) {
			count1 += v
		}
	}
	fmt.Println("Part 1:", count1)

	algToIng := map[string]string{}
	ingToAlg := map[string]string{}
	ingList := []string{}

	for len(ingToAlg) < len(allergens.elements) {
		for k, v := range possible {
			if _, ok := algToIng[k]; ok {
				continue
			}
			if v.Len() == 0 {
				log.Fatal("has zero options")
			}
			if v.Len() != 1 {
				continue
			}
			x := v.First()
			algToIng[k] = x
			ingToAlg[x] = k
			ingList = append(ingList, x)
			for k2, v2 := range possible {
				if k2 == k {
					continue
				}
				v2.Rm(x)
			}
		}
	}

	sort.Slice(ingList, func(i, j int) bool {
		return ingToAlg[ingList[i]] < ingToAlg[ingList[j]]
	})
	fmt.Println("Part 2:", strings.Join(ingList, ","))
}
