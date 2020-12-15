package main

import (
	"fmt"
)

var (
	puzzleInput = []int{19, 20, 14, 0, 9, 1}
)

type (
	Num struct {
		cur  int
		prev int
	}

	Hist struct {
		hist map[int]Num
		idx  int
	}
)

func NewHist() *Hist {
	return &Hist{
		hist: map[int]Num{},
		idx:  0,
	}
}

func (h *Hist) Say(num int) {
	h.idx++
	if v, ok := h.hist[num]; ok {
		h.hist[num] = Num{
			cur:  h.idx,
			prev: v.cur,
		}
	} else {
		h.hist[num] = Num{
			cur:  h.idx,
			prev: 0,
		}
	}
}

func (h *Hist) Diff(num int) int {
	if v, ok := h.hist[num]; ok && v.prev != 0 {
		return v.cur - v.prev
	}
	return 0
}

func main() {
	hist := NewHist()
	var prev int
	for _, i := range puzzleInput {
		prev = i
		hist.Say(prev)
	}
	for hist.idx < 30000000 {
		prev = hist.Diff(prev)
		hist.Say(prev)
		if hist.idx == 2020 {
			fmt.Println("Part 1:", prev)
		}
	}
	fmt.Println("Part 2:", prev)
}
