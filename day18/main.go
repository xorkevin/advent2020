package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	puzzleInput = "input.txt"
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

	sum := 0
	sum2 := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens, err := tokenize([]byte(scanner.Text()))
		if err != nil {
			log.Fatal(err)
		}
		tree, t, err := parseExpr(tokens)
		if err != nil {
			log.Fatal(err)
		}
		if len(t) != 0 {
			log.Fatal("Remaining tokens")
		}
		sum += tree.Eval()
		tree2, t, err := parseExpr2(tokens)
		if err != nil {
			log.Fatal(err)
		}
		if len(t) != 0 {
			log.Fatal("Remaining tokens")
		}
		sum2 += tree2.Eval()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", sum)
	fmt.Println("Part 2:", sum2)
}

const (
	OpAdd = iota
	OpMul
)

const (
	TokenAdd = iota
	TokenMul
	TokenOpenParen
	TokenCloseParen
	TokenNum
)

type (
	Expr interface {
		Eval() int
	}

	Val struct {
		v int
	}

	Node struct {
		op   int
		arg1 Expr
		arg2 Expr
	}

	Token1 struct {
		kind int
		num  int
	}
)

func (v *Val) Eval() int {
	return v.v
}

func (n *Node) Eval() int {
	v1 := n.arg1.Eval()
	v2 := n.arg2.Eval()
	switch n.op {
	case OpAdd:
		return v1 + v2
	case OpMul:
		return v1 * v2
	default:
		log.Fatal("Invalid tree")
		return 0
	}
}

func TokenKindToString(kind int) string {
	switch kind {
	case TokenAdd:
		return "+"
	case TokenMul:
		return "*"
	case TokenOpenParen:
		return "("
	case TokenCloseParen:
		return ")"
	default:
		return ""
	}
}

func (t Token1) String() string {
	if t.kind == TokenNum {
		return strconv.Itoa(t.num)
	}
	return TokenKindToString(t.kind)
}

const (
	LexStateRoot = iota
	LexStateNum
)

func tokenize(line []byte) ([]Token1, error) {
	tokens := []Token1{}
	state := LexStateRoot
	buf := []byte{}
	for len(line) > 0 {
		i := line[0]
		switch state {
		case LexStateRoot:
			if i >= '0' && i <= '9' {
				state = LexStateNum
			} else {
				switch i {
				case '+':
					tokens = append(tokens, Token1{
						kind: TokenAdd,
					})
				case '*':
					tokens = append(tokens, Token1{
						kind: TokenMul,
					})
				case '(':
					tokens = append(tokens, Token1{
						kind: TokenOpenParen,
					})
				case ')':
					tokens = append(tokens, Token1{
						kind: TokenCloseParen,
					})
				case ' ':
				default:
					return nil, fmt.Errorf("Invalid character")
				}
				line = line[1:]
			}
		case LexStateNum:
			if i >= '0' && i <= '9' {
				buf = append(buf, i)
				line = line[1:]
			} else {
				num, err := strconv.Atoi(string(buf))
				if err != nil {
					return nil, err
				}
				tokens = append(tokens, Token1{
					kind: TokenNum,
					num:  num,
				})
				buf = []byte{}
				state = LexStateRoot
			}
		}
	}
	if len(buf) > 0 {
		num, err := strconv.Atoi(string(buf))
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, Token1{
			kind: TokenNum,
			num:  num,
		})
		buf = []byte{}
		state = LexStateRoot
	}
	return tokens, nil
}

func parseParen(tokens []Token1) (Expr, []Token1, error) {
	val, tokens, err := parseExpr(tokens)
	if err != nil {
		return nil, nil, err
	}
	if len(tokens) == 0 {
		return nil, nil, fmt.Errorf("EOF")
	}
	i := tokens[0]
	tokens = tokens[1:]
	if i.kind != TokenCloseParen {
		return nil, nil, fmt.Errorf("No close paren")
	}
	return val, tokens, nil
}

func parseVal(tokens []Token1) (Expr, []Token1, error) {
	if len(tokens) == 0 {
		return nil, nil, fmt.Errorf("EOF")
	}
	i := tokens[0]
	tokens = tokens[1:]
	switch i.kind {
	case TokenNum:
		return &Val{
			v: i.num,
		}, tokens, nil
	case TokenOpenParen:
		return parseParen(tokens)
	default:
		return nil, nil, fmt.Errorf("Invalid token")
	}
}

func parseExpr(tokens []Token1) (Expr, []Token1, error) {
	val, tokens, err := parseVal(tokens)
	if err != nil {
		return nil, nil, err
	}
	for len(tokens) != 0 {
		i := tokens[0]
		var op int
		switch i.kind {
		case TokenAdd:
			op = OpAdd
		case TokenMul:
			op = OpMul
		case TokenCloseParen:
			return val, tokens, nil
		default:
			return nil, nil, fmt.Errorf("Invalid token")
		}
		tokens = tokens[1:]
		var arg2 Expr
		arg2, tokens, err = parseVal(tokens)
		if err != nil {
			return nil, nil, err
		}
		val = &Node{
			op:   op,
			arg1: val,
			arg2: arg2,
		}
	}
	return val, tokens, nil
}

func parseParen2(tokens []Token1) (Expr, []Token1, error) {
	val, tokens, err := parseExpr2(tokens)
	if err != nil {
		return nil, nil, err
	}
	if len(tokens) == 0 {
		return nil, nil, fmt.Errorf("EOF")
	}
	i := tokens[0]
	tokens = tokens[1:]
	if i.kind != TokenCloseParen {
		return nil, nil, fmt.Errorf("No close paren")
	}
	return val, tokens, nil
}

func parseVal2(tokens []Token1) (Expr, []Token1, error) {
	if len(tokens) == 0 {
		return nil, nil, fmt.Errorf("EOF")
	}
	i := tokens[0]
	tokens = tokens[1:]
	switch i.kind {
	case TokenNum:
		return &Val{
			v: i.num,
		}, tokens, nil
	case TokenOpenParen:
		return parseParen2(tokens)
	default:
		return nil, nil, fmt.Errorf("Invalid token")
	}
}

func parseAddExpr2(tokens []Token1) (Expr, []Token1, error) {
	val, tokens, err := parseVal2(tokens)
	if err != nil {
		return nil, nil, err
	}
	for len(tokens) != 0 {
		i := tokens[0]
		var op int
		switch i.kind {
		case TokenAdd:
			op = OpAdd
		case TokenMul, TokenCloseParen:
			return val, tokens, nil
		default:
			return nil, nil, fmt.Errorf("Invalid token")
		}
		tokens = tokens[1:]
		var arg2 Expr
		arg2, tokens, err = parseVal2(tokens)
		if err != nil {
			return nil, nil, err
		}
		val = &Node{
			op:   op,
			arg1: val,
			arg2: arg2,
		}
	}
	return val, tokens, nil
}

func parseExpr2(tokens []Token1) (Expr, []Token1, error) {
	val, tokens, err := parseAddExpr2(tokens)
	if err != nil {
		return nil, nil, err
	}
	for len(tokens) != 0 {
		i := tokens[0]
		var op int
		switch i.kind {
		case TokenMul:
			op = OpMul
		case TokenCloseParen:
			return val, tokens, nil
		default:
			return nil, nil, fmt.Errorf("Invalid token")
		}
		tokens = tokens[1:]
		var arg2 Expr
		arg2, tokens, err = parseAddExpr2(tokens)
		if err != nil {
			return nil, nil, err
		}
		val = &Node{
			op:   op,
			arg1: val,
			arg2: arg2,
		}
	}
	return val, tokens, nil
}
