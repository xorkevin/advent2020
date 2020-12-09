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

const (
	CodeNoop = iota
	CodeAcc
	CodeJmp
)

type (
	Instr struct {
		code int
		arg  int
	}

	Machine struct {
		instrs []Instr
		ip     int
		acc    int
	}
)

func NewMachine(instrs []Instr) *Machine {
	return &Machine{
		instrs: instrs,
		ip:     0,
		acc:    0,
	}
}

func (m *Machine) Step() (bool, error) {
	if m.ip == len(m.instrs) {
		return true, nil
	}
	if m.ip < 0 || m.ip > len(m.instrs) {
		return false, fmt.Errorf("Ip out of bounds: %d", m.ip)
	}
	instr := m.instrs[m.ip]
	switch instr.code {
	case CodeAcc:
		m.acc += instr.arg
		m.ip++
	case CodeJmp:
		m.ip += instr.arg
	case CodeNoop:
		m.ip++
	}
	return false, nil
}

func (m *Machine) Run(loopLimit int) (int, error) {
	runSet := map[int]int{}
	for {
		if _, ok := runSet[m.ip]; !ok {
			runSet[m.ip] = 0
		}
		runSet[m.ip]++
		if runSet[m.ip] > loopLimit {
			return m.acc, fmt.Errorf("Looped over %d at %d", loopLimit, m.ip)
		}
		term, err := m.Step()
		if err != nil {
			return m.acc, err
		}
		if term {
			return m.acc, nil
		}
	}
}

func parseInstr(line string) (*Instr, error) {
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return nil, fmt.Errorf("Syntax err")
	}
	var code int
	switch fields[0] {
	case "nop":
		code = CodeNoop
	case "acc":
		code = CodeAcc
	case "jmp":
		code = CodeJmp
	default:
		return nil, fmt.Errorf("Invalid code: %s", fields[0])
	}
	arg, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil, err
	}
	return &Instr{
		code: code,
		arg:  arg,
	}, nil
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

	instrs := []Instr{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		instr, err := parseInstr(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		instrs = append(instrs, *instr)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	m := NewMachine(instrs)
	part1, _ := m.Run(1)
	fmt.Println("Part 1:", part1)

	for n, i := range instrs {
		if i.code == CodeAcc {
			continue
		}
		swapToJmp := i.code == CodeNoop
		if swapToJmp {
			instrs[n] = Instr{
				code: CodeJmp,
				arg:  i.arg,
			}
		} else {
			instrs[n] = Instr{
				code: CodeNoop,
				arg:  i.arg,
			}
		}
		m := NewMachine(instrs)
		part2, err := m.Run(1)
		if err == nil {
			fmt.Println("Part 2:", part2)
			return
		}
		if swapToJmp {
			instrs[n] = Instr{
				code: CodeNoop,
				arg:  i.arg,
			}
		} else {
			instrs[n] = Instr{
				code: CodeJmp,
				arg:  i.arg,
			}
		}
	}
}
