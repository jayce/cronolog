package cronolog

import "errors"

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

func isDigit(s byte) bool {
	return s <= '9' && s >= '0'
}

// ParseSize strings: -1, 2.2b, 2
func ParseSize(s string) (Size, error) {
	origin := s
	length := len(origin)
	if length < 2 {
		return 0, errors.New("invalid " + s)
	}

	negative := false
	c := s[0]
	if c == '-' || c == '+' {
		s = s[1:]
		negative = c == '-'
	}

	digit := cutDigit(s)
	if digit == "" {
		return 0, errors.New("invalid " + s)
	}

	var (
		decimal string
		unit    string
	)

	s = s[len(digit):]
	if s[0] == '.' {
		s = s[1:]
		decimal = cutDigit(s)
		if decimal == "" {
			return 0, errors.New("invalid decimal " + s)
		}
		unit = s[len(decimal):]
	} else {
		unit = s
	}

	Unit := Byte
	if unit != "" {
		v, ok := unitMap[unit]
		if !ok {
			return 0, errors.New("invalid unit " + s)
		}
		Unit = v
	}

	size := atoi(digit)
	size *= int64(Unit)

	if decimal != "" {
		d := atoi(decimal)
		d1 := float64(d) / float64(square(10, len(decimal)))
		v := int64(d1 * (float64(Unit) / 1))
		size += v
	}

	if negative {
		size = -size
	}

	return Size(size), nil
}

func cutDigit(s string) string {
	l := len(s)
	if l == 1 && isDigit(s[0]) {
		return s
	}

	for i := 0; i < l; i++ {
		if !isDigit(s[i]) {
			return s[:i]
		}
	}

	return ""
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
