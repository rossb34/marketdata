package marketdata

import "math"

type Quantity Decimal8

func QuantityInit(q *Quantity) {
	q.Exponent = -8
}

func (q *Quantity) String() string {
	return Dtoa(q.Mantissa, -8)
}

func ParseQuantity(qty string) (Quantity, error) {
	qty8 := Quantity{Mantissa: math.MinInt64, Exponent: -8}

	// Parse the quantity string into a decimal with precision according to string format
	d := ParseDecimal(qty)

	// Convert from arbitrary precision decimal to decimal with fixed exponent
	d8, err := GetDecimal8(&d)
	if err != nil {
		return qty8, err
	}
	qty8.Mantissa = d8.Mantissa
	return qty8, nil
}
