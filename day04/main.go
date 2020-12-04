package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	Passport struct {
		byr int
		iyr int
		eyr int
		hgt string
		hcl string
		ecl string
		pid string
	}
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

	passports := []Passport{}

	p := Passport{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			passports = append(passports, p)
			p = Passport{}
			continue
		}
		f := strings.Fields(line)
		for _, i := range f {
			k := strings.Split(i, ":")
			if len(k) != 2 {
				log.Fatal("invalid kv pair")
			}
			v := k[1]
			var err error
			switch k[0] {
			case "byr":
				p.byr, err = strconv.Atoi(v)
				if err != nil {
					log.Fatal(err)
				}
			case "iyr":
				p.iyr, err = strconv.Atoi(v)
				if err != nil {
					log.Fatal(err)
				}
			case "eyr":
				p.eyr, err = strconv.Atoi(v)
				if err != nil {
					log.Fatal(err)
				}
			case "hgt":
				p.hgt = v
			case "hcl":
				p.hcl = v
			case "ecl":
				p.ecl = v
			case "pid":
				p.pid = v
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	passports = append(passports, p)

	count := 0
	count2 := 0
	for _, i := range passports {
		if passIsValid(i) {
			count++
			if passIsValid2(i) {
				count2++
			}
		}
	}

	fmt.Println("Part 1:", count)
	fmt.Println("Part 2:", count2)
}

func passIsValid(p Passport) bool {
	return p.byr != 0 &&
		p.iyr != 0 &&
		p.eyr != 0 &&
		p.hgt != "" &&
		p.hcl != "" &&
		p.ecl != "" &&
		p.pid != ""
}

var (
	reHeight    = regexp.MustCompile(`^([0-9]+)(cm|in)$`)
	reColor     = regexp.MustCompile(`^#[0-9a-f]{6}$`)
	rePid       = regexp.MustCompile(`^[0-9]{9}$`)
	eyeColorSet = map[string]struct{}{
		"amb": {},
		"blu": {},
		"brn": {},
		"gry": {},
		"grn": {},
		"hzl": {},
		"oth": {},
	}
)

func passIsValid2(p Passport) bool {
	if !inRange(p.byr, 1920, 2002) {
		return false
	}
	if !inRange(p.iyr, 2010, 2020) {
		return false
	}
	if !inRange(p.eyr, 2020, 2030) {
		return false
	}
	hm := reHeight.FindStringSubmatch(p.hgt)
	if hm == nil {
		return false
	}
	h, _ := strconv.Atoi(hm[1])
	switch hm[2] {
	case "cm":
		if !inRange(h, 150, 193) {
			return false
		}
	case "in":
		if !inRange(h, 59, 76) {
			return false
		}
	default:
		return false
	}
	if !reColor.MatchString(p.hcl) {
		return false
	}
	if _, ok := eyeColorSet[p.ecl]; !ok {
		return false
	}
	if !rePid.MatchString(p.pid) {
		return false
	}
	return true
}

func inRange(i, a, b int) bool {
	return i >= a && i <= b
}
