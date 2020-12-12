package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

const (
	puzzleInput = "input.txt"
)

type (
	Action struct {
		act string
		v   int
	}

	Ship struct {
		dir string
		x   int
		y   int
		wx  int
		wy  int
	}
)

func NewShip() *Ship {
	return &Ship{
		dir: "E",
		x:   0,
		y:   0,
		wx:  10,
		wy:  -1,
	}
}

func (s *Ship) turnLeft() {
	switch s.dir {
	case "N":
		s.dir = "W"
	case "E":
		s.dir = "N"
	case "S":
		s.dir = "E"
	case "W":
		s.dir = "S"
	}
}

func (s *Ship) turnRight() {
	switch s.dir {
	case "N":
		s.dir = "E"
	case "E":
		s.dir = "S"
	case "S":
		s.dir = "W"
	case "W":
		s.dir = "N"
	}
}

func (s *Ship) Step(a Action) {
	switch a.act {
	case "N":
		s.y -= a.v
	case "S":
		s.y += a.v
	case "E":
		s.x += a.v
	case "W":
		s.x -= a.v
	case "L":
		turns := a.v / 90
		for i := 0; i < turns; i++ {
			s.turnLeft()
		}
	case "R":
		turns := a.v / 90
		for i := 0; i < turns; i++ {
			s.turnRight()
		}
	case "F":
		switch s.dir {
		case "N":
			s.y -= a.v
		case "E":
			s.x += a.v
		case "S":
			s.y += a.v
		case "W":
			s.x -= a.v
		}
	}
}

func (s *Ship) turnWaypointLeft() {
	a := s.wx
	b := s.wy
	s.wx = b
	s.wy = -a
}

func (s *Ship) turnWaypointRight() {
	a := s.wx
	b := s.wy
	s.wx = -b
	s.wy = a
}

func (s *Ship) StepWaypoint(a Action) {
	switch a.act {
	case "N":
		s.wy -= a.v
	case "S":
		s.wy += a.v
	case "E":
		s.wx += a.v
	case "W":
		s.wx -= a.v
	case "L":
		turns := a.v / 90
		for i := 0; i < turns; i++ {
			s.turnWaypointLeft()
		}
	case "R":
		turns := a.v / 90
		for i := 0; i < turns; i++ {
			s.turnWaypointRight()
		}
	case "F":
		s.x += s.wx * a.v
		s.y += s.wy * a.v
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
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

	re := regexp.MustCompile(`^([A-Z])([0-9]+)$`)

	s := NewShip()
	s2 := NewShip()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m := re.FindStringSubmatch(scanner.Text())
		if m == nil {
			log.Fatal("Invalid line")
		}
		act := m[1]
		num, err := strconv.Atoi(m[2])
		if err != nil {
			log.Fatal(err)
		}
		a := Action{
			act: act,
			v:   num,
		}
		s.Step(a)
		s2.StepWaypoint(a)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", abs(s.x)+abs(s.y))
	fmt.Println("Part 2:", abs(s2.x)+abs(s2.y))
}
