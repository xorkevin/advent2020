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

type (
	Card struct {
		val  int
		next *Card
	}

	Deck struct {
		first *Card
		last  *Card
		size  int
	}
)

func NewDeck() *Deck {
	return &Deck{
		first: nil,
		last:  nil,
		size:  0,
	}
}

func (d *Deck) Len() int {
	return d.size
}

func (d *Deck) Push(v int) {
	if d.last == nil {
		k := &Card{
			val:  v,
			next: nil,
		}
		d.first = k
		d.last = k
	} else {
		d.last.next = &Card{
			val:  v,
			next: nil,
		}
		d.last = d.last.next
	}
	d.size++
}

func (d *Deck) First() int {
	if d.first == nil {
		return 0
	}
	return d.first.val
}

func (d *Deck) Pop() int {
	if d.first == d.last {
		k := d.first
		d.first = nil
		d.last = nil
		if k != nil {
			d.size--
			return k.val
		}
		return 0
	}
	k := d.first
	d.first = k.next
	d.size--
	return k.val
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

	deck11 := NewDeck()
	deck12 := NewDeck()
	deck21 := NewDeck()
	deck22 := NewDeck()

	var mode bool
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if line == "Player 1:" {
			mode = false
			continue
		}
		if line == "Player 2:" {
			mode = true
			continue
		}
		num, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		}
		if mode {
			deck12.Push(num)
			deck22.Push(num)
		} else {
			deck11.Push(num)
			deck21.Push(num)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for deck11.Len() > 0 && deck12.Len() > 0 {
		v1 := deck11.Pop()
		v2 := deck12.Pop()
		if v1 > v2 {
			deck11.Push(v1)
			deck11.Push(v2)
		} else {
			deck12.Push(v2)
			deck12.Push(v1)
		}
	}

	sum1 := 0
	if deck11.Len() == 0 {
		for n := deck12.Len(); n > 0; n-- {
			sum1 += deck12.Pop() * n
		}
	} else {
		for n := deck11.Len(); n > 0; n-- {
			sum1 += deck11.Pop() * n
		}
	}
	fmt.Println("Part 1:", sum1)
}
