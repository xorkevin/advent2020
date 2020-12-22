package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

const (
	puzzleInput = "input.txt"
)

type (
	Edge struct {
		f, b int
	}

	Tile struct {
		id         int
		r, t, l, b Edge
	}

	Transform struct {
		flip     bool
		rotation int
	}
)

var (
	allTransforms = []Transform{
		{
			flip:     false,
			rotation: 0,
		},
		{
			flip:     false,
			rotation: 1,
		},
		{
			flip:     false,
			rotation: 2,
		},
		{
			flip:     false,
			rotation: 3,
		},
		{
			flip:     true,
			rotation: 0,
		},
		{
			flip:     true,
			rotation: 1,
		},
		{
			flip:     true,
			rotation: 2,
		},
		{
			flip:     true,
			rotation: 3,
		},
	}
)

func getEdge(tile Tile, tr Transform, side int) int {
	angle := (((side - tr.rotation) % 4) + 4) % 4
	var e Edge
	switch angle {
	case 0:
		e = tile.r
	case 1:
		e = tile.t
	case 2:
		e = tile.l
	case 3:
		e = tile.b
	}
	if tr.flip {
		return e.b
	}
	return e.f
}

type (
	TrTile struct {
		tile Tile
		tr   Transform
	}
)

func (t TrTile) Edge(side int) int {
	return getEdge(t.tile, t.tr, side)
}

type (
	TileIter struct {
		idx   int
		trIdx int
		tiles []Tile
	}
)

func NewTileIter(tiles []Tile) *TileIter {
	return &TileIter{
		idx:   0,
		trIdx: 0,
		tiles: tiles,
	}
}

func (t *TileIter) Next() (TrTile, bool) {
	if t.trIdx >= len(allTransforms) {
		t.trIdx = 0
		t.idx++
	}
	if t.idx >= len(t.tiles) {
		return TrTile{}, false
	}
	k := TrTile{
		tile: t.tiles[t.idx],
		tr:   allTransforms[t.trIdx],
	}
	t.trIdx++
	return k, true
}

type (
	ValidPeers struct {
		peers []map[int][]TrTile
	}
)

func NewValidPeers() *ValidPeers {
	peers := make([]map[int][]TrTile, 0, 4)
	for i := 0; i < 4; i++ {
		peers = append(peers, map[int][]TrTile{})
	}
	return &ValidPeers{
		peers: peers,
	}
}

func (p *ValidPeers) Get(i int, id int) []TrTile {
	if i >= len(p.peers) {
		return []TrTile{}
	}
	k, ok := p.peers[i][id]
	if !ok {
		return []TrTile{}
	}
	return k
}

func (p *ValidPeers) AddTrTile(t TrTile) {
	for i := 0; i < 4; i++ {
		e := t.Edge((i + 2) % 4)
		if _, ok := p.peers[i][e]; !ok {
			p.peers[i][e] = []TrTile{}
		}
		p.peers[i][e] = append(p.peers[i][e], t)
	}
}

type (
	TilePool struct {
		pool map[int]struct{}
	}
)

func NewTilePool(tiles []Tile) *TilePool {
	pool := make(map[int]struct{}, len(tiles))
	for _, i := range tiles {
		pool[i.id] = struct{}{}
	}
	return &TilePool{
		pool: pool,
	}
}

func (p *TilePool) Has(i int) bool {
	_, ok := p.pool[i]
	return ok
}

func (p *TilePool) Pop(i int) {
	delete(p.pool, i)
}

func (p *TilePool) Push(i int) {
	p.pool[i] = struct{}{}
}

type (
	Grid struct {
		dimx int
		dimy int
		grid [][]TrTile
	}
)

func NewGrid(dimx, dimy int) *Grid {
	grid := make([][]TrTile, 0, dimy)
	for i := 0; i < dimx; i++ {
		grid = append(grid, make([]TrTile, dimx))
	}
	return &Grid{
		dimx: dimx,
		dimy: dimy,
		grid: grid,
	}
}

func (g *Grid) SetTopLeft(t TrTile) {
	g.grid[0][0] = t
}

func (g *Grid) SetTop(idx int, pool *TilePool, peers *ValidPeers) bool {
	if idx >= g.dimx {
		return true
	}
	e := g.grid[0][idx-1].Edge(0)
	for _, i := range peers.Get(0, e) {
		if !pool.Has(i.tile.id) {
			continue
		}
		g.grid[0][idx] = i
		pool.Pop(i.tile.id)
		if g.SetTop(idx+1, pool, peers) {
			pool.Push(i.tile.id)
			return true
		}
		pool.Push(i.tile.id)
	}
	return false
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

	tiles := []Tile{}
	tileID := 0
	tile := [][]byte{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if tileID == 0 {
			num, err := strconv.Atoi(line[5 : len(line)-1])
			if err != nil {
				log.Fatal(err)
			}
			tileID = num
			continue
		}
		if line == "" {
			h := len(tile)
			w := len(tile[0])
			if h != w {
				log.Fatal("tile is not square")
			}
			rf := 0
			rb := 0
			lf := 0
			lb := 0
			for i := 0; i < h; i++ {
				rf <<= 1
				if tile[i][w-1] == '#' {
					rf += 1
				}
				rb <<= 1
				if tile[h-i-1][w-1] == '#' {
					rb += 1
				}
				lf <<= 1
				if tile[h-i-1][0] == '#' {
					lf += 1
				}
				lb <<= 1
				if tile[i][0] == '#' {
					lb += 1
				}
			}
			tf := 0
			tb := 0
			bf := 0
			bb := 0
			for i := 0; i < w; i++ {
				tf <<= 1
				if tile[0][i] == '#' {
					tf += 1
				}
				bb <<= 1
				if tile[0][w-i-1] == '#' {
					bb += 1
				}
				bf <<= 1
				if tile[h-1][w-i-1] == '#' {
					bf += 1
				}
				tb <<= 1
				if tile[h-1][i] == '#' {
					tb += 1
				}
			}
			tiles = append(tiles, Tile{
				id: tileID,
				r:  Edge{f: rf, b: rb},
				t:  Edge{f: tf, b: tb},
				l:  Edge{f: lf, b: lb},
				b:  Edge{f: bf, b: bb},
			})
			tileID = 0
			tile = [][]byte{}
			continue
		}
		tile = append(tile, []byte(line))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	squareSize := math.Sqrt(float64(len(tiles)))
	if math.Floor(squareSize) != squareSize {
		log.Fatal("Not a square")
	}
	dim := int(squareSize)

	validPeers := NewValidPeers()
	tileIter := NewTileIter(tiles)
	for t, ok := tileIter.Next(); ok; t, ok = tileIter.Next() {
		validPeers.AddTrTile(t)
	}

	_ = NewTilePool(tiles)

	_ = NewGrid(dim, dim)

	fmt.Println(dim)
}
