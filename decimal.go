package marketdata

import (
	"fmt"
	"math"
)

func Dtoa(mantissa int64, exponent int8) string {
	isNegative := mantissa < 0
	if isNegative {
		mantissa = -mantissa
	}

	buf := make([]rune, 0)

	for {
		a := '0' + rune(mantissa%10)
		buf = append(buf, a)
		mantissa /= 10
		exponent++
		if exponent == 0 {
			buf = append(buf, '.')
		}
		// Continue iterating while mantissa is greater than 0 or exponet is less than 1
		if !(mantissa > 0 || exponent < 1) {
			break
		}
	}
	if isNegative {
		buf = append(buf, '-')
	}

	return string(reverse(buf))
}

// Decimal type with arbitrary precision
type Decimal struct {
	Mantissa int64 `json:"mantissa"`
	Exponent int8  `json:"exponent"`
}

func (d *Decimal) String() string {
	return Dtoa(d.Mantissa, d.Exponent)
}

func reverse(s []rune) []rune {
	rev := make([]rune, 0, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		rev = append(rev, s[i])
	}
	return rev
}

// Convert ascii representation of a number to a decimal
// logic borrowed from https://github.com/jamesdbrock/hffix/blob/master/include/hffix.hpp function `atod`
func ParseDecimal(value string) Decimal {
	mantissa := int64(0)
	exponent := int8(0)

	isNegative := false
	offset := 0
	if value[0] == '-' {
		isNegative = true
		offset = 1
	}

	isDecimal := false
	for _, v := range value[offset:] {
		if v == '.' {
			isDecimal = true
		} else {
			mantissa *= 10
			mantissa += int64(v - '0')
			if isDecimal {
				exponent--
			}
		}
	}

	if isNegative {
		mantissa = -mantissa
	}
	return Decimal{Mantissa: mantissa, Exponent: exponent}
}

type Decimal9 struct {
	Mantissa int64 `json:"mantissa"`
	Exponent int8  `json:"exponent"`
}

func NewDecimal9() *Decimal9 {
	d := Decimal9{Exponent: -9}
	return &d
}

func (d9 *Decimal9) String() string {
	return Dtoa(d9.Mantissa, -9)
}

type Decimal8 struct {
	Mantissa int64 `json:"mantissa"`
	Exponent int8  `json:"exponent"`
}

func NewDecimal8() *Decimal8 {
	d := Decimal8{Exponent: -8}
	return &d
}

func (d8 *Decimal8) String() string {
	return Dtoa(d8.Mantissa, -8)
}

// Gets a Decimal9 from a Decimal
func GetDecimal9(d *Decimal) (Decimal9, error) {
	d9 := Decimal9{Exponent: -9}

	switch d.Exponent {
	case 0:
		d9.Mantissa = d.Mantissa * 1_000_000_000
		return d9, nil
	case -1:
		d9.Mantissa = d.Mantissa * 100_000_000
		return d9, nil
	case -2:
		d9.Mantissa = d.Mantissa * 10_000_000
		return d9, nil
	case -3:
		d9.Mantissa = d.Mantissa * 1_000_000
		return d9, nil
	case -4:
		d9.Mantissa = d.Mantissa * 100_000
		return d9, nil
	case -5:
		d9.Mantissa = d.Mantissa * 10_000
		return d9, nil
	case -6:
		d9.Mantissa = d.Mantissa * 1_000
		return d9, nil
	case -7:
		d9.Mantissa = d.Mantissa * 100
		return d9, nil
	case -8:
		d9.Mantissa = d.Mantissa * 10
		return d9, nil
	case -9:
		d9.Mantissa = d.Mantissa
		return d9, nil
	}
	d9.Mantissa = math.MinInt64
	return d9, fmt.Errorf("Invalid exponent %v", d.Exponent)
}

// Gets a Decimal8 from a Decimal
func GetDecimal8(d *Decimal) (Decimal8, error) {
	d8 := Decimal8{Exponent: -8}

	switch d.Exponent {
	case 0:
		d8.Mantissa = d.Mantissa * 100_000_000
		return d8, nil
	case -1:
		d8.Mantissa = d.Mantissa * 10_000_000
		return d8, nil
	case -2:
		d8.Mantissa = d.Mantissa * 1_000_000
		return d8, nil
	case -3:
		d8.Mantissa = d.Mantissa * 100_000
		return d8, nil
	case -4:
		d8.Mantissa = d.Mantissa * 10_000
		return d8, nil
	case -5:
		d8.Mantissa = d.Mantissa * 1_000
		return d8, nil
	case -6:
		d8.Mantissa = d.Mantissa * 100
		return d8, nil
	case -7:
		d8.Mantissa = d.Mantissa * 10
		return d8, nil
	case -8:
		d8.Mantissa = d.Mantissa
		return d8, nil
	}
	d8.Mantissa = math.MinInt64
	return d8, fmt.Errorf("Invalid exponent %v", d.Exponent)
}
