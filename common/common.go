package common

import (
	"io/ioutil"
	"os"
	"strings"
)

func OpenRuneFile(name string) (runes [][]rune) {
	file, err := os.Open(name)
	if err != nil {
		panic(err) // TODO
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err) // TODO
		}
	}()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err) // TODO
	}

	lines := strings.Split(string(b), "\n")

	for _, line := range lines {
		if line == "\n" || line == "" {
			break
		}
		var xs []rune
		for _, c := range line {
			xs = append(xs, c)
		}
		runes = append(runes, xs)
	}

	for i, rs := range runes {
		runes[i] = []rune(strings.ReplaceAll(string(rs), "░", " "))
	}

	return
}

func ReflectHorizontal(runes [][]rune) [][]rune {
	runes = reverseRunes(runes)
	for i, rs := range runes {
		for j, r := range rs {
			switch r {
			case '`':
				r = '´'
			case '(':
				r = ')'
			case ')':
				r = '('
			case '▙':
				r = '▟'
			case '▟':
				r = '▙'
			case '▜':
				r = '▛'
			case '▛':
				r = '▜'
			case '\\':
				r = '/'
			case '/':
				r = '\\'
			}
			runes[i][j] = r
		}
	}

	return runes
}

func reverseRunes(runes [][]rune) [][]rune {
	for i, rs := range runes {
		for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
			rs[i], rs[j] = rs[j], rs[i]
		}
		runes[i] = rs
	}
	return runes
}
