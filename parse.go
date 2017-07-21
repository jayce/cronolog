package cronolog

import (
	"errors"
	"fmt"
)

type Size int64

const (
	Byte Size = 1
	KiB       = 1024 * Byte
	MiB       = 1024 * KiB
	GiB       = 1024 * MiB
)

const floatlen = 2

func square(n, s int) int {
	m := n
	for s > 1 {
		s--
		m *= n
	}

	return m
}

func fmtSize(buf []byte, s Size, unit Size) int {
	l := len(buf)
	float := float64(s%unit) / float64(unit)
	digit := int64(0)

	if float != 0 {
		first := true
		digit = int64(float * float64(square(10, floatlen)))

		for digit > 0 {
			n := digit % 10
			digit /= 10

			if first && n == 0 {
				continue
			}

			l--
			buf[l] = byte(n) + '0'
			first = false
		}

		if !first {
			l--
			buf[l] = '.'
		}
	}

	value := s / unit
	for value > 0 {
		l--
		buf[l] = byte(value%10) + '0'
		value /= 10
	}
	return l
}

func (s Size) String() string {
	var buf [32]byte

	l := len(buf)
	dup := s
	negative := s < 0
	if negative {
		dup = -dup
	}

	switch {
	case dup == 0:
		return "0b"
	case dup < KiB:
		l--
		buf[l] = 'b'
		l = fmtSize(buf[:l], dup, Byte)
	case dup < MiB:
		l -= 2
		copy(buf[l:], "kb")
		l = fmtSize(buf[:l], dup, KiB)
	case dup < GiB:
		l -= 2
		copy(buf[l:], "mb")
		l = fmtSize(buf[:l], dup, MiB)
	default:
		l -= 2
		copy(buf[l:], "gb")
		l = fmtSize(buf[:l], dup, GiB)
	}

	if negative {
		l--
		buf[l] = '-'
	}
	return string(buf[l:])
}

var unitMap = map[string]Size{
	"b":  Byte,
	"kb": KiB,
	"mb": MiB,
	"gb": GiB,
}

type postion struct {
	negative bool // -
	head     int  // 1
	point    int  // .
	tail     int  // 1
	unit     int  // kb
}

func parse(s string) (p *postion, err error) {
	length := len(s)
	if length == 0 {
		return nil, errors.New("empty string")
	}

	p = new(postion)
	switch s[0] {
	case '-':
		p.negative = true
		p.head++
	case '+':
		p.head++
	}

	for i := p.head; i < length; i++ {
		if s[i] > '9' || s[i] < '0' {
			if s[i] == '.' {
				p.point = i
			}
			break
		}
		p.tail = i
	}

	return
}

func ParseSize(s string) (Size, error) {
	p, err := parse(s)
	if err != nil {
		return 0, err
	}

	fmt.Printf("%s %#v", s, p)

	var (
		unit = Byte
		size = int64(0)
		head int
		tail int
	)

	head = p.unit
	if head > 0 {
		v, ok := unitMap[s[head:]]
		if !ok {
			return 0, fmt.Errorf("invalid unit '%s'", s)
		}
		unit = v
	}

	head = p.head
	tail = p.point
	if tail > 0 {
		tail--
		digit := atoi(s[head:tail])
		fmt.Println(digit, unit)
	}

	if p.negative {
		size = -size
	}

	return Size(size), nil
}

func atoi(s string) int64 {
	l := len(s)
	n := int64(0)
	for i := 0; i < l; i++ {
		n += int64(s[i] - '0')
		if i+1 != l {
			n *= 10
		}
	}
	return n
}
