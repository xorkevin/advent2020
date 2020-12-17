package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	puzzleInput = "input.txt"
)

type (
	Point struct {
		x, y, z, w int
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

	active := map[Point]struct{}{}
	active2 := map[Point]struct{}{}

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		for j, c := range []byte(scanner.Text()) {
			if c == '#' {
				active[Point{j, i, 0, 0}] = struct{}{}
				active2[Point{j, i, 0, 0}] = struct{}{}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 6; i++ {
		active = step(active, false)
	}
	fmt.Println("Part 1:", len(active))
	for i := 0; i < 6; i++ {
		active2 = step(active2, true)
	}
	fmt.Println("Part 2:", len(active2))
}

func step(active map[Point]struct{}, dim4 bool) map[Point]struct{} {
	next := map[Point]struct{}{}
	nbc := map[Point]int{}
	for p := range active {
		var nb []Point
		if dim4 {
			nb = neighborsW(p)
		} else {
			nb = neighbors(p)
		}
		c := countNeighbors(nb, active)
		if c == 2 || c == 3 {
			next[p] = struct{}{}
		}
		for _, n := range nb {
			if _, ok := nbc[n]; !ok {
				nbc[n] = 0
			}
			nbc[n]++
		}
	}
	for p, c := range nbc {
		if c == 3 {
			next[p] = struct{}{}
		}
	}
	return next
}

func countNeighbors(points []Point, active map[Point]struct{}) int {
	count := 0
	for _, p := range points {
		if _, ok := active[p]; ok {
			count++
		}
	}
	return count
}

func neighbors(p Point) []Point {
	n := make([]Point, 0, 26)
	for i := p.x - 1; i <= p.x+1; i++ {
		for j := p.y - 1; j <= p.y+1; j++ {
			for k := p.z - 1; k <= p.z+1; k++ {
				x := Point{i, j, k, 0}
				if x != p {
					n = append(n, x)
				}
			}
		}
	}
	return n
}

func neighborsW(p Point) []Point {
	n := make([]Point, 0, 26)
	for i := p.x - 1; i <= p.x+1; i++ {
		for j := p.y - 1; j <= p.y+1; j++ {
			for k := p.z - 1; k <= p.z+1; k++ {
				for l := p.w - 1; l <= p.w+1; l++ {
					x := Point{i, j, k, l}
					if x == p {
						continue
					}
					n = append(n, x)
				}
			}
		}
	}
	return n
}
