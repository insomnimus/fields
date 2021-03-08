package args

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type parser struct {
	text         []rune
	ch           rune
	pos, readpos int
}

// Parse parses any number of numbers given as a string, returning integers.
func Parse(s string) ([]int, error) {
	if s == "" {
		return nil, nil
	}
	p := &parser{
		text: []rune(s),
	}
	p.read()
	vals, err := p.parse()
	if err != nil {
		return nil, err
	}
	sort.Ints(vals)
	return vals, nil
}

func (p *parser) read() {
	if p.readpos >= len(p.text) {
		p.ch = 0
	} else {
		p.ch = p.text[p.readpos]
	}
	p.pos = p.readpos
	p.readpos++
}

func (p *parser) peekN(n int) rune {
	n = p.pos + n
	if n >= len(p.text) {
		return 0
	}
	return p.text[n]
}

func (p *parser) peek() rune {
	return p.peekN(1)
}

func (p *parser) parse() ([]int, error) {
	var nums []int
LOOP:
	for {
		switch {
		case p.ch == 0:
			break LOOP
		case unicode.IsDigit(p.ch):
			vals, err := p.readExpression()
			if err != nil {
				return nil, err
			}
			nums = append(nums, vals...)
		case p.ch == ',':
		default:
			return nil, fmt.Errorf("%q: not a number", p.ch)
		}
		p.read()
	}
	return nums, nil
}

func (p *parser) readExpression() ([]int, error) {
	var buff strings.Builder
	hasRange := false
LOOP:
	for {
		switch {
		case unicode.IsDigit(p.ch):
			buff.WriteString(p.readSingleNumber())
			continue LOOP
		case p.ch == '.':
			if p.peek() != '.' {
				fmt.Printf("char is %q", p.peek())
				return nil, fmt.Errorf("invalid range statement 'num.num', missing second '.'")
			}
			if hasRange {
				return nil, fmt.Errorf("invalid range expression syntax")
			}
			p.read()
			hasRange = true
			buff.WriteString("..")
		case p.ch == 0 || p.ch == ',':
			break LOOP
		default:
			return nil, fmt.Errorf("%q: not a number", p.ch)
		}
		p.read()
	}
	text := buff.String()
	if hasRange {
		split := strings.Split(text, "..")
		n1, err := strconv.Atoi(split[0])
		if err != nil {
			return nil, err
		}
		n2, err := strconv.Atoi(split[1])
		if err != nil {
			return nil, err
		}
		var nums []int
		if n1 == n2 {
			nums = append(nums, n1)
		} else if n1 < n2 {
			for i := n1; i <= n2; i++ {
				nums = append(nums, i)
			}
		} else {
			for i := n2; i <= n1; i++ {
				nums = append(nums, i)
			}
		}
		return nums, nil
	}
	n, err := strconv.Atoi(text)
	if err != nil {
		return nil, err
	}
	return []int{n}, nil
}

func (p *parser) readSingleNumber() string {
	var buff strings.Builder
	for unicode.IsDigit(p.ch) {
		buff.WriteRune(p.ch)
		p.read()
	}
	return buff.String()
}
