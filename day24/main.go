package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"xorkevin.dev/gnom"
)

const (
	puzzleInput = "input.txt"
)

type (
	HexCoord struct {
		x, z int
	}
)

func (c *HexCoord) MoveE() {
	c.x++
}

func (c *HexCoord) MoveW() {
	c.x--
}

func (c *HexCoord) MoveNE() {
	c.x++
	c.z--
}

func (c *HexCoord) MoveSW() {
	c.x--
	c.z++
}

func (c *HexCoord) MoveNW() {
	c.z--
}

func (c *HexCoord) MoveSE() {
	c.z++
}

func (c *HexCoord) Neighbors() []HexCoord {
	return []HexCoord{
		{x: c.x + 1, z: c.z},
		{x: c.x - 1, z: c.z},
		{x: c.x + 1, z: c.z - 1},
		{x: c.x - 1, z: c.z + 1},
		{x: c.x, z: c.z - 1},
		{x: c.x, z: c.z + 1},
	}
}

type (
	TileSet struct {
		blackTiles map[HexCoord]struct{}
	}
)

func NewTileSet() *TileSet {
	return &TileSet{
		blackTiles: map[HexCoord]struct{}{},
	}
}

func (s *TileSet) Flip(c HexCoord) {
	if _, ok := s.blackTiles[c]; ok {
		delete(s.blackTiles, c)
	} else {
		s.blackTiles[c] = struct{}{}
	}
}

func (s *TileSet) Step() {
	next := map[HexCoord]struct{}{}
	nbc := map[HexCoord]int{}
	for c := range s.blackTiles {
		nb := c.Neighbors()
		count := 0
		for _, i := range nb {
			if _, ok := s.blackTiles[i]; ok {
				count++
			}
			if _, ok := nbc[i]; !ok {
				nbc[i] = 0
			}
			nbc[i]++
		}
		if count == 1 || count == 2 {
			next[c] = struct{}{}
		}
	}
	for k, v := range nbc {
		if v == 2 {
			next[k] = struct{}{}
		}
	}
	s.blackTiles = next
}

const (
	tokenKindDefault = iota
	tokenKindEOF
	tokenKindE
	tokenKindW
	tokenKindNE
	tokenKindSW
	tokenKindNW
	tokenKindSE
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

	dfa := gnom.NewDfa(tokenKindDefault)
	dfa.AddPath([]rune("e"), tokenKindE, tokenKindDefault)
	dfa.AddPath([]rune("w"), tokenKindW, tokenKindDefault)
	dfa.AddPath([]rune("ne"), tokenKindNE, tokenKindDefault)
	dfa.AddPath([]rune("sw"), tokenKindSW, tokenKindDefault)
	dfa.AddPath([]rune("nw"), tokenKindNW, tokenKindDefault)
	dfa.AddPath([]rune("se"), tokenKindSE, tokenKindDefault)
	lexer := gnom.NewDfaLexer(dfa, tokenKindDefault, tokenKindEOF, map[int]struct{}{})

	ts := NewTileSet()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens, err := lexer.Tokenize([]rune(scanner.Text()))
		if err != nil {
			log.Fatal(err)
		}
		pt := HexCoord{0, 0}
		for _, i := range tokens {
			switch i.Kind() {
			case tokenKindE:
				pt.MoveE()
			case tokenKindW:
				pt.MoveW()
			case tokenKindNE:
				pt.MoveNE()
			case tokenKindSW:
				pt.MoveSW()
			case tokenKindNW:
				pt.MoveNW()
			case tokenKindSE:
				pt.MoveSE()
			}
		}
		ts.Flip(pt)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", len(ts.blackTiles))
	for i := 0; i < 100; i++ {
		ts.Step()
	}
	fmt.Println("Part 2:", len(ts.blackTiles))
}
