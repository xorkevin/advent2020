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

type (
	Parser struct {
		table map[int]map[int][]GrammarSym
	}

	GrammarSym struct {
		term bool
		kind int
	}

	GrammarRule struct {
		from int
		to   []GrammarSym
	}
)

func NewGrammarTerm(kind int) GrammarSym {
	return GrammarSym{
		term: true,
		kind: kind,
	}
}

func NewGrammarNonTerm(kind int) GrammarSym {
	return GrammarSym{
		term: false,
		kind: kind,
	}
}

func NewGrammarRule(from int, to []GrammarSym) GrammarRule {
	return GrammarRule{
		from: from,
		to:   to,
	}
}

var (
	ErrGrammar = errors.New("grammar error")
)

func NewParser(rules []GrammarRule, start, eof int) (*Parser, error) {
	// nullable set
	nonTerminals := map[int]struct{}{}
	nullableSet := map[int]struct{}{}
	for _, i := range rules {
		nonTerminals[i.from] = struct{}{}
		if len(i.to) == 0 {
			nullableSet[i.from] = struct{}{}
			continue
		}
	}
	for {
		change := false
	nullouter:
		for _, i := range rules {
			for _, j := range i.to {
				if j.term {
					continue nullouter
				}
				if _, ok := nullableSet[j.kind]; !ok {
					continue nullouter
				}
			}
			if _, ok := nullableSet[i.from]; !ok {
				change = true
				nullableSet[i.from] = struct{}{}
			}
		}
		if !change {
			break
		}
	}

	firstSet := map[int]map[int]struct{}{}
	followSet := map[int]map[int]struct{}{}
	parseTable := map[int]map[int][]GrammarSym{}
	for nt := range nonTerminals {
		firstSet[nt] = map[int]struct{}{}
		followSet[nt] = map[int]struct{}{}
		parseTable[nt] = map[int][]GrammarSym{}
	}

	// first set
	for _, i := range rules {
		if len(i.to) != 0 && i.to[0].term {
			firstSet[i.from][i.to[0].kind] = struct{}{}
		}
	}
	for {
		change := false
		for _, i := range rules {
			for _, j := range i.to {
				if j.term {
					if _, ok := firstSet[i.from][j.kind]; !ok {
						change = true
						firstSet[i.from][j.kind] = struct{}{}
					}
					break
				}
				for k := range firstSet[j.kind] {
					if _, ok := firstSet[i.from][k]; !ok {
						change = true
						firstSet[i.from][k] = struct{}{}
					}
				}
				if _, ok := nullableSet[j.kind]; !ok {
					break
				}
			}
		}
		if !change {
			break
		}
	}

	// follow set
	followSet[start][eof] = struct{}{}
	for {
		change := false
		for _, i := range rules {
		followouter:
			for n, j := range i.to {
				if j.term {
					continue
				}
				for _, k := range i.to[n+1:] {
					if k.term {
						if _, ok := followSet[j.kind][k.kind]; !ok {
							change = true
							followSet[j.kind][k.kind] = struct{}{}
						}
						continue followouter
					}
					for l := range firstSet[k.kind] {
						if _, ok := followSet[j.kind][l]; !ok {
							change = true
							followSet[j.kind][l] = struct{}{}
						}
					}
					if _, ok := nullableSet[k.kind]; !ok {
						continue followouter
					}
				}
				for k := range followSet[i.from] {
					if _, ok := followSet[j.kind][k]; !ok {
						change = true
						followSet[j.kind][k] = struct{}{}
					}
				}
			}
		}
		if !change {
			break
		}
	}

	// parsing table
	for _, i := range rules {
		for _, j := range i.to {
			if j.term {
				if _, ok := parseTable[i.from][j.kind]; ok {
					return nil, fmt.Errorf("Grammar is not LL1: %w", ErrGrammar)
				}
				parseTable[i.from][j.kind] = i.to
				break
			}
			for k := range firstSet[j.kind] {
				if _, ok := parseTable[i.from][k]; ok {
					return nil, fmt.Errorf("Grammar is not LL1: %w", ErrGrammar)
				}
				parseTable[i.from][k] = i.to
			}
			if _, ok := nullableSet[j.kind]; !ok {
				break
			}
		}
		if _, ok := nullableSet[i.from]; !ok {
			continue
		}
		for j := range followSet[i.from] {
			if _, ok := parseTable[i.from][j]; ok {
				return nil, fmt.Errorf("Grammar is not LL1: %w", ErrGrammar)
			}
			parseTable[i.from][j] = i.to
		}
	}

	return &Parser{
		table: parseTable,
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

	rules := []GrammarRule{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		ruleLine := strings.Split(line, ": ")
		if len(ruleLine) != 2 {
			log.Fatal("invalid rule line")
		}
		from, err := strconv.Atoi(ruleLine[0])
		if err != nil {
			log.Fatal(err)
		}
		tos := strings.Split(ruleLine[1], " | ")
		for _, t := range tos {
			to := []GrammarSym{}
			for _, i := range strings.Split(t, " ") {
				if len(i) == 0 {
					log.Fatal("invalid rule line")
				}
				if i[0] == '"' {
					to = append(to, NewGrammarTerm(int(i[1])))
					continue
				}
				num, err := strconv.Atoi(i)
				if err != nil {
					log.Fatal(err)
				}
				to = append(to, NewGrammarNonTerm(num))
			}
			rules = append(rules, NewGrammarRule(from, to))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if _, err := NewParser(rules, 0, -1); err != nil {
		log.Fatal(err)
	}
}
