package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	puzzleInput = "389547612"
)

type (
	Ring struct {
		head int
		data []int
		min  int
		max  int
	}
)

func NewRing(data []int, head, min, max int) *Ring {
	return &Ring{
		head: head,
		data: data,
		min:  min,
		max:  max,
	}
}

func (r *Ring) Step() {
	a := r.data[r.head]
	b := r.data[a]
	c := r.data[b]
	nextSet := map[int]struct{}{
		a: {},
		b: {},
		c: {},
	}
	ptr := r.data[c]
	dest := r.head - 1
	if dest < r.min {
		dest = r.max
	}
	for {
		if _, ok := nextSet[dest]; !ok {
			break
		}
		dest--
		if dest < r.min {
			dest = r.max
		}
	}
	r.data[r.head] = ptr
	succ := r.data[dest]
	r.data[dest] = a
	r.data[c] = succ
	r.head = ptr
}

func (r *Ring) Order() string {
	s := &strings.Builder{}
	for i := r.data[1]; i != 1; i = r.data[i] {
		s.WriteString(strconv.Itoa(i))
	}
	return s.String()
}

func (r *Ring) String() string {
	s := &strings.Builder{}
	s.WriteString(strconv.Itoa(r.head))
	for i := r.data[r.head]; i != r.head; i = r.data[i] {
		s.WriteByte(' ')
		s.WriteString(strconv.Itoa(i))
	}
	return s.String()
}

func (r *Ring) DoubleSuccessor() int {
	a := r.data[1]
	b := r.data[a]
	return a * b
}

func main() {
	nums := []int{}
	for _, i := range puzzleInput {
		num, err := strconv.Atoi(string(i))
		if err != nil {
			log.Fatal(err)
		}
		nums = append(nums, num)
	}
	{
		buf := make([]int, len(nums)+1)
		head := nums[0]
		min := head
		max := head
		for n, i := range nums {
			prev := nums[(n-1+len(nums))%len(nums)]
			buf[prev] = i
			if i < min {
				min = i
			}
			if i > max {
				max = i
			}
		}
		ring := NewRing(buf, head, min, max)
		for i := 0; i < 100; i++ {
			ring.Step()
		}
		fmt.Println("Part 1:", ring.Order())
	}
	{
		buf := make([]int, 1_000_001)
		head := nums[0]
		min := head
		max := head
		for n, i := range nums {
			prev := nums[(n-1+len(nums))%len(nums)]
			buf[prev] = i
			if i < min {
				min = i
			}
			if i > max {
				max = i
			}
		}
		for i := max + 1; i < 1_000_000; i++ {
			buf[i] = i + 1
		}
		first := nums[0]
		last := nums[len(nums)-1]
		buf[last] = max + 1
		buf[1_000_000] = first
		ring := NewRing(buf, head, min, 1_000_000)
		for i := 0; i < 10_000_000; i++ {
			ring.Step()
		}
		fmt.Println("Part 2:", ring.DoubleSuccessor())
	}
}
