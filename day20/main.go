package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
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

func getEdgeReverse(tile Tile, tr Transform, side int) int {
	angle := (((side - tr.rotation) % 4) + 4) % 4
	var e Edge
	switch angle {
	case 0:
		e = tile.r
	case 1:
		e = tile.b
	case 2:
		e = tile.l
	case 3:
		e = tile.t
	}
	if tr.flip {
		return e.f
	}
	return e.b
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

func (t TrTile) EdgeReverse(side int) int {
	return getEdgeReverse(t.tile, t.tr, side)
}

type (
	TileIter struct {
		idx     int
		trIdx   int
		tiles   []Tile
		current TrTile
	}
)

func NewTileIter(tiles []Tile) *TileIter {
	return &TileIter{
		idx:   0,
		trIdx: 0,
		tiles: tiles,
	}
}

func (t *TileIter) Next() bool {
	if t.trIdx >= len(allTransforms) {
		t.trIdx = 0
		t.idx++
	}
	if t.idx >= len(t.tiles) {
		return false
	}
	t.current = TrTile{
		tile: t.tiles[t.idx],
		tr:   allTransforms[t.trIdx],
	}
	t.trIdx++
	return true
}

func (t *TileIter) Get() TrTile {
	return t.current
}

type (
	Tuple2 struct {
		a, b int
	}

	ValidPeers struct {
		right     map[int][]TrTile
		down      map[int][]TrTile
		rightDown map[Tuple2][]TrTile
	}
)

func NewValidPeers() *ValidPeers {
	return &ValidPeers{
		right:     map[int][]TrTile{},
		down:      map[int][]TrTile{},
		rightDown: map[Tuple2][]TrTile{},
	}
}

func (p *ValidPeers) GetRight(id int) []TrTile {
	k, ok := p.right[id]
	if !ok {
		return []TrTile{}
	}
	return k
}

func (p *ValidPeers) GetDown(id int) []TrTile {
	k, ok := p.down[id]
	if !ok {
		return []TrTile{}
	}
	return k
}

func (p *ValidPeers) GetRightDown(r, d int) []TrTile {
	k, ok := p.rightDown[Tuple2{a: r, b: d}]
	if !ok {
		return []TrTile{}
	}
	return k
}

func (p *ValidPeers) AddTrTile(t TrTile) {
	r := t.Edge(2)
	d := t.Edge(1)
	rd := Tuple2{a: r, b: d}
	if _, ok := p.right[r]; !ok {
		p.right[r] = []TrTile{}
	}
	p.right[r] = append(p.right[r], t)
	if _, ok := p.down[d]; !ok {
		p.down[d] = []TrTile{}
	}
	p.down[d] = append(p.down[d], t)
	if _, ok := p.rightDown[rd]; !ok {
		p.rightDown[rd] = []TrTile{}
	}
	p.rightDown[rd] = append(p.rightDown[rd], t)
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

func (g *Grid) SetTopLeft(t TrTile, pool *TilePool, peers *ValidPeers) bool {
	g.grid[0][0] = t
	return g.SetTop(1, pool, peers)
}

func (g *Grid) SetTop(idx int, pool *TilePool, peers *ValidPeers) bool {
	if idx >= g.dimx {
		return g.SetLeft(1, pool, peers)
	}
	e := g.grid[0][idx-1].EdgeReverse(0)
	for _, i := range peers.GetRight(e) {
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

func (g *Grid) SetLeft(idx int, pool *TilePool, peers *ValidPeers) bool {
	if idx >= g.dimy {
		return g.SetCenter(1, 1, pool, peers)
	}
	e := g.grid[idx-1][0].EdgeReverse(3)
	for _, i := range peers.GetDown(e) {
		if !pool.Has(i.tile.id) {
			continue
		}
		g.grid[idx][0] = i
		pool.Pop(i.tile.id)
		if g.SetLeft(idx+1, pool, peers) {
			pool.Push(i.tile.id)
			return true
		}
		pool.Push(i.tile.id)
	}
	return false
}

func (g *Grid) SetCenter(idxr int, idxc int, pool *TilePool, peers *ValidPeers) bool {
	if idxc >= g.dimx {
		idxc = 1
		idxr++
	}
	if idxr >= g.dimy {
		return true
	}
	le := g.grid[idxr][idxc-1].EdgeReverse(0)
	te := g.grid[idxr-1][idxc].EdgeReverse(3)
	for _, i := range peers.GetRightDown(le, te) {
		if !pool.Has(i.tile.id) {
			continue
		}
		g.grid[idxr][idxc] = i
		pool.Pop(i.tile.id)
		if g.SetCenter(idxr, idxc+1, pool, peers) {
			pool.Push(i.tile.id)
			return true
		}
		pool.Push(i.tile.id)
	}
	return false
}

func mapCoords(grid [][]byte, tr Transform, i, j int, dim int) byte {
	ip := dim - i - 1
	jp := dim - j - 1
	if tr.flip {
		switch tr.rotation % 4 {
		case 0:
			return grid[ip][j]
		case 1:
			return grid[jp][ip]
		case 2:
			return grid[i][jp]
		case 3:
			return grid[j][i]
		}
	} else {
		switch tr.rotation % 4 {
		case 0:
			return grid[i][j]
		case 1:
			return grid[j][ip]
		case 2:
			return grid[ip][jp]
		case 3:
			return grid[jp][i]
		}
	}
	return 0
}

func matchMonster(grid [][]byte, midx [][]int, tr Transform, i, j int, dim int) bool {
	for r, row := range midx {
		for _, c := range row {
			if mapCoords(grid, tr, i+r, j+c, dim) != '#' {
				return false
			}
		}
	}
	return true
}

const (
	monster = `                  # 
#    ##    ##    ###
 #  #  #  #  #  #   `
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

	tileMap := map[int][][]byte{}
	tiles := []Tile{}
	tileID := 0
	tile := [][]byte{}
	var tileDim int

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
			tileDim = h
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
			tileMap[tileID] = tile
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
	for iter := NewTileIter(tiles); iter.Next(); {
		validPeers.AddTrTile(iter.Get())
	}

	pool := NewTilePool(tiles)
	grid := NewGrid(dim, dim)

	foundLayout := false
	for iter := NewTileIter(tiles); iter.Next(); {
		k := iter.Get()
		pool.Pop(k.tile.id)
		if grid.SetTopLeft(k, pool, validPeers) {
			pool.Push(k.tile.id)
			foundLayout = true
			break
		}
		pool.Push(k.tile.id)
	}
	if !foundLayout {
		log.Fatal("Failed to find grid layout")
	}
	product := grid.grid[0][0].tile.id * grid.grid[0][grid.dimx-1].tile.id * grid.grid[grid.dimy-1][0].tile.id * grid.grid[grid.dimy-1][grid.dimx-1].tile.id
	fmt.Println("Part 1:", product)

	pixelBlockDim := tileDim - 2
	pixelsDim := pixelBlockDim * dim
	pixels := make([][]byte, 0, pixelsDim)
	for i := 0; i < pixelsDim; i++ {
		pixels = append(pixels, make([]byte, pixelsDim))
	}

	pixelCount := 0
	for gy, r := range grid.grid {
		for gx, c := range r {
			for i := 0; i < pixelBlockDim; i++ {
				for j := 0; j < pixelBlockDim; j++ {
					k := mapCoords(tileMap[c.tile.id], c.tr, i+1, j+1, tileDim)
					pixels[gy*pixelBlockDim+i][gx*pixelBlockDim+j] = k
					if k == '#' {
						pixelCount++
					}
				}
			}
		}
	}

	monsterPixelCount := 0
	monsterIdx := [][]int{}
	var monsterDimX int
	for _, i := range strings.Split(monster, "\n") {
		monsterDimX = len(i)
		row := []int{}
		for n, j := range i {
			if j == '#' {
				row = append(row, n)
				monsterPixelCount++
			}
		}
		monsterIdx = append(monsterIdx, row)
	}
	monsterDimY := len(monsterIdx)

	monsterCount := 0
	for _, tr := range allTransforms {
		for i := 0; i < pixelsDim-monsterDimY; i++ {
			for j := 0; j < pixelsDim-monsterDimX; j++ {
				if matchMonster(pixels, monsterIdx, tr, i, j, pixelsDim) {
					monsterCount++
				}
			}
		}
		if monsterCount > 0 {
			break
		}
	}

	if monsterCount == 0 {
		log.Fatal("No monsters found")
	}

	fmt.Println("Part 2:", pixelCount-monsterPixelCount*monsterCount)
}

func printPixels(pixels [][]byte, tr Transform, dim, tileDim int) {
	for i := 0; i < dim; i++ {
		if i%tileDim == 0 {
			fmt.Println()
		}
		for j := 0; j < dim; j++ {
			if j%tileDim == 0 {
				fmt.Print(" ")
			}
			fmt.Print(string(mapCoords(pixels, tr, i, j, dim)))
		}
		fmt.Println()
	}
}
