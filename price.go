package marketdata

import "math"

type Price Decimal9

func PriceInit(p *Price) {
	p.Exponent = -9
}

func (p *Price) String() string {
	return Dtoa(p.Mantissa, -9)
}

func ParsePrice(px string) (Price, error) {
	px9 := Price{Mantissa: math.MinInt64, Exponent: -9}

	// Parse the price string into a decimal with precision according to string format
	d := ParseDecimal(px)

	// Convert from arbitrary precision decimal to decimal with fixed exponent
	d9, err := GetDecimal9(&d)
	if err != nil {
		return px9, err
	}
	px9.Mantissa = d9.Mantissa
	return px9, nil
}
