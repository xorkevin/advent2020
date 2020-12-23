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

type (
	Deck struct {
		first    int
		last     int
		size     int
		capacity int
		data     []int
	}
)

func NewDeck(capacity int) *Deck {
	return &Deck{
		first:    0,
		last:     0,
		size:     0,
		capacity: capacity + 1,
		data:     make([]int, capacity+1),
	}
}

func (d *Deck) Len() int {
	return d.size
}

func (d *Deck) Push(v int) bool {
	n := (d.last + 1) % d.capacity
	if n == d.first {
		return false
	}
	d.data[d.last] = v
	d.last = n
	d.size++
	return true
}

func (d *Deck) Pop() (int, bool) {
	if d.first == d.last {
		return 0, false
	}
	k := d.data[d.first]
	d.first = (d.first + 1) % d.capacity
	d.size--
	return k, true
}

func (d *Deck) Copy(n int) *Deck {
	next := NewDeck(d.capacity)
	for i := d.first; next.Len() < n; i = (i + 1) % d.capacity {
		next.Push(d.data[i])
	}
	return next
}

func (d *Deck) String() string {
	s := &strings.Builder{}
	for i := d.first; i != d.last; i = (i + 1) % d.capacity {
		s.WriteString(strconv.Itoa(d.data[i]))
		s.WriteString(",")
	}
	return s.String()
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

	cards1 := []int{}
	cards2 := []int{}

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
			cards2 = append(cards2, num)
		} else {
			cards1 = append(cards1, num)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	{
		deck1 := NewDeck(len(cards1) + len(cards2))
		deck2 := NewDeck(len(cards1) + len(cards2))
		for _, i := range cards1 {
			deck1.Push(i)
		}
		for _, i := range cards2 {
			deck2.Push(i)
		}
		for deck1.Len() > 0 && deck2.Len() > 0 {
			v1, _ := deck1.Pop()
			v2, _ := deck2.Pop()
			if v1 > v2 {
				deck1.Push(v1)
				deck1.Push(v2)
			} else {
				deck2.Push(v2)
				deck2.Push(v1)
			}
		}
		sum1 := 0
		if deck1.Len() == 0 {
			for n := deck2.Len(); n > 0; n-- {
				v, _ := deck2.Pop()
				sum1 += v * n
			}
		} else {
			for n := deck1.Len(); n > 0; n-- {
				v, _ := deck1.Pop()
				sum1 += v * n
			}
		}
		fmt.Println("Part 1:", sum1)
	}
	{
		deck1 := NewDeck(len(cards1) + len(cards2))
		deck2 := NewDeck(len(cards1) + len(cards2))
		for _, i := range cards1 {
			deck1.Push(i)
		}
		for _, i := range cards2 {
			deck2.Push(i)
		}

		sum2 := 0
		if playGame(deck1, deck2) {
			for n := deck2.Len(); n > 0; n-- {
				v, _ := deck2.Pop()
				sum2 += v * n
			}
		} else {
			for n := deck1.Len(); n > 0; n-- {
				v, _ := deck1.Pop()
				sum2 += v * n
			}
		}
		fmt.Println("Part 2:", sum2)
	}
}

func playGame(deck1 *Deck, deck2 *Deck) bool {
	states := map[string]struct{}{}
	for deck1.Len() > 0 && deck2.Len() > 0 {
		current := deck1.String() + ":" + deck2.String()
		if _, ok := states[current]; ok {
			return false
		}
		states[current] = struct{}{}
		v1, _ := deck1.Pop()
		v2, _ := deck2.Pop()
		if deck1.Len() >= v1 && deck2.Len() >= v2 {
			subdeck1 := deck1.Copy(v1)
			subdeck2 := deck2.Copy(v2)
			if playGame(subdeck1, subdeck2) {
				deck2.Push(v2)
				deck2.Push(v1)
			} else {
				deck1.Push(v1)
				deck1.Push(v2)
			}
		} else {
			if v1 > v2 {
				deck1.Push(v1)
				deck1.Push(v2)
			} else {
				deck2.Push(v2)
				deck2.Push(v1)
			}
		}
	}
	return deck1.Len() == 0
}
