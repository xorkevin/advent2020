package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"xorkevin.dev/gnom"
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

	rules := []gnom.GrammarRule{}

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
		fromKind, err := strconv.Atoi(ruleLine[0])
		if err != nil {
			log.Fatal(err)
		}
		from := gnom.NewGrammarNonTerm(fromKind)
		tos := strings.Split(ruleLine[1], " | ")
		for _, t := range tos {
			to := []gnom.GrammarSym{}
			for _, i := range strings.Split(t, " ") {
				if len(i) == 0 {
					log.Fatal("invalid rule line")
				}
				if i[0] == '"' {
					to = append(to, gnom.NewGrammarTerm(int(i[1])))
					continue
				}
				num, err := strconv.Atoi(i)
				if err != nil {
					log.Fatal(err)
				}
				to = append(to, gnom.NewGrammarNonTerm(num))
			}
			rules = append(rules, gnom.NewGrammarRule(from, to...))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	const (
		tokenKindDefault = -2
		tokenKindEOF     = -1
	)

	start := gnom.NewGrammarNonTerm(0)
	eof := gnom.NewGrammarTerm(tokenKindEOF)

	dfa := gnom.NewDfa(tokenKindDefault)
	dfa.AddPath([]rune("a"), 'a', tokenKindDefault)
	dfa.AddPath([]rune("b"), 'b', tokenKindDefault)
	lexer := gnom.NewDfaLexer(dfa, tokenKindDefault, tokenKindEOF, map[int]struct{}{})

	parser := gnom.NewPEGParser(rules, start, eof)
	rules = append(rules,
		gnom.NewGrammarRule(gnom.NewGrammarNonTerm(8), gnom.NewGrammarNonTerm(42), gnom.NewGrammarNonTerm(8)),
		gnom.NewGrammarRule(gnom.NewGrammarNonTerm(11), gnom.NewGrammarNonTerm(42), gnom.NewGrammarNonTerm(11), gnom.NewGrammarNonTerm(31)),
	)
	parser2 := gnom.NewPEGParser(rules, start, eof)

	count := 0
	count2 := 0
	for scanner.Scan() {
		tokens, err := lexer.Tokenize([]rune(scanner.Text()))
		if err != nil {
			log.Fatal(err)
		}
		_, err = parser.Parse(tokens)
		if err == nil {
			count++
		}
		_, err = parser2.Parse(tokens)
		if err == nil {
			count2++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", count)
	fmt.Println("Part 2:", count2)
}
