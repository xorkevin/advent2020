package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"xorkevin.dev/gnom"
)

const (
	puzzleInput = "input.txt"
)

const (
	tokenKindDefault = iota
	tokenKindEOF
	tokenKindWspace
	tokenKindNum
	tokenKindLparen
	tokenKindRparen
	tokenKindPlus
	tokenKindStar
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
	dfaWspace := gnom.NewDfa(tokenKindWspace)
	dfa.AddDfa([]rune(" "), dfaWspace)
	dfaWspace.AddDfa([]rune(" "), dfaWspace)
	dfaNum := gnom.NewDfa(tokenKindNum)
	dfa.AddDfa([]rune("0123456789"), dfaNum)
	dfaNum.AddDfa([]rune("0123456789"), dfaNum)
	dfa.AddPath([]rune("("), tokenKindLparen, tokenKindDefault)
	dfa.AddPath([]rune(")"), tokenKindRparen, tokenKindDefault)
	dfa.AddPath([]rune("+"), tokenKindPlus, tokenKindDefault)
	dfa.AddPath([]rune("*"), tokenKindStar, tokenKindDefault)
	lexer := gnom.NewDfaLexer(dfa, tokenKindDefault, tokenKindEOF, map[int]struct{}{
		tokenKindWspace: {},
	})

	sum := 0
	sum2 := 0
	scanner := bufio.NewScanner(file)
	for n := 1; scanner.Scan(); n++ {
		tokens, err := lexer.Tokenize([]rune(scanner.Text()))
		if err != nil {
			log.Fatal(err)
		}
		tokens = tokens[:len(tokens)-1]
		t1, err := parseExpr(tokens)
		if err != nil {
			log.Fatal("parse 1 ", n, err)
		}
		v1, err := NewRPNStack().Eval(t1)
		if err != nil {
			log.Fatal("eval 1 ", n, err)
		}
		sum += v1
		t2, err := parseExpr2(tokens)
		if err != nil {
			log.Fatal("parse 2 ", n, err)
		}
		v2, err := NewRPNStack().Eval(t2)
		if err != nil {
			log.Fatal("eval 2 ", n, err)
		}
		sum2 += v2
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", sum)
	fmt.Println("Part 2:", sum2)
}

func parseExpr(tokens []gnom.Token) ([]gnom.Token, error) {
	return NewShuntingYard(map[int]int{
		tokenKindStar: 1,
		tokenKindPlus: 1,
	}, map[int]int{
		tokenKindLparen: tokenKindRparen,
	}).RPN(tokens)
}

func parseExpr2(tokens []gnom.Token) ([]gnom.Token, error) {
	return NewShuntingYard(map[int]int{
		tokenKindStar: 2,
		tokenKindPlus: 1,
	}, map[int]int{
		tokenKindLparen: tokenKindRparen,
	}).RPN(tokens)
}

type (
	RPNStack struct {
		stack []int
	}
)

func NewRPNStack() *RPNStack {
	return &RPNStack{
		stack: []int{},
	}
}

func (s *RPNStack) Push(k int) {
	s.stack = append(s.stack, k)
}

func (s *RPNStack) Pop() (int, error) {
	if len(s.stack) == 0 {
		return 0, fmt.Errorf("No elements on the stack")
	}
	k := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return k, nil
}

func (s *RPNStack) Eval(tokens []gnom.Token) (int, error) {
	for _, i := range tokens {
		if i.Kind() == tokenKindNum {
			num, err := strconv.Atoi(i.Val())
			if err != nil {
				return 0, err
			}
			s.Push(num)
			continue
		}
		v2, err := s.Pop()
		if err != nil {
			return 0, err
		}
		v1, err := s.Pop()
		if err != nil {
			return 0, err
		}
		switch i.Kind() {
		case tokenKindPlus:
			s.Push(v1 + v2)
		case tokenKindStar:
			s.Push(v1 * v2)
		default:
			return 0, fmt.Errorf("Invalid op")
		}
	}
	k, err := s.Pop()
	if err != nil {
		return 0, err
	}
	if len(s.stack) != 0 {
		return 0, fmt.Errorf("Elements left on the stack")
	}
	return k, nil
}

type (
	ShuntingYard struct {
		precedence map[int]int
		pair       map[int]int
		pairOpen   map[int]struct{}
		out        []gnom.Token
		ops        []gnom.Token
		max        int
	}
)

func NewShuntingYard(precedence map[int]int, pair map[int]int) *ShuntingYard {
	max := 0
	for _, v := range precedence {
		if v > max {
			max = v + 1
		}
	}
	rpair := map[int]int{}
	pairOpen := map[int]struct{}{}
	for k, v := range pair {
		rpair[v] = k
		pairOpen[k] = struct{}{}
	}
	return &ShuntingYard{
		precedence: precedence,
		pair:       rpair,
		pairOpen:   pairOpen,
		out:        []gnom.Token{},
		ops:        []gnom.Token{},
		max:        max,
	}
}

func (s *ShuntingYard) Out() []gnom.Token {
	return s.out
}

func (s *ShuntingYard) Peek() (gnom.Token, bool) {
	if len(s.ops) == 0 {
		return gnom.Token{}, false
	}
	return s.ops[len(s.ops)-1], true
}

func (s *ShuntingYard) Pop() {
	if len(s.ops) == 0 {
		return
	}
	k := s.ops[len(s.ops)-1]
	s.ops = s.ops[:len(s.ops)-1]
	s.out = append(s.out, k)
}

func (s *ShuntingYard) PopUntil(p int) {
	for {
		top, ok := s.Peek()
		if !ok {
			return
		}
		if v, ok := s.precedence[top.Kind()]; !ok || v > p {
			return
		}
		s.Pop()
	}
}

func (s *ShuntingYard) Find(k int) error {
	for {
		top, ok := s.Peek()
		if !ok {
			return fmt.Errorf("No matching pair")
		}
		if top.Kind() == k {
			return nil
		}
		s.Pop()
	}
}

func (s *ShuntingYard) Push(token gnom.Token) error {
	if _, ok := s.pairOpen[token.Kind()]; ok {
		s.ops = append(s.ops, token)
		return nil
	}
	if t, ok := s.pair[token.Kind()]; ok {
		if err := s.Find(t); err != nil {
			return err
		}
		s.ops = s.ops[:len(s.ops)-1]
		return nil
	}
	if p, ok := s.precedence[token.Kind()]; ok {
		s.PopUntil(p)
		s.ops = append(s.ops, token)
		return nil
	}
	s.out = append(s.out, token)
	return nil
}

func (s *ShuntingYard) RPN(tokens []gnom.Token) ([]gnom.Token, error) {
	for _, i := range tokens {
		if err := s.Push(i); err != nil {
			return nil, err
		}
	}
	s.PopUntil(s.max)
	return s.Out(), nil
}
