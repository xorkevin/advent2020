package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	Dfa struct {
		kind  int
		nodes map[byte]*Dfa
	}
)

func NewDfa(kind int) *Dfa {
	return &Dfa{
		kind:  kind,
		nodes: map[byte]*Dfa{},
	}
}

func (d *Dfa) AddDfa(s []byte, dfa *Dfa) {
	for _, c := range s {
		d.nodes[c] = dfa
	}
}

func (d *Dfa) AddPath(path []byte, kind int, def int) *Dfa {
	if len(path) == 0 {
		d.kind = kind
		return d
	}
	c := path[0]
	path = path[1:]
	if _, ok := d.nodes[c]; !ok {
		d.nodes[c] = NewDfa(def)
	}
	return d.nodes[c].AddPath(path, kind, def)
}

func (d *Dfa) Match(c byte) (*Dfa, bool) {
	next, ok := d.nodes[c]
	if !ok {
		return nil, false
	}
	return next, true
}

func (d *Dfa) Kind() int {
	return d.kind
}

func (d *Dfa) writeIndent(s *strings.Builder, indent int) {
	for i := 0; i < indent; i++ {
		s.WriteByte(' ')
	}
}

func (d *Dfa) ToString(s *strings.Builder, indent int) {
	for k, v := range d.nodes {
		d.writeIndent(s, indent)
		s.WriteByte(k)
		s.WriteByte('\n')
		v.ToString(s, indent+2)
	}
}

func (d Dfa) String() string {
	s := &strings.Builder{}
	d.ToString(s, 0)
	return s.String()
}

type (
	Lexer struct {
		dfa     *Dfa
		def     int
		eof     int
		ignored map[int]struct{}
	}

	TokenStack struct {
		tokens []Token
	}

	Token struct {
		kind int
		val  string
	}
)

func NewToken(kind int, val string) *Token {
	return &Token{
		kind: kind,
		val:  val,
	}
}

func (t *Token) Kind() int {
	return t.kind
}

func (t *Token) Val() string {
	return t.val
}

func NewTokenStack() *TokenStack {
	return &TokenStack{
		tokens: []Token{},
	}
}

func (s *TokenStack) Empty() bool {
	return len(s.tokens) == 0
}

func (s *TokenStack) Push(t Token) {
	s.tokens = append(s.tokens, t)
}

func (s *TokenStack) Pop() (*Token, bool) {
	if s.Empty() {
		return nil, false
	}
	k := s.tokens[len(s.tokens)-1]
	s.tokens = s.tokens[:len(s.tokens)-1]
	return &k, true
}

func NewLexer(dfa *Dfa, def, eof int, ignored map[int]struct{}) *Lexer {
	return &Lexer{
		dfa:     dfa,
		def:     def,
		eof:     eof,
		ignored: ignored,
	}
}

var (
	ErrLex = errors.New("lexer error")
)

func (l *Lexer) Next(chars []byte) (*Token, []byte, error) {
	s := &strings.Builder{}
	n := l.dfa
	for {
		if len(chars) == 0 {
			break
		}
		c := chars[0]
		next, ok := n.Match(c)
		if !ok {
			break
		}
		n = next
		s.WriteByte(c)
		chars = chars[1:]
	}
	if n.Kind() == l.def {
		if s.Len() == 0 {
			return NewToken(l.eof, ""), chars, nil
		}
		return nil, nil, fmt.Errorf("Invalid token: %s: %w", s.String(), ErrLex)
	}
	return NewToken(n.Kind(), s.String()), chars, nil
}

func (l *Lexer) Tokenize(chars []byte) ([]Token, error) {
	tokens := []Token{}
	for {
		t, next, err := l.Next(chars)
		if err != nil {
			return nil, err
		}
		if _, ok := l.ignored[t.Kind()]; !ok {
			tokens = append(tokens, *t)
		}
		chars = next
		if t.Kind() == l.eof {
			return tokens, nil
		}
	}
}

const (
	TokenTypeDefault = iota
	TokenTypeEOF
	TokenTypeWSpace
	TokenTypeNum
	TokenTypeLParen
	TokenTypeRParen
	TokenTypeAdd
	TokenTypeMul
)

func createLangDfa() *Dfa {
	d := NewDfa(TokenTypeDefault)
	wspace := NewDfa(TokenTypeWSpace)
	d.AddDfa([]byte(" "), wspace)
	wspace.AddDfa([]byte(" "), wspace)
	number := NewDfa(TokenTypeNum)
	d.AddDfa([]byte("0123456789"), number)
	number.AddDfa([]byte("0123456789"), number)
	d.AddPath([]byte("("), TokenTypeLParen, TokenTypeDefault)
	d.AddPath([]byte(")"), TokenTypeRParen, TokenTypeDefault)
	d.AddPath([]byte("+"), TokenTypeAdd, TokenTypeDefault)
	d.AddPath([]byte("*"), TokenTypeMul, TokenTypeDefault)
	return d
}

func langTokensToString(tokens []Token) string {
	s := &strings.Builder{}
	for _, i := range tokens {
		s.WriteByte(' ')
		switch i.Kind() {
		case TokenTypeEOF:
			s.WriteString("EOF_TOKEN")
		case TokenTypeWSpace:
			s.WriteString("WSPACE_TOKEN")
		case TokenTypeNum:
			s.WriteString(i.Val())
		case TokenTypeLParen:
			s.WriteString("(")
		case TokenTypeRParen:
			s.WriteString(")")
		case TokenTypeAdd:
			s.WriteString("+")
		case TokenTypeMul:
			s.WriteString("*")
		default:
			s.WriteString("DEFAULT_TOKEN")
		}
	}
	return s.String()
}

func createLangLexer() *Lexer {
	return NewLexer(createLangDfa(), TokenTypeDefault, TokenTypeEOF, map[int]struct{}{
		TokenTypeWSpace: {},
	})
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

	sum := 0
	sum2 := 0

	lexer := createLangLexer()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens, err := lexer.Tokenize([]byte(scanner.Text()))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(langTokensToString(tokens))

		{
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
